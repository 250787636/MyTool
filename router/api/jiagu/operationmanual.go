package jiagu

import (
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 10.1创建手册
func HandBookCreate(c *gin.Context) {
	var operationManual model.JiaguOperationManual
	// 获取服务名 和文件
	num, serviceIdOK := c.GetPostForm("service_id")
	serviceId, _ := strconv.Atoi(num)
	if !serviceIdOK {
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

	// 存储文件
	fileuuid := uuid.New().String()
	filePath := fmt.Sprintf("media/handbook/%s/%s", fileuuid, file.Filename)
	// 创建文件目录
	err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	// 为存入数据库的结构体赋值
	operationManual.ServiceId = serviceId
	operationManual.FileName = file.Filename
	operationManual.FilePath = filePath
	if err := app.DB.Model(model.JiaguOperationManual{}).Create(&operationManual).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  model.CreSuccess,
	})
}

// 10.2获取手册信息
func HandBookGetAll(c *gin.Context) {
	var operationMessage []map[string]interface{}

	page, size, offNum := utils.GetFromDataPageSizeOffNum(c)
	if page+size+offNum == -3 {
		return
	}

	var count int64
	if err := app.DB.Table("jiagu_operation_manual").
		Select("jiagu_operation_manual.created_at,jiagu_operation_manual.id,jiagu_operation_manual.file_name,jiagu_operation_manual.file_path,service_type.service_type").
		Joins("INNER JOIN service_type ON jiagu_operation_manual.service_id = service_type.id").
		Where("jiagu_operation_manual.deleted_at is NULL ").
		Count(&count).
		Offset(offNum).
		Limit(size).
		Order("Created_at desc").
		Scan(&operationMessage).
		Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	for _, v := range operationMessage {
		v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": gin.H{
			"datalist": operationMessage,
			"total":    count,
		},
		"msg": model.ReqSuccess,
	})
}

// 10.3下载手册
func HandBookDownland(c *gin.Context) {
	filePath, filePathOK := c.GetQuery("file_path")
	fileName, fileNameOK := c.GetQuery("file_name")
	if !(filePathOK && fileNameOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}

	if _, err := os.Open(filePath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.FileNotExist,
		})
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Transfer-Encoding", "binary")
	c.File(filePath)
}

// 10.4删除手册
func HandBookDelete(c *gin.Context) {
	var manualTask model.JiaguOperationManual
	num, handBookIdOK := c.GetPostForm("id")
	handBookId, _ := strconv.Atoi(num)
	filePath, filePathOK := c.GetPostForm("file_path")
	if !(handBookIdOK && filePathOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	manualTask.ID = uint(handBookId)
	if err := app.DB.Delete(&manualTask).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.FileAlreadyDelete,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  model.DelSuccess,
	})
}

// 10.4 获取服务名
func HandBookGetServiceType(c *gin.Context) {
	var serviceType []model.ServiceType
	if err := app.DB.Model(model.ServiceType{}).Find(&serviceType).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": gin.H{
			"datalist": serviceType,
		},
		"msg": model.ReqSuccess,
	})
}
