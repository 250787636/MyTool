package jiagu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 加固内置账号使用的数据
var JIAGUUSERNAME, APIKEY, APISECRET, JIAGUIP, FILENAME = JiaGuLoading()

// 请求频繁，共用连接
var client = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: false,
	},
}

// 加固内置账户
func JiaGuLoading() (string, string, string, string, string) {
	user := app.Conf.JiaGu
	return user.UserName, user.ApiKey, user.ApiSecret, user.Ip, user.FileName
}

// android部分
// 5.1上传apk加固包
func WebBoxV5Upload(c *gin.Context) {
	// 任务表
	var jiaGuTask model.JiaGuTask
	// 获取前端传入的参数
	num3, policyIdOK := c.GetPostForm("policy_id")
	apkPath, apkPathOK := c.GetPostForm("apk_path")
	channelPath, channelPathOK := c.GetPostForm("channel_path")

	if !(policyIdOK && apkPathOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	// 获取策略id
	policyId, err := strconv.Atoi(num3)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	urlString := fmt.Sprintf("%s/v5/protect/upload?username=%s&policy_id=%d&upload_type=%d",
		JIAGUIP, JIAGUUSERNAME, policyId, 2)
	m := map[string]interface{}{
		"username":    JIAGUUSERNAME,
		"policy_id":   policyId,
		"upload_type": 2,
	}

	file, err := os.Open(apkPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": "获取apk包文件失败",
		})
		return
	}

	var res map[string]interface{}
	if apkPathOK && channelPathOK {
		// 非必传渠道加固包文件
		file2, err := os.Open(channelPath)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  "获取渠道文件失败",
			})
			return
		}
		fileMap := map[string]interface{}{
			"apk_file":     file,
			"channel_file": file2,
		}
		res, err = postFile(urlString, m, fileMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
	} else {
		fileMap := map[string]interface{}{
			"apk_file": file,
		}
		res, err = postFile(urlString, m, fileMap)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
	}
	_, exist := res["info"]
	// res["code"]为1时为出错
	if !exist || res["code"].(float64) == 1 {
		c.JSON(http.StatusOK, gin.H{
			"code": res["code"],
			"info": res["info"],
			"err":  res["msg"],
		})
		return
	}
	//var datalist struct {
	//	ApkName     string `json:"apk_name"`
	//	ApkSize     string `json:"apk_size"`
	//	Filename    string `json:"filename"`
	//	Id          int    `json:"id"`
	//	PackageName string `json:"package_name"`
	//	SrcApkMd5   string `json:"src_apk_md5"`
	//	StatusCode  int    `json:"status_code"`
	//	Version     string `json:"version"`
	//}
	datalist := res["info"].(map[string]interface{})

	// 获取当前用户id
	userId, _ := c.Get("userId")

	jiaGuTask.UserID = userId.(uint)
	jiaGuTask.PolicyID = policyId
	jiaGuTask.TaskStatus = "加固中"
	jiaGuTask.ApkName = datalist["apk_name"].(string)
	jiaGuTask.ApkSize = datalist["apk_size"].(string)
	jiaGuTask.Filename = datalist["filename"].(string)
	jiaGuTask.Version = datalist["version"].(string)

	// 进行策略表计数操作
	var policy model.JiaguPolicy
	if err := app.DB.Model(model.JiaguPolicy{}).
		Where("id = ?", policyId).First(&policy).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.PolicyNotExist,
		})
		return
	}
	if err := app.DB.Model(model.JiaguPolicy{}).
		Where("id = ?", policyId).Update("number_of_use", policy.NumberOfUse+1).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.PolicyCountError,
		})
		return
	}

	if arr, ok := res["info"].(map[string]interface{}); ok {
		jiaGuTask.TaskID = int(arr["id"].(float64))
	}

	if err := app.DB.Save(&jiaGuTask).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": res["code"],
		"info": gin.H{
			"datalist": res["info"],
		},
		"msg": res["msg"],
	})

	// 轮询状态
	go getStateOK(jiaGuTask.TaskID, "android")
}

