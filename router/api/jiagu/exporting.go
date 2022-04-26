package jiagu

import (
	"bytes"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"net/http"
	"net/url"
	"strconv"
)

type JiaGuData struct {
	ID           uint   `json:"id"`
	AppName      string `json:"app_name"`
	AppType      string `json:"app_type"`
	AppVersion   string `json:"app_version"`
	CreatedAt    string `json:"created_at"`
	FinishTime   string `json:"finish_time"`
	PolicyId     int    `json:"policy_id"`
	PolicyReason string `json:"policy_reason"`
	TaskId       int    `json:"task_id"`
	TaskStatus   string `json:"task_status"`
	UserName     string `json:"user_name"`
}

// 加固记录导出功能
func Exporting(c *gin.Context) {
	var list []JiaGuData
	var fileName string
	num, typeIdOK := c.GetQuery("type_id")
	typeIdOK = utils.IsStringEmpty(num, typeIdOK)
	if !typeIdOK {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	typeId, _ := strconv.Atoi(num)
	// 文件名赋值
	if typeId == 1 {
		fileName = "android加固记录.xlsx"
	}

	if err := app.DB.Table("jia_gu_task").
		Select(" jia_gu_task.id,user.user_name,jia_gu_task.task_id,jia_gu_task.policy_id,jia_gu_task.policy_reason,jia_gu_task.created_at,jia_gu_task.task_status,jia_gu_task.finish_time").
		Joins("INNER JOIN user ON user.id = jia_gu_task.user_id").
		Scan(&list).
		Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	//导出
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	row = sheet.AddRow()
	row.AddCell().Value = "记录ID"
	row.AddCell().Value = "用户名"
	row.AddCell().Value = "策略id"
	row.AddCell().Value = "不使用推荐策略理由"
	row.AddCell().Value = "任务id"
	row.AddCell().Value = "任务状态"
	row.AddCell().Value = "提交时间"
	row.AddCell().Value = "完成时间"
	for _, v := range list {
		row = sheet.AddRow()
		row.AddCell().Value = strconv.Itoa(int(v.ID))
		row.AddCell().Value = v.UserName
		row.AddCell().Value = strconv.Itoa(v.PolicyId)
		row.AddCell().Value = v.PolicyReason
		row.AddCell().Value = strconv.Itoa(v.TaskId)
		row.AddCell().Value = v.TaskStatus
		row.AddCell().Value = v.CreatedAt
		row.AddCell().Value = v.FinishTime
	}
	buf := new(bytes.Buffer)
	err = file.Write(buf)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, url.QueryEscape(fileName)))
	c.Data(http.StatusOK, "text/xlsx", buf.Bytes())
}
