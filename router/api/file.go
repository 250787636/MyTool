package api

import (
	"archive/zip"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/utils"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
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
	// 拼接路径保存文件
	filePathString := fmt.Sprintf("media/%s/%s/%s", useType, fileUuid, file.Filename)
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

// NewExcel 导入excel
func NewExcel(c *gin.Context) {

	excelFile, err := c.FormFile("file")
	if err != nil {
		log.Error(err.Error())
		response.FailWithMessage("获取上传文件失败", c)
		return
	}

	// 保存文件到本地
	fileUuid := uuid.New()
	// 拼接路径保存文件
	filePathString := fmt.Sprintf("media/%s/%s", fileUuid, excelFile.Filename)
	err = os.MkdirAll(filepath.Dir(filePathString), os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		response.FailWithMessage("保存文件失败", c)
		return
	}

	if err := c.SaveUploadedFile(excelFile, filePathString); err != nil {
		log.Error(err.Error())
		response.FailWithMessage("保存文件失败", c)
		return
	}

	f, err := excelize.OpenFile(filePathString)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Error(err.Error())
		}
	}()

	// 0.获取总数据
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Println("len", len(rows))

	// 1.获取excel有效行数
	effectiveRow := 0
	for rowIn, row := range rows {
		if len(row) != 6 {
			continue
		}
		fmt.Printf("第 %d 行有%d个数据\n", rowIn, len(row))
		effectiveRow++
	}
	// 2.创建空二维数组
	temArray := make([][]string, effectiveRow)
	lieTotal := len(rows[0])
	for i := range temArray {
		temArray[i] = make([]string, lieTotal)
	}
	fmt.Println("1111", temArray)
	temRow := 0
	// 3.将excel信息转存到二维数组
	for rowIn, row := range rows {
		// 如果该行数据为1表示只有序号值，就不对该行进行获取
		if len(row) != 6 {
			continue
		}
		fmt.Printf("第 %d 行有%d个数据\n", rowIn, len(row))

		for i := 0; i < len(row); i++ {
			fmt.Printf("第%d个值为 %s\n", i, row[i])
			temArray[temRow][i] = row[i]
		}
		temRow++
	}
	fmt.Println("获取到的数据", temArray)

	// 4.获取二维数组信息
	for rowIndex, row := range temArray {
		//跳过第一行表头信息
		if rowIndex == 0 {
			continue
		}
		//遍历每一个单元
		application := model.Application{}
		application.AppName = row[1]
		application.ModelName = row[2]
		application.AppCnName = row[3]
		application.AppVersion = row[4]
		application.ModelCnName = row[5]
		application.LastChangeTime = time.Now().Format("2006-01-02 15:04:05")
		application.TheApp = application.AppName + "-" + application.AppCnName
		application.TheModel = application.ModelName + "-" + application.ModelCnName

		// 5.写入mysql
		if err := app.DB.Debug().Model(model.Application{}).Create(&application).Error; err != nil {
			log.Error(err.Error())
			response.FailWithMessage("根据模版创建系统失败", c)
			return
		}

	}
	response.OkWithMessage("导入excel文件成功", c)
}