// 5.2查询apk加固状态
func WebBoxV5GetState(c *gin.Context) {
	// 获取任务id
	taskID, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	urlString := fmt.Sprintf("%s/v5/protect/get_state?username=%s&apkinfo_id=%d",
		JIAGUIP, JIAGUUSERNAME, taskID)
	m := map[string]interface{}{
		"username":   JIAGUUSERNAME,
		"apkinfo_id": taskID,
	}
	res, err := PostWithoutFile(urlString, m)
	ress := map[string]interface{}{}
	_ = jsoniter.Unmarshal(res, &ress)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": ress["code"],
			"info": gin.H{
				"datalist": ress["info"],
			},
			"msg": ress["msg"],
		})
	}
}

// 5.3下载apk加固包
func WebBoxV5Download(c *gin.Context) {
	// 获取任务id
	num, taskIDOk := c.GetQuery("id")
	// 下载方式:0-nfs,1-二进制
	num2, downloadTypeOK := c.GetQuery("download_type")
	if !(taskIDOk && downloadTypeOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
	}
	taskID, _ := strconv.Atoi(num)
	downloadType, _ := strconv.Atoi(num2)

	urlString := fmt.Sprintf("%s/v5/protect/download?username=%s&apkinfo_id=%d&download_type=%d",
		JIAGUIP, JIAGUUSERNAME, taskID, downloadType)
	m := map[string]interface{}{
		"username":      JIAGUUSERNAME,
		"apkinfo_id":    taskID,
		"download_type": downloadType,
	}
	response, err := postDowloand(urlString, m)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
	} else {
		// 从response 获取下载的文件头
		header := response.Header.Get("Content-Disposition")
		c.Writer.Header().Set("Content-Disposition", header)
		body, _ := ioutil.ReadAll(response.Body)
		_, _ = c.Writer.Write(body)
	}
}

// 5.4下载apk加固日志
func WebBoxV5DownloadLog(c *gin.Context) {
	// 获取任务id
	num, taskIDOK := c.GetQuery("id")
	if !taskIDOK {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
	}
	taskID, _ := strconv.Atoi(num)

	urlString := fmt.Sprintf("%s/v5/protect/download_log?username=%s&apkinfo_id=%d",
		JIAGUIP, JIAGUUSERNAME, taskID)
	m := map[string]interface{}{
		"username":   JIAGUUSERNAME,
		"apkinfo_id": taskID,
	}
	response, err := postDowloand(urlString, m)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
	} else {
		// 从response 获取下载的文件头
		header := response.Header.Get("Content-Disposition")
		c.Writer.Header().Set("Content-Disposition", header)
		body, _ := ioutil.ReadAll(response.Body)
		_, _ = c.Writer.Write(body)
	}
}

// 5.5删除apk加固记录
func WebBoxV5Delete(c *gin.Context) {
	// 获取任务id
	taskArrIds, taskArrIdOK := c.GetPostForm("ids")
	taskArrIdOK = utils.IsStringEmpty(taskArrIds, taskArrIdOK)
	if !taskArrIdOK {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}

	// 将String数组转为Int数组
	str1 := strings.ReplaceAll(taskArrIds, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArrayStr := strings.Split(str2, ",")
	idArrayInt := make([]int, len(idArrayStr))
	for i := 0; i < len(idArrayStr); i++ {
		idArrayInt[i], _ = strconv.Atoi(idArrayStr[i])
	}

	log.Info(idArrayInt)

	// 进行任务删除
	for _, v := range idArrayInt {
		if err := app.DB.Where("task_id = ?", v).Delete(&model.JiaGuTask{}).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  model.DelSuccess,
	})
}

