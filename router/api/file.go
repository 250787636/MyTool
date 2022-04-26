package api

import (
	"archive/zip"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/router/api/jiagu"
	"example.com/m/utils"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 获取文件信息接口
func GetFile(c *gin.Context) {
	var data struct {
		FileName string `json:"file_name"`
		FileSize int64  `json:"file_size"`
		FilePath string `json:"file_path"`
	}
	useType, useTypeOk := c.GetPostForm("use_type")
	useTypeOk = utils.IsStringEmpty(useType, useTypeOk)
	if !useTypeOk {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	// 保存文件到本地
	fileUuid := uuid.New()
	var filePathString string
	// 拼接路径保存文件
	if useType == "ios" {
		filePathString = fmt.Sprintf("media/%s/%s", useType, file.Filename)
	} else {
		filePathString = fmt.Sprintf("datastorage/ak/upload/scan_upload/%s/%s/%s", useType, fileUuid, file.Filename)
	}
	err = os.MkdirAll(filepath.Dir(filePathString), os.ModePerm)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	if err := c.SaveUploadedFile(file, filePathString); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	data.FileName = file.Filename
	data.FileSize = file.Size
	data.FilePath = filePathString

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": data,
		"msg":  model.ReqSuccess,
	})
}

// DownloadFile 批量下载源文件
func DownloadFile(c *gin.Context) {

	taskIdString := c.Query("file")
	fileType := c.Query("type")
	str1 := strings.ReplaceAll(taskIdString, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")

	var TaskInfo []model.CePingUserTask
	if err := app.DB.Debug().Model(&model.CePingUserTask{}).Where("task_id in (?)", idArray).Find(&TaskInfo).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 2.准备zip
	out, err := os.Create("test.zip")
	if err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)
		return
	}

	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
		}
	}(out)

	writer := zip.NewWriter(out)

	sufFIx := ".apk"
	switch fileType {
	case "apk":
		sufFIx = ".apk"
	case "ios":
		sufFIx = ".ios"
	case "mp":
		sufFIx = ".jpg"
	}

	var TaskMap = make(map[string]string)

	// 1.把应用名去后缀
	for _, task := range TaskInfo {
		_, ok := TaskMap[task.AppName]
		// 如果没有就放进去
		if !ok {
			TaskMap[task.AppName] = task.FilePath
		}
	}

	for key, task := range TaskMap {
		appName := strings.Split(key, ".")
		//fmt.Println("task", task.FilePath)
		fileWriter, err := writer.Create(appName[0] + sufFIx)
		if err != nil {
			if os.IsPermission(err) {
				log.Error("err:", err.Error())
				response.FailWithMessage(err.Error(), c)
				return
			}
			log.Error("Create file %s error: %s\n", appName[0], err.Error())
			response.FailWithMessage(err.Error(), c)

			return
		}

		fileInfo, err := os.Open(task)
		if err != nil {
			log.Error("Open file %s error: %s\n", task, err.Error())
			response.FailWithMessage(err.Error(), c)

			return
		}
		fileBody, err := ioutil.ReadAll(fileInfo)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			log.Error("Read file %s error: %s\n", task, err.Error())
			return
		}

		_, err = fileWriter.Write(fileBody)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			log.Error("err:", err.Error())
			return
		}
	}

	if err := writer.Close(); err != nil {
		log.Error("err:", err.Error())
		response.FailWithMessage(err.Error(), c)

		fmt.Println("Close error: ", err)
		return
	}
	fileContentDisposition := "inline;filename=测评源文件下载.zip"
	c.Header("Content-Type", "application/zip") // 这里是压缩文件类型 .zip
	c.Header("Content-Disposition", fileContentDisposition)

	c.File("test.zip")

}
func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

// 6.3 源码加固下载客户端
func DownlandIosFile(c *gin.Context) {
	browserName, browserNameOK := c.GetQuery("browser_name")
	systemType, systemTypeOK := c.GetQuery("system_type")

	if !(systemTypeOK && browserNameOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}

	if _, err := os.Open("media/ios/" + jiagu.FILENAME); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.FileNotExist,
		})
		return
	}
	// 用户名
	_, ok := c.Get("userName")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.UserNotExist,
		})
		return
	}
	userName, _ := c.Get("userName")
	// 下载时间
	downloadTime := time.Now().Format("2006-01-02 15:04:05")
	// ip地址
	_, ok = c.Get("ip")
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.UserNotExist,
		})
		return
	}
	ip, _ := c.Get("ip")

	data := model.Source_Code_Reinforcement_Log{
		UserName:     userName.(string),
		DownloadTime: downloadTime,
		FileName:     jiagu.FILENAME,
		IpAddress:    ip.(string),
		BrowserName:  browserName,
		SystemType:   systemType,
	}
	if err := app.DB.Model(model.Source_Code_Reinforcement_Log{}).Save(&data).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	//application/octet-stream
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment;filename="+jiagu.FILENAME)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File("media/ios/" + jiagu.FILENAME)
}
