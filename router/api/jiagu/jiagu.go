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
var JIAGUUSERNAME, APIKEY, APISECRET, JIAGUIP = JiaGuLoading()

// 请求频繁，共用连接
var client = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: false,
	},
}

// 加固内置账户
func JiaGuLoading() (string, string, string, string) {
	user := app.Conf.JiaGu
	return user.UserName, user.ApiKey, user.ApiSecret, user.Ip
}

// android部分
// 5.1上传apk加固包
func WebBoxV5Upload(c *gin.Context) {
	// 任务表
	var jiaGuTask model.JiaGuTask
	// 获取前端传入的参数
	num, appIdOK := c.GetPostForm("app_id")
	appId, _ := strconv.Atoi(num)
	num2, appTypeIDOK := c.GetPostForm("app_type_id")
	appTypeID, _ := strconv.Atoi(num2)
	num3, policyIdOK := c.GetPostForm("policy_id")
	policyReason, policyReasonOK := c.GetPostForm("policy_reason")
	apkPath, apkPathOK := c.GetPostForm("apk_path")
	channelPath, channelPathOK := c.GetPostForm("channel_path")

	if !(policyIdOK && appTypeIDOK && appIdOK && apkPathOK) {
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
	_, exist := res["status"].(float64)
	if exist {
		c.JSON(http.StatusOK, gin.H{
			"code": res["status"],
			"err":  res["msg"],
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

	// 获取当前用户id
	userId, _ := c.Get("userId")

	jiaGuTask.AppID = appId
	jiaGuTask.UserID = userId.(uint)
	jiaGuTask.PolicyID = policyId
	jiaGuTask.AppTypeID = appTypeID
	jiaGuTask.TaskStatus = "加固中"
	// 如果存在选择策略理由则添加
	if policyReasonOK {
		jiaGuTask.PolicyReason = policyReason
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

	// 进行应用表存储
	var application model.Application
	if err := app.DB.Model(model.Application{}).Where("id = ?", appId).First(&application).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	// 使用用户id拼接
	if application.UseUser != "" {
		application.UseUser = application.UseUser + "," + strconv.Itoa(int(userId.(uint)))
	} else {
		application.UseUser += strconv.Itoa(int(userId.(uint)))
	}
	// 存储使用用户id拼接
	if err := app.DB.Save(&application).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
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
	res, err := postWithoutFile(urlString, m)
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

	// 不进行加固平台删除
	//urlString := fmt.Sprintf("%s/v5/protect/delete?apkinfo_id=%d",
	//	JIAGUIP, taskID)
	//m := map[string]interface{}{
	//	"apkinfo_id": taskID,
	//}
	//res, err := postWithoutFile(urlString, m)
	//ress := map[string]interface{}{}
	//_ = jsoniter.Unmarshal(res, &ress)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": http.StatusInternalServerError,
	//		"err":  err.Error(),
	//	})
	//	return
	//} else {
	//	c.JSON(http.StatusOK, ress)
	//}
}

// 5.6获取当前用户加固任务
func WebBoxV5AllTask(c *gin.Context) {
	// 开始和结束日期
	var startTime time.Time
	var endTime time.Time
	var err error
	// 获取传入的值
	appName, appNameOK := c.GetPostForm("app_name")
	status, statusOK := c.GetPostForm("app_status")
	Time1, startTimeOK := c.GetPostForm("start_time")
	Time2, endTimeOk := c.GetPostForm("end_time")
	userName, userNameOk := c.GetPostForm("user_name")

	// 判断是否为空字符
	appNameOK = utils.IsStringEmpty(appName, appNameOK)
	statusOK = utils.IsStringEmpty(status, statusOK)
	startTimeOK = utils.IsStringEmpty(Time1, startTimeOK)
	endTimeOk = utils.IsStringEmpty(Time2, endTimeOk)
	userNameOk = utils.IsStringEmpty(userName, userNameOk)

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

	// 是否是admin用户
	isAdmin, _ := c.Get("isAdmin")
	departmentId, _ := c.Get("departmentId")
	superAdmin, _ := c.Get("superAdmin")

	// 总数
	var total int64
	var taskMessage []map[string]interface{}
	// 拼接sql语句
	var sqlString bytes.Buffer

	layout := "2006-01-02 15:04:05"

	sqlString.WriteString("jia_gu_task.deleted_at is NULL")
	if !superAdmin.(bool) {
		// 不为超级管理员 进行拼接条件
		if isAdmin.(bool) { // 为部门管理员账户则
			sqlString.WriteString(" AND user.department_id = ")
			sqlString.WriteString(strconv.Itoa(int(departmentId.(uint))))

			// 添加用户名检索功能(只有管理员能进行用户名检索)
			if userNameOk {
				sqlString.WriteString(" AND user.user_name like '")
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
	} else {
		// 添加用户名检索功能(只有管理员能进行用户名检索)
		if userNameOk {
			sqlString.WriteString(" AND user.user_name like '")
			sqlString.WriteString("%" + userName + "%")
			sqlString.WriteString("'")
		}
	}

	//  当前为android加固 所有 app_type_id 为 1
	sqlString.WriteString(" AND jia_gu_task.app_type_id = ")
	sqlString.WriteString(strconv.Itoa(1))
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
	//  添加应用名检索
	if appNameOK {
		sqlString.WriteString(" AND application.app_name like '")
		sqlString.WriteString("%" + appName + "%")
		sqlString.WriteString("'")

	}
	// 添加加固状态检索
	if statusOK {
		sqlString.WriteString(" AND jia_gu_task.task_status = ")
		sqlString.WriteString(`"`)
		sqlString.WriteString(status)
		sqlString.WriteString(`"`)
	}

	// 获取分页查询数据 | 获取总数total
	if err := app.DB.Table("jia_gu_task").
		Select("user.user_name,application.app_name,application.app_version,application.the_app,application.the_model,application.app_cn_name," +
			"jia_gu_task.task_id,jiagu_policy.name,jia_gu_task.policy_reason,jia_gu_task.created_at,jia_gu_task.task_status,jia_gu_task.finish_time," +
			"application_type.app_type").
		Joins("INNER JOIN user ON user.id = jia_gu_task.user_id").
		Joins("INNER JOIN application_type ON application_type.id = jia_gu_task.app_type_id").
		Joins("INNER JOIN application ON application.id = jia_gu_task.app_id").
		Joins("INNER JOIN jiagu_policy ON jia_gu_task.policy_id = jiagu_policy.id").
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

//// 加固系统版本号
//func WebBoxV5ReinforceVer(c *gin.Context) {
//	urlString := fmt.Sprintf("%s/v5/policy/system/qry_reinforce_ver",
//		JIAGUIP)
//	// 没有需要传入的值传入 空map
//	m := map[string]interface{}{}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}
//
//// so加固系统版本号
//func WebBoxV5SoVer(c *gin.Context) {
//	urlString := fmt.Sprintf("%s/v5/policy/system/qry_soversion",
//		JIAGUIP)
//	// 没有需要传入的值传入 空map
//	m := map[string]interface{}{}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}
//
//// 策略使用者
//func WebBoxV5PolicyUser(c *gin.Context) {
//	urlString := fmt.Sprintf("%s/v5/policy/get_policy_user?policy_id=%d",
//		JIAGUIP, 94)
//	// 没有需要传入的值传入 空map
//	m := map[string]interface{}{
//		"policy_id": 94,
//	}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}
//
//// 添加策略
//func WebBoxV5PolicyAdd(c *gin.Context) {
//	policyConfig := "{\"addition\":{\"antixml\":false,\"compresspng\":false,\"nolog\":false},\"antiTamp\":{\"enable\":false,\"polling\":true,\"ptrace\":false},\"cpu\":[\"armeabi\",\"armeabi-v7a\",\"arm64-v8a\"],\"dataEnc\":{\"enable\":false,\"dataencxor\":false,\"dataencbind\":false,\"rule\":[]},\"dex\":{\"v1\":{\"enable\":false,\"dexfast\":true,\"fun\":[]},\"v2\":{\"enable\":false,\"dexdcrp\":true,\"dexdmrp\":false,\"fakeclass\":false,\"fun\":[]},\"v4\":{\"enable\":false,\"dexbind\":false,\"allvmp\":false,\"so\":\"\",\"fun\":[]},\"obfstr\":{\"enable\":false,\"obffilter\":[]},\"antiRe\":{\"dexhunter\":false,\"antidex2jar\":true,\"antijadx\":true,\"shelldexhelper\":true}},\"integrity\":{\"enable\":false,\"rule\":[\"*\",\"!AndroidManifest.xml\"]},\"resEnc\":{\"enable\":false,\"zipres\":false,\"rule\":[]},\"runtime\":{\"pappmo\":false,\"proot\":false,\"psimulator\":false,\"antibooster\":false,\"pmem\":false,\"hooktools\":false,\"sig\":false,\"antijnject\":true,\"hijack\":{\"enable\":false,\"activity\":[]},\"safescreen\":{\"enable\":false,\"activity\":[]}},\"so\":{\"enable\":false,\"ver\":\"\",\"bind\":false,\"clear\":false,\"enc\":false,\"file\":[]},\"u3d\":{\"enable\":false,\"version\":\"\",\"unityver\":\"\",\"dll\":[]},\"version\":\"ver7.1.2_TAND_cc_211104.1.docker\"}"
//	urlString := fmt.Sprintf("%s/v5/policy/add?username=%s&policy_name=%s&policy_status=%s&policy_config=%s",
//		JIAGUIP, JIAGUUSERNAME, "中间平台通用策略1", "已启用", policyConfig)
//	// 没有需要传入的值传入 空map
//	m := map[string]interface{}{
//		"username":      JIAGUUSERNAME,
//		"policy_name":   "中间平台通用策略1",
//		"policy_status": "已启用",
//	}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}
//
//// 修改策略
//func WebBoxV5PolicyModify(c *gin.Context) {
//	policyConfig := "{\"addition\":{\"antixml\":false,\"compresspng\":false,\"nolog\":false},\"antiTamp\":{\"enable\":false,\"polling\":true,\"ptrace\":false},\"cpu\":[\"armeabi\",\"armeabi-v7a\",\"arm64-v8a\"],\"dataEnc\":{\"enable\":false,\"dataencxor\":false,\"dataencbind\":false,\"rule\":[]},\"dex\":{\"v1\":{\"enable\":false,\"dexfast\":true,\"fun\":[]},\"v2\":{\"enable\":false,\"dexdcrp\":true,\"dexdmrp\":false,\"fakeclass\":false,\"fun\":[]},\"v4\":{\"enable\":false,\"dexbind\":false,\"allvmp\":false,\"so\":\"\",\"fun\":[]},\"obfstr\":{\"enable\":false,\"obffilter\":[]},\"antiRe\":{\"dexhunter\":false,\"antidex2jar\":true,\"antijadx\":true,\"shelldexhelper\":true}},\"integrity\":{\"enable\":false,\"rule\":[\"*\",\"!AndroidManifest.xml\"]},\"resEnc\":{\"enable\":false,\"zipres\":false,\"rule\":[]},\"runtime\":{\"pappmo\":false,\"proot\":false,\"psimulator\":false,\"antibooster\":false,\"pmem\":false,\"hooktools\":false,\"sig\":false,\"antijnject\":true,\"hijack\":{\"enable\":false,\"activity\":[]},\"safescreen\":{\"enable\":false,\"activity\":[]}},\"so\":{\"enable\":false,\"ver\":\"\",\"bind\":false,\"clear\":false,\"enc\":false,\"file\":[]},\"u3d\":{\"enable\":false,\"version\":\"\",\"unityver\":\"\",\"dll\":[]},\"version\":\"ver7.1.2_TAND_cc_211104.1.docker\"}"
//	urlString := fmt.Sprintf("%s/v5/policy/modify?policy_name=%s&policy_status=%s&policy_config=%s&policy_id=%d",
//		JIAGUIP, "中间平台通用策略2", "已启用", policyConfig, 1483)
//	// 没有需要传入的值传入 空map
//	m := map[string]interface{}{
//		"policy_name":   "中间平台通用策略2",
//		"policy_status": "已启用",
//		"policy_id":     1483,
//	}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}
//
//// 删除策略
//func WebBoxV5PolicyDelete(c *gin.Context) {
//	urlString := fmt.Sprintf("%s/v5/policy/delete?policy_id=%d",
//		JIAGUIP, 1483)
//	// 没有需要传入的值传入 空map
//	m := map[string]interface{}{
//		"policy_id": 1483,
//	}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}

// html5部分
// 6.1上传H5加固包
func WebBoxH5Upload(c *gin.Context) {
	// 任务表
	var jiaGuTask model.JiaGuTask
	// 获取前端传入的参数
	num, appIdOK := c.GetPostForm("app_id")
	appId, _ := strconv.Atoi(num)
	num2, appTypeIDOK := c.GetPostForm("app_type_id")
	appTypeID, _ := strconv.Atoi(num2)

	num3, policyIdOK := c.GetPostForm("policy_id")
	policyReason, policyReasonOK := c.GetPostForm("policy_reason")
	h5Path, h5PathOk := c.GetPostForm("h5_path")

	if !(policyIdOK && appTypeIDOK && appIdOK && h5PathOk) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	policyId, err := strconv.Atoi(num3)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	// 向加固平台发送请求
	urlString := fmt.Sprintf("%s/h5/protect/upload?username=%s&policy_id=%d&upload_type=%d",
		JIAGUIP, JIAGUUSERNAME, policyId, 2)
	m := map[string]interface{}{
		"username":    JIAGUUSERNAME,
		"policy_id":   policyId,
		"upload_type": 2,
	}

	file, err := os.Open(h5Path)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  "获取文件失败",
		})
		return
	}
	fileMap := map[string]interface{}{
		"h5_file": file,
	}
	res, err := postFile(urlString, m, fileMap)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
	}
	_, exist := res["status"].(float64)
	if exist {
		c.JSON(http.StatusOK, gin.H{
			"code": res["status"],
			"err":  res["msg"],
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

	// 获取用户id
	userId, _ := c.Get("userId")

	jiaGuTask.AppID = appId
	jiaGuTask.UserID = userId.(uint)
	jiaGuTask.PolicyID = policyId
	jiaGuTask.AppTypeID = appTypeID
	jiaGuTask.TaskStatus = "加固中"
	// 如果存在选择策略理由则添加
	if policyReasonOK {
		jiaGuTask.PolicyReason = policyReason
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

	// 进行应用表存储
	var application model.Application
	if err := app.DB.Model(model.Application{}).Where("id = ?", appId).First(&application).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	// 使用用户id拼接
	if application.UseUser != "" {
		application.UseUser = application.UseUser + "," + strconv.Itoa(int(userId.(uint)))
	} else {
		application.UseUser += strconv.Itoa(int(userId.(uint)))
	}
	// 存储使用用户
	if err := app.DB.Save(&application).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	// 轮询获取最新加固状态
	go getStateOK(jiaGuTask.TaskID, "h5")
}

// 6.2查询H5加固状态
func WebBoxH5GetState(c *gin.Context) {
	// 任务id
	taskID, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	urlString := fmt.Sprintf("%s/h5/protect/get_state?username=%s&h5info_id=%d",
		JIAGUIP, JIAGUUSERNAME, taskID)
	m := map[string]interface{}{
		"username":  JIAGUUSERNAME,
		"h5info_id": taskID,
	}
	res, err := postWithoutFile(urlString, m)
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

// 6.3下载H5加固包
func WebBoxH5Download(c *gin.Context) {
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
	urlString := fmt.Sprintf("%s/h5/protect/download?username=%s&h5info_id=%d&download_type=%d&apkinfo_id=%d",
		JIAGUIP, JIAGUUSERNAME, taskID, downloadType, taskID)
	m := map[string]interface{}{
		"username":      JIAGUUSERNAME,
		"h5info_id":     taskID,
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

// 6.4下载H5加固日志
func WebBoxH5DownloadLog(c *gin.Context) {
	// 获取任务id
	num, taskIDOK := c.GetQuery("id")
	if !taskIDOK {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
	}
	taskID, _ := strconv.Atoi(num)

	urlString := fmt.Sprintf("%s/h5/protect/download_log?username=%s&h5info_id=%d&apkinfo_id=%d",
		JIAGUIP, JIAGUUSERNAME, taskID, taskID)
	m := map[string]interface{}{
		"username":   JIAGUUSERNAME,
		"apkinfo_id": taskID,
		"h5info_id":  taskID,
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

// 6.5删除H5加固记录
func WebBoxH5Delete(c *gin.Context) {
	// 任务id
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
	//urlString := fmt.Sprintf("%s/h5/protect/delete?h5info_id=%d",
	//	JIAGUIP, taskID)
	//m := map[string]interface{}{
	//	"h5info_id": taskID,
	//}
	//res, err := postWithoutFile(urlString, m)
	//ress := map[string]interface{}{}
	//_ = jsoniter.Unmarshal(res, &ress)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": http.StatusInternalServerError,
	//		"err":  err.Error(),
	//	})
	//} else {
	//	c.JSON(http.StatusOK, ress)
	//}
}

// 6.6 h5获取用户下的所有任务
func WebBoxH5TaskAllTask(c *gin.Context) {
	// 开始和结束日期
	var startTime time.Time
	var endTime time.Time
	var err error
	// 获取传入的值
	appName, appNameOK := c.GetPostForm("app_name")
	status, statusOK := c.GetPostForm("app_status")
	Time1, startTimeOK := c.GetPostForm("start_time")
	Time2, endTimeOk := c.GetPostForm("end_time")
	userName, userNameOk := c.GetPostForm("user_name")

	// 判断是否为空字符
	appNameOK = utils.IsStringEmpty(appName, appNameOK)
	statusOK = utils.IsStringEmpty(status, statusOK)
	startTimeOK = utils.IsStringEmpty(Time1, startTimeOK)
	endTimeOk = utils.IsStringEmpty(Time2, endTimeOk)
	userNameOk = utils.IsStringEmpty(userName, userNameOk)
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
	// 将page size offNum封装成工具方法
	page, size, offNum := utils.GetFromDataPageSizeOffNum(c)
	if page+size+offNum == -3 {
		return
	}

	// 是否是admin用户
	isAdmin, _ := c.Get("isAdmin")
	departmentId, _ := c.Get("departmentId")
	superAdmin, _ := c.Get("superAdmin")

	// 拼接sql语句
	var sqlString bytes.Buffer

	layout := "2006-01-02 15:04:05"
	// 总数
	var total int64
	var taskMessage []map[string]interface{}

	sqlString.WriteString("jia_gu_task.deleted_at is NULL")
	// 不为超级管理员进行条件拼接
	if !superAdmin.(bool) {
		if isAdmin.(bool) { // 为部门管理员账户则
			sqlString.WriteString(" AND user.department_id = ")
			sqlString.WriteString(strconv.Itoa(int(departmentId.(uint))))

			// 添加用户名检索功能(只有管理员能进行用户名检索)
			if userNameOk {
				sqlString.WriteString(" AND user.user_name like '")
				sqlString.WriteString("%" + userName + "%")
				sqlString.WriteString("'")
			}
		} else { // 为普通账户则
			// 获取当前用户id
			num3, _ := c.Get("userId")
			userId := int(num3.(uint))
			sqlString.WriteString(" AND jia_gu_task.user_id = ")
			sqlString.WriteString(strconv.Itoa(userId))
		}
	} else {
		// 添加用户名检索功能(只有管理员能进行用户名检索)
		if userNameOk {
			sqlString.WriteString(" AND user.user_name like '")
			sqlString.WriteString("%" + userName + "%")
			sqlString.WriteString("'")
		}
	}

	//  当前为h5加固 所有 app_type_id 为 2
	sqlString.WriteString(" AND jia_gu_task.app_type_id = ")
	sqlString.WriteString(strconv.Itoa(2))
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
	//  添加应用名检索
	if appNameOK {
		sqlString.WriteString(" AND application.app_name like '")
		sqlString.WriteString("%" + appName + "%")
		sqlString.WriteString("'")

	}
	// 添加加固状态检索
	if statusOK {
		sqlString.WriteString(" AND jia_gu_task.task_status = ")
		sqlString.WriteString(`"`)
		sqlString.WriteString(status)
		sqlString.WriteString(`"`)
	}

	// 获取分页查询数据 | 获取总数total
	if err := app.DB.Table("jia_gu_task").
		Select("user.user_name,application.app_name,application.app_version,application.the_app,application.the_model,application.app_cn_name," +
			"jia_gu_task.task_id,jiagu_policy.name,jia_gu_task.policy_reason,jia_gu_task.created_at,jia_gu_task.task_status,jia_gu_task.finish_time," +
			"application_type.app_type").
		Joins("INNER JOIN user ON user.id = jia_gu_task.user_id").
		Joins("INNER JOIN application_type ON application_type.id = jia_gu_task.app_type_id").
		Joins("INNER JOIN application ON application.id = jia_gu_task.app_id").
		Joins("INNER JOIN jiagu_policy ON jia_gu_task.policy_id = jiagu_policy.id").
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

//// h5加固系统版本号
//func WebBoxH5Ver(c *gin.Context) {
//	urlString := fmt.Sprintf("%s/h5/policy/system/qry_reinforce_ver",
//		JIAGUIP)
//	m := map[string]interface{}{}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}
//
//// h5添加策略
//func WebBoxH5PolicyAdd(c *gin.Context) {
//	policyConfig := "{\"bindApp\":{\"app_list\":[],\"enable\":false,\"level\":\"0.3\"},\"compact\":{\"enable\":true,\"selfDefending\":{\"enable\":false}},\"controlFlowFlattening\":{\"controlFlowFlatteningThreshold\":\"0.5\",\"enable\":true},\"convDefine\":{\"enable\":true},\"deadCodeInjection\":{\"deadCodeInjectionThreshold\":\"0.5\",\"enable\":true},\"debugProtection\":{\"debugProtectionInterval\":{\"enable\":false,\"entry\":\"\"},\"enable\":true,\"rule\":[\"*.js\",\"!*-min.js\",\"!jJqQuery*.js\",\"!qrcode.js\",\"!cordova*.js\",\"!ionic.js\"]},\"disableConsoleOutput\":false,\"doPacked\":{\"enable\":true},\"domainLock\":{\"domain\":[],\"enable\":false},\"encode\":{\"enable\":true,\"encode\":\"UTF-8\"},\"html\":{\"enable\":true,\"rule\":[\"*.html\"]},\"identifierNamesGenerator\":\"hexadecimal\",\"imageCompress\":{\"enable\":true},\"js\":{\"rule\":[\"**\",\"!*.min.js\",\"!*-min.js\",\"!jJqQuery*.js\",\"!qrcode.js\",\"!cordova*.js\",\"!ionic*.js\"],\"enable\":true},\"numbersToExpressions\":{\"enable\":true},\"renameGlobals\":{\"enable\":false},\"renameProperties\":{\"enable\":false},\"reservedFunctions\":{\"enable\":true,\"names\":[]},\"reservedNames\":{\"enable\":true,\"names\":[]},\"seed\":\"0\",\"simplify\":{\"enable\":true},\"sourceMap\":{\"enable\":false,\"sourceMapBaseUrl\":\"\",\"sourceMapFileName\":\"\",\"sourceMapMode\":\"inline\"},\"splitStrings\":{\"enable\":false,\"splitStringsChunkLength\":\"2\"},\"stringArray\":{\"enable\":true,\"rotateStringArray\":false,\"stringArrayEncoding\":\"base64\",\"stringArrayThreshold\":\"0.5\"},\"transformObjectKeys\":{\"enable\":false},\"unicodeEscapeSequence\":true,\"version\":\"ver3.3.20210324_h5\",\"vmp\":{\"enable\":true,\"rule\":[\"**\",\"!*.min.js\",\"!*-min.js\",\"!jJqQuery*.js\",\"!qrcode.js\",\"!cordova*.js\",\"!ionic*.js\"]}}"
//	urlString := fmt.Sprintf("%s/h5/policy/add?username=%s&policy_name=%s&policy_status=%s&policy_config=%s",
//		JIAGUIP, JIAGUUSERNAME, "中间平台通用策略", "已启用", policyConfig)
//	m := map[string]interface{}{
//		"username":      JIAGUUSERNAME,
//		"policy_name":   "中间平台通用策略",
//		"policy_status": "已启用",
//	}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}
//
//// h5修改策略
//func WebBoxH5PolicyModify(c *gin.Context) {
//	policyConfig := "{\"bindApp\":{\"app_list\":[],\"enable\":false,\"level\":\"0.3\"},\"compact\":{\"enable\":true,\"selfDefending\":{\"enable\":false}},\"controlFlowFlattening\":{\"controlFlowFlatteningThreshold\":\"0.5\",\"enable\":true},\"convDefine\":{\"enable\":true},\"deadCodeInjection\":{\"deadCodeInjectionThreshold\":\"0.5\",\"enable\":true},\"debugProtection\":{\"debugProtectionInterval\":{\"enable\":false,\"entry\":\"\"},\"enable\":true,\"rule\":[\"*.js\",\"!*-min.js\",\"!jJqQuery*.js\",\"!qrcode.js\",\"!cordova*.js\",\"!ionic.js\"]},\"disableConsoleOutput\":false,\"doPacked\":{\"enable\":true},\"domainLock\":{\"domain\":[],\"enable\":false},\"encode\":{\"enable\":true,\"encode\":\"UTF-8\"},\"html\":{\"enable\":true,\"rule\":[\"*.html\"]},\"identifierNamesGenerator\":\"hexadecimal\",\"imageCompress\":{\"enable\":true},\"js\":{\"rule\":[\"**\",\"!*.min.js\",\"!*-min.js\",\"!jJqQuery*.js\",\"!qrcode.js\",\"!cordova*.js\",\"!ionic*.js\"],\"enable\":true},\"numbersToExpressions\":{\"enable\":true},\"renameGlobals\":{\"enable\":false},\"renameProperties\":{\"enable\":false},\"reservedFunctions\":{\"enable\":true,\"names\":[]},\"reservedNames\":{\"enable\":true,\"names\":[]},\"seed\":\"0\",\"simplify\":{\"enable\":true},\"sourceMap\":{\"enable\":false,\"sourceMapBaseUrl\":\"\",\"sourceMapFileName\":\"\",\"sourceMapMode\":\"inline\"},\"splitStrings\":{\"enable\":false,\"splitStringsChunkLength\":\"2\"},\"stringArray\":{\"enable\":true,\"rotateStringArray\":false,\"stringArrayEncoding\":\"base64\",\"stringArrayThreshold\":\"0.5\"},\"transformObjectKeys\":{\"enable\":false},\"unicodeEscapeSequence\":true,\"version\":\"ver3.3.20210324_h5\",\"vmp\":{\"enable\":true,\"rule\":[\"**\",\"!*.min.js\",\"!*-min.js\",\"!jJqQuery*.js\",\"!qrcode.js\",\"!cordova*.js\",\"!ionic*.js\"]}}"
//	urlString := fmt.Sprintf("%s/h5/policy/modify?policy_id=%d&policy_name=%s&policy_status=%s&policy_config=%s",
//		JIAGUIP, 1485, "中间平台通用策略", "已启用", policyConfig)
//	m := map[string]interface{}{
//		"policy_id":     1485,
//		"policy_name":   "中间平台通用策略",
//		"policy_status": "已启用",
//	}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}
//
//// h5删除策略
//func WebBoxH5PolicyDelete(c *gin.Context) {
//	urlString := fmt.Sprintf("%s/h5/policy/delete?policy_id=%d",
//		JIAGUIP, 1480)
//	m := map[string]interface{}{
//		"policy_id": 1480,
//	}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}
//
//// h5查询策略使用用户
//func WebBoxH5PolicyUser(c *gin.Context) {
//	urlString := fmt.Sprintf("%s/h5/policy/get_policy_user?policy_id=%d",
//		JIAGUIP, 1485)
//	m := map[string]interface{}{
//		"policy_id": 1485,
//	}
//	res, err := postWithoutFile(urlString, m)
//	ress := map[string]interface{}{}
//	_ = jsoniter.Unmarshal(res, &ress)
//	if err != nil {
//		c.JSON(http.StatusOK, gin.H{
//			"code": http.StatusInternalServerError,
//			"err":  err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, ress)
//	}
//}

// 加固策略部分
// 13.1获取加固策略
func JiaguPolicyFind(c *gin.Context) {
	// 策略的原数据
	var datas []model.JiaguPolicy
	var sqlString bytes.Buffer
	jiaguType, jiaguTypeOK := c.GetPostForm("app_type_id")
	policyId, policyIdOK := c.GetPostForm("policy_id")
	// 判断是否为空
	jiaguTypeOK = utils.IsStringEmpty(jiaguType, jiaguTypeOK)
	policyIdOK = utils.IsStringEmpty(policyId, policyIdOK)

	// 拼接sql语句
	sqlString.WriteString(`status = "已启用" `)
	// 有加固类型时
	if jiaguTypeOK {
		var appType string
		if err := app.DB.Select("app_type").Model(model.ApplicationType{}).Where("id = ?", jiaguType).First(&appType).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
		sqlString.WriteString(" AND type = ")
		sqlString.WriteString(`"`)
		sqlString.WriteString(appType)
		sqlString.WriteString(`"`)
	}

	// 有策略id时
	if policyIdOK {
		sqlString.WriteString(" AND id = ")
		sqlString.WriteString(policyId)
	}
	//  查出id对应的数据
	if err := app.DB.Model(model.JiaguPolicy{}).Where(sqlString.String()).Find(&datas).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": datas,
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
	jiaguType, jiaguTypeOK := c.GetPostForm("app_type_id")
	// 分页
	page, size, offNum := utils.GetFromDataPageSizeOffNum(c)
	if page+size+offNum == -3 {
		return
	}
	// 判断是否为空
	jiaguTypeOK = utils.IsStringEmpty(jiaguType, jiaguTypeOK)

	// 拼接sql语句
	sqlString.WriteString(`status = "已启用" `)
	// 有加固类型时
	if jiaguTypeOK {
		var appType string
		if err := app.DB.Select("app_type").Model(model.ApplicationType{}).Where("id = ?", jiaguType).First(&appType).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
		sqlString.WriteString(" AND type = ")
		sqlString.WriteString(`"`)
		sqlString.WriteString(appType)
		sqlString.WriteString(`"`)
	}

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

// post不带参数和文件
func postWithoutFile(urlString string, m map[string]interface{}) ([]byte, error) {
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
	case "h5":
		urlString = fmt.Sprintf("%s/h5/protect/get_state?username=%s&apkinfo_id=%d&h5info_id=%d",
			JIAGUIP, JIAGUUSERNAME, taskID, taskID)
		m = map[string]interface{}{
			"username":   JIAGUUSERNAME,
			"apkinfo_id": taskID,
			"h5info_id":  taskID,
		}
		break
	default:
		break
	}
	for {
		// 任务完成的标识量
		isFinished := true
		res, err := postWithoutFile(urlString, m)
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
		} else if taskType == "h5" {
			if statusCode == 9308 {
				status = "加固失败"
			} else if statusCode == 9309 {
				status = "加固成功"
			} else {
				isFinished = false
			}
		}

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
		//if types == "android" {
		urlString = fmt.Sprintf("%s/v5/policy/get_list?policy_id=%d&username=%s",
			JIAGUIP, 1, JIAGUUSERNAME)
		//}
		//else if types == "h5" {
		//	urlString = fmt.Sprintf("%s/h5/policy/get_list?policy_id=%d&username=%s",
		//		JIAGUIP, 1, JIAGUUSERNAME)
		//}
		// 没有需要传入的值传入 空map
		m := map[string]interface{}{
			"policy_id": 1,
			"username":  JIAGUUSERNAME,
		}
		res, _ := postWithoutFile(urlString, m)
		_ = jsoniter.Unmarshal(res, &dataList)

		// 获取所有的加固策略
		if err := app.DB.Find(&policyList).Error; err != nil {
			log.Error(err.Error)
		}

		// 如果查询出来的策略数量 不等于当前表的数据时 进行数据删除
		if len(policyList) != len(dataList.Info) {
			// 清空所有策略所有策略
			if err := app.DB.Delete(model.JiaguPolicy{}, "status like ? ", "%已%").Error; err != nil {
				log.Error(err.Error)
			}
		}

		for _, val := range dataList.Info {
			// 拼接加固策略类型
			if val.LicKeyHelper == "AppShield" {
				val.Type = "android加固"
			} else if val.LicKeyHelper == "H5Shield" {
				val.Type = "h5加固"
			}
			if err := app.DB.Save(&val).Error; err != nil {
				log.Error(err.Error)
			}
		}

		log.Info("更新了加固策略")
		time.Sleep(5 * time.Minute)
	}
}