// 5.6获取当前用户加固任务
func WebBoxV5AllTask(c *gin.Context) {
	// 开始和结束日期
	var startTime time.Time
	var endTime time.Time
	var err error
	// 获取传入的值
	status, statusOK := c.GetPostForm("app_status")
	Time1, startTimeOK := c.GetPostForm("start_time")
	Time2, endTimeOk := c.GetPostForm("end_time")
	userName, userNameOk := c.GetPostForm("user_name")
	fileName, fileNameOk := c.GetPostForm("filename")        // 应用名称
	apkName, apkNameOk := c.GetPostForm("apk_name")          // apk文件名
	policyName, policyNameOk := c.GetPostForm("policy_name") // 策略名称

	// 判断是否为空字符
	statusOK = utils.IsStringEmpty(status, statusOK)
	startTimeOK = utils.IsStringEmpty(Time1, startTimeOK)
	endTimeOk = utils.IsStringEmpty(Time2, endTimeOk)
	userNameOk = utils.IsStringEmpty(userName, userNameOk)
	fileNameOk = utils.IsStringEmpty(fileName, fileNameOk)
	apkNameOk = utils.IsStringEmpty(apkName, apkNameOk)
	policyNameOk = utils.IsStringEmpty(policyName, policyNameOk)

	// 将string转为datetime
	if startTimeOK && endTimeOk {
		startTime, err = time.Parse("2006-01-02 15:04", Time1)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
		endTime, err = time.Parse("2006-01-02 15:04", Time2)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
	}

	// 将page size offnum封装成工具方法
	page, size, offNum := utils.GetFromDataPageSizeOffNum(c)
	if page+size+offNum == -3 {
		return
	}

	// 用户权限
	//isAdmin, _ := c.Get("isAdmin")
	//superAdmin, _ := c.Get("superAdmin")
	var userType interface{}
	if _, exist := c.Get("userType"); exist {
		userType, _ = c.Get("userType")
	}
	isAdmin := userType.(string) == "1"

	// 总数
	var total int64
	var taskMessage []map[string]interface{}
	// 拼接sql语句
	var sqlString bytes.Buffer

	layout := "2006-01-02 15:04:05"

	sqlString.WriteString("jia_gu_task.deleted_at is NULL")
	// 不为超级管理员 进行拼接条件
	if isAdmin { // 为部门管理员账户则
		// 添加用户名检索功能(只有管理员能进行用户名检索)
		if userNameOk {
			sqlString.WriteString(" AND user_info.username like '")
			sqlString.WriteString("%" + userName + "%")
			sqlString.WriteString("'")
		}
	} else { // 为普通用户
		// 获取当前用户id
		num3, _ := c.Get("userId")
		userId := int(num3.(uint))
		sqlString.WriteString(" AND jia_gu_task.user_id = ")
		sqlString.WriteString(strconv.Itoa(userId))
	}

	// 添加时间条件
	if startTimeOK && endTimeOk {
		sqlString.WriteString(" AND jia_gu_task.created_at BETWEEN ")
		sqlString.WriteString(`"`)
		sqlString.WriteString(startTime.Format(layout))
		sqlString.WriteString(`"`)
		sqlString.WriteString(" AND ")
		sqlString.WriteString(`"`)
		sqlString.WriteString(endTime.Format(layout))
		sqlString.WriteString(`"`)
	}
	// 添加加固状态检索
	if statusOK {
		sqlString.WriteString(" AND jia_gu_task.task_status = ")
		sqlString.WriteString(`"`)
		sqlString.WriteString(status)
		sqlString.WriteString(`"`)
	}
	// apk文件名模糊查询
	if fileNameOk {
		sqlString.WriteString(" AND jia_gu_task.filename like '")
		sqlString.WriteString("%" + fileName + "%")
		sqlString.WriteString("'")
	}

	// fileName文件名模糊查询
	if apkNameOk {
		sqlString.WriteString(" AND jia_gu_task.apk_name like '")
		sqlString.WriteString("%" + apkName + "%")
		sqlString.WriteString("'")
	}
	// 策略名称
	if policyNameOk {
		sqlString.WriteString(" AND jiagu_policy.name like '")
		sqlString.WriteString("%" + policyName + "%")
		sqlString.WriteString("'")
	}

	// 获取分页查询数据 | 获取总数total
	if err := app.DB.Table("jia_gu_task").
		Select("user_info.username," +
			"jia_gu_task.task_id,jia_gu_task.apk_name,jia_gu_task.apk_size,jia_gu_task.filename,jia_gu_task.version" +
			",jiagu_policy.name,jia_gu_task.created_at,jia_gu_task.task_status,jia_gu_task.finish_time").
		Joins("INNER JOIN user_info ON user_info.id = jia_gu_task.user_id").
		Joins("INNER JOIN jiagu_policy ON jiagu_policy.id = jia_gu_task.policy_id").
		Where(sqlString.String()).
		Count(&total).
		Offset(offNum).
		Limit(size).
		Order("jia_gu_task.created_at desc").
		Scan(&taskMessage).
		Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	// 格式化创建时间与完成时间
	for _, v := range taskMessage {
		v["created_at"] = v["created_at"].(time.Time).Format(layout)
		if v["finish_time"] != nil {
			v["finish_time"] = v["finish_time"].(time.Time).Format(layout)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": taskMessage,
			"total":    total,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})
}

// 源码加固数据
// 6.1获取源码加固数据
func WebBoxFindIosEdition(c *gin.Context) {
	var data model.IosClientDownloadPath
	if err := app.DB.First(&data).Error; err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": data,
		"msg":  model.ReqSuccess,
	})
}

// 6.2修改源码加固数据
func WebBoxModifyIosEdition(c *gin.Context) {
	var data model.IosClientDownloadPath
	edition := c.PostForm("edition")
	windowsName := c.PostForm("windows_name")
	windowsPath := c.PostForm("windows_path")
	macName := c.PostForm("mac_Name")
	macPath := c.PostForm("mac_path")

	// 只有一条数据 所有默认赋值id
	data.Id = 1
	data.IosClientEdition = edition
	data.IosWindowsName = windowsName
	data.IosWindowsDownloadPath = windowsPath
	data.IosMacName = macName
	data.IosMacDownloadPath = macPath

	if err := app.DB.Save(&data).Error; err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  model.ModSuccess,
	})
}

// 6.4 客户端下载记录
func IosDownlandList(c *gin.Context) {
	var datalist []model.Source_Code_Reinforcement_Log
	var query string
	var args []interface{}
	startTime, isStatTime := c.GetPostForm("start_time")
	if utils.IsStringEmpty(startTime, isStatTime) {
		query = "AND download_time >= ? "
		args = append(args, startTime)
	}
	endTime, isEndTime := c.GetPostForm("end_time")
	if utils.IsStringEmpty(endTime, isEndTime) {
		query += "AND download_time < ? "
		args = append(args, endTime)
	}
	userName, isUserNmae := c.GetPostForm("username")
	if utils.IsStringEmpty(userName, isUserNmae) {
		query += "AND user_name like ? "
		args = append(args, fmt.Sprintf("%%%s%%", utils.SpecialString(userName)))
	}

	query = strings.TrimPrefix(query, "AND")

	// 将page size offnum封装成工具方法
	page, size, offNum := utils.GetFromDataPageSizeOffNum(c)
	if page+size+offNum == -3 {
		return
	}
	var total int64
	if err := app.DB.
		Model(model.Source_Code_Reinforcement_Log{}).
		Where(query, args...).
		Order("id DESC").
		Count(&total).
		Offset(offNum).
		Limit(size).
		Find(&datalist).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": gin.H{
			"datalist": datalist,
			"total":    total,
		},
		"msg": model.ReqSuccess,
	})
}

// 加固策略部分
// 13.1获取加固策略
func JiaguPolicyFind(c *gin.Context) {
	// 除标准策略以外的原数据
	var datas []model.JiaguPolicy
	// 最终的数据
	var finalDatas []model.JiaguPolicy
	var sqlString bytes.Buffer

	// 拼接sql语句
	sqlString.WriteString(`status = "1" `)

	var userID interface{}
	// 用户id
	if _, exist := c.Get("userId"); exist {
		userID, _ = c.Get("userId")
	}
	sUserId := strconv.Itoa(int(userID.(uint)))

	// 用户权限
	var userType interface{}
	if _, exist := c.Get("userType"); exist {
		userType, _ = c.Get("userType")
	}
	isAdmin := userType.(string) == "1"
	// 有应用id时
	if !isAdmin {
		sqlString.WriteString(" AND user_ids like '")
		sqlString.WriteString("%" + sUserId + "%")
		sqlString.WriteString("'")
	}
	//  查出id对应的数据
	if err := app.DB.Model(model.JiaguPolicy{}).Where(sqlString.String()).Find(&datas).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	if !isAdmin {
		// 筛选出包含appid的数据
		for _, v := range datas {
			arr := strings.Split(v.UserIds, ",")
			for _, id := range arr {
				if sUserId == id {
					finalDatas = append(finalDatas, v)
					break
				}
			}
		}
	} else {
		finalDatas = datas
	}

	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": finalDatas,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})
}

// 13.2获取加固策略 并分页
func JiaguPolicyFindWithPage(c *gin.Context) {
	var datas []model.JiaguPolicy
	var sqlString bytes.Buffer
	var total int64

	// 分页
	page, size, offNum := utils.GetFromDataPageSizeOffNum(c)
	if page+size+offNum == -3 {
		return
	}

	// 拼接sql语句
	sqlString.WriteString(`status = "1" `)

	if err := app.DB.Model(model.JiaguPolicy{}).
		Where(sqlString.String()).
		Count(&total).
		Offset(offNum).
		Limit(size).
		Find(&datas).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": datas,
			"total":    total,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})
}

// 13.3加固策略进行关联用户
func JiaguPolicyRelated(c *gin.Context) {
	userIds, userIdsOK := c.GetPostForm("user_ids")
	policyId, policyIdOk := c.GetPostForm("policyId")
	policyIdOk = utils.IsStringEmpty(policyId, policyIdOk)
	userIdsOK = utils.IsStringEmpty(userIds, userIdsOK)

	if !(policyIdOk && userIdsOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	if err := app.DB.Model(model.JiaguPolicy{}).
		Where("id = ?", policyId).Update("user_ids", userIds).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})
}

// post不带参数和文件
func PostWithoutFile(urlString string, m map[string]interface{}) ([]byte, error) {
	req, err := utils.NewFormDataRequest(urlString, nil, nil)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	req.Header.Set("api_key", APIKEY)
	req.Header.Set("sign", getSign(m))
	resp, err := client.Do(req)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

// post请求下载文件
func postDowloand(urlString string, m map[string]interface{}) (*http.Response, error) {
	req, err := utils.NewFormDataRequest(urlString, nil, nil)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	req.Header.Set("api_key", APIKEY)
	req.Header.Set("sign", getSign(m))
	resp, err := client.Do(req)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	return resp, nil
}

// 发送post不带参数 但带文件
func postFile(urlString string, m map[string]interface{}, fileMap map[string]interface{}) (map[string]interface{}, error) {
	req, err := utils.NewFormDataRequest(urlString, nil, fileMap)
	if err != nil {
		return nil, err
	}
	req.Header.Set("api_key", APIKEY)
	req.Header.Set("sign", getSign(m))
	resp, err := client.Do(req)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	response := make(map[string]interface{})
	_ = jsoniter.Unmarshal(body, &response)
	return response, nil
}

// 获取sign
func getSign(m map[string]interface{}) (result string) {
	result = hmacSha1(APISECRET, concatParam(m, APIKEY))
	return result
}

// 获取sign的两个工具方法
func hmacSha1(secret, text string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(text))
	return hex.EncodeToString(mac.Sum(nil))
}
func concatParam(m map[string]interface{}, apiKey string) string {
	result := apiKey
	keyList := make([]string, 0)
	for k, _ := range m {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, k := range keyList {
		result = result + fmt.Sprintf("%v", m[k])
	}
	return strings.Trim(result, "&")
}

type Info struct {
	ApkName     string `json:"apk_name"`
	ApkSize     string `json:"apk_size"`
	Id          int    `json:"id"`
	PackageName string `json:"package_name"`
	SrcApkMd5   string `json:"src_apk_md5"`
	StatusCode  int    `json:"status_code"`
	Version     string `json:"version"`
}

// 轮询获取加固状态码
func getStateOK(taskID int, taskType string) {
	// 请求地址
	var urlString string
	// 请求参数
	var m map[string]interface{}
	// 加固状态
	var status string

	switch taskType {
	case "android":
		urlString = fmt.Sprintf("%s/v5/protect/get_state?username=%s&apkinfo_id=%d",
			JIAGUIP, JIAGUUSERNAME, taskID)
		m = map[string]interface{}{
			"username":   JIAGUUSERNAME,
			"apkinfo_id": taskID,
		}
		break
	//case "h5":
	//	urlString = fmt.Sprintf("%s/h5/protect/get_state?username=%s&apkinfo_id=%d&h5info_id=%d",
	//		JIAGUIP, JIAGUUSERNAME, taskID, taskID)
	//	m = map[string]interface{}{
	//		"username":   JIAGUUSERNAME,
	//		"apkinfo_id": taskID,
	//		"h5info_id":  taskID,
	//	}
	//	break
	default:
		break
	}
	for {
		// 任务完成的标识量
		isFinished := true
		res, err := PostWithoutFile(urlString, m)
		ress := map[string]interface{}{}
		_ = jsoniter.Unmarshal(res, &ress)
		myRes, err := jsoniter.Marshal(ress["info"])
		if err != nil {
			log.Info(err)
		}
		var final Info
		err = jsoniter.Unmarshal(myRes, &final)
		statusCode := final.StatusCode
		// 通过查询详情接口参数进行状态查询
		if taskType == "android" {
			if statusCode == 9008 {
				status = "加固失败"
			} else if statusCode == 9009 {
				status = "加固成功"
			} else {
				isFinished = false
			}
		}
		//else if taskType == "h5" {
		//	if statusCode == 9308 {
		//		status = "加固失败"
		//	} else if statusCode == 9309 {
		//		status = "加固成功"
		//	} else {
		//		isFinished = false
		//	}
		//}

		if isFinished {
			times := time.Now()
			if err := app.DB.Model(model.JiaGuTask{}).
				Where("task_id = ?", taskID).
				Updates(model.JiaGuTask{TaskStatus: status, FinishTime: &times}).
				Error; err != nil {
				log.Info(err.Error())
			}
			break
		}

		time.Sleep(6 * time.Second)
	}

}

// 获取策略的结构体
type DataList struct {
	Code int                 `json:"code"`
	Info []model.JiaguPolicy `json:"info"`
	Msg  string              `json:"msg"`
}

// 获取策略通过类型
func GetPolicyListByType() {
	for {
		var policyList []model.JiaguPolicy
		//var idList []int
		var dataList DataList
		var urlString string
		//urlString = fmt.Sprintf("%s/v5/policy/get_list?policy_id=%d&username=%s",
		//	JIAGUIP, 1, JIAGUUSERNAME)
		urlString = fmt.Sprintf("%s/v5/policy/get_list?username=%s",
			JIAGUIP, JIAGUUSERNAME)
		m := map[string]interface{}{
			//"policy_id": 1,
			"username": JIAGUUSERNAME,
		}
		res, err := PostWithoutFile(urlString, m)
		if err != nil {
			log.Error(err.Error())
		}
		err = jsoniter.Unmarshal(res, &dataList)
		if err != nil {
			log.Error(err.Error())
		}

		// 获取所有的加固策略
		if err := app.DB.Find(&policyList).Error; err != nil {
			log.Error(err.Error)
		}

		// 如果查询出来的策略数量 不等于当前表的数据时 进行数据删除
		if len(policyList) != len(dataList.Info) {
			// 清空所有策略所有策略
			if err := app.DB.Delete(model.JiaguPolicy{}, "status = ? or status = ?", 0, 1).Error; err != nil {
				log.Error(err.Error)
			}
		}

		for _, val := range dataList.Info {
			// 进行匹配 是否有app_ids 匹配当前策略
			for _, val2 := range policyList {
				if val2.UserIds != "" && val2.Id == val.Id {
					val.UserIds = val2.UserIds
					break
				}
			}
			// 进行匹配 是否有number_of_use不为零时 赋值使用数
			for _, val2 := range policyList {
				if val2.NumberOfUse != 0 && val2.Id == val.Id {
					val.NumberOfUse = val2.NumberOfUse
					break
				}
			}

			if err := app.DB.Save(&val).Error; err != nil {
				log.Error(err.Error)
			}
		}

		log.Info("更新了加固策略")
		time.Sleep(5 * time.Minute)
	}
}
