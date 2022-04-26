package jiagu

import (
	"bytes"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 7.1添加系统
func CreateApplication(c *gin.Context) {
	// 检测必传参数是否赋值
	appTypeID, appTypeOK := c.GetPostForm("app_type_id")
	appName, appNameOK := c.GetPostForm("app_name")
	appVer, appVerOK := c.GetPostForm("app_version")
	appCnName, appCnNameOK := c.GetPostForm("app_cn_name")
	modelName, modelNameOK := c.GetPostForm("model_name")
	modelCnName, modelCnNameOK := c.GetPostForm("model_cn_name")
	policyId, policyIdOK := c.GetPostForm("policy_id")

	appTypeOK = utils.IsStringEmpty(appTypeID, appTypeOK)
	appNameOK = utils.IsStringEmpty(appName, appNameOK)
	appVerOK = utils.IsStringEmpty(appVer, appVerOK)
	appCnNameOK = utils.IsStringEmpty(appCnName, appCnNameOK)
	modelNameOK = utils.IsStringEmpty(modelName, modelNameOK)
	modelCnNameOK = utils.IsStringEmpty(modelCnName, modelCnNameOK)
	policyIdOK = utils.IsStringEmpty(policyId, policyIdOK)

	if !(appTypeOK && appNameOK && appVerOK && appCnNameOK && modelNameOK && modelCnNameOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	// 给应用表赋值
	var application model.Application
	application.AppName = appName
	application.AppVersion = appVer
	application.AppTypeID, _ = strconv.Atoi(appTypeID)
	application.ModelName = modelName
	application.AppCnName = appCnName
	application.ModelCnName = modelCnName
	num2, _ := strconv.Atoi(policyId)
	application.RecommendPolicy = num2
	application.LastChangeTime = time.Now().Format("2006-01-02 15:04:05")
	application.TheApp = appName + "-" + appCnName
	application.TheModel = modelName + "-" + modelCnName
	if err := app.DB.Model(model.Application{}).Create(&application).Error; err != nil {
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

// 7.2编辑系统
func EditApplication(c *gin.Context) {
	// 检测必传参数是否赋值
	appID, appIdOK := c.GetPostForm("app_id")
	appTypeID, appTypeOK := c.GetPostForm("app_type_id")
	appName, appNameOK := c.GetPostForm("app_name")
	appVer, appVerOK := c.GetPostForm("app_version")
	appCnName, appCnNameOK := c.GetPostForm("app_cn_name")
	modelName, modelNameOK := c.GetPostForm("model_name")
	modelCnName, modelCnNameOK := c.GetPostForm("model_cn_name")
	policyId, policyIdOK := c.GetPostForm("policy_id")

	appIdOK = utils.IsStringEmpty(appID, appIdOK)
	appTypeOK = utils.IsStringEmpty(appTypeID, appTypeOK)
	appNameOK = utils.IsStringEmpty(appName, appNameOK)
	appVerOK = utils.IsStringEmpty(appVer, appVerOK)
	appCnNameOK = utils.IsStringEmpty(appCnName, appCnNameOK)
	modelNameOK = utils.IsStringEmpty(modelName, modelNameOK)
	modelCnNameOK = utils.IsStringEmpty(modelCnName, modelCnNameOK)
	policyIdOK = utils.IsStringEmpty(policyId, policyIdOK)

	if !(appIdOK && appTypeOK && appNameOK && appVerOK && appCnNameOK && modelNameOK && modelCnNameOK && policyIdOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	// 给应用表赋值
	var application model.Application
	num, _ := strconv.Atoi(appID)
	application.ID = uint(num)
	application.AppName = appName
	application.AppVersion = appVer
	application.AppTypeID, _ = strconv.Atoi(appTypeID)
	application.ModelName = modelName
	application.AppCnName = appCnName
	application.ModelCnName = modelCnName
	num2, _ := strconv.Atoi(policyId)
	application.RecommendPolicy = num2
	application.LastChangeTime = time.Now().Format("2006-01-02 15:04:05")
	application.TheApp = appName + "-" + appCnName
	application.TheModel = modelName + "-" + modelCnName
	if err := app.DB.Where("id = ?", application.ID).Updates(&application).
		Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  model.ModSuccess,
	})
}

// 7.3删除系统
func DelApplication(c *gin.Context) {
	// 检测必传参数是否赋值
	appID, appIDOK := c.GetPostForm("app_id")
	if !appIDOK {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}

	// 给应用表赋值
	var application model.Application
	num, _ := strconv.Atoi(appID)
	application.ID = uint(num)
	if err := app.DB.Model(model.Application{}).First(&application).
		Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	// 判断是否有用户关联
	if application.UseUser != "" {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  "删除失败,该应用已存在关联信息",
		})
		return
	}

	// 进行删除操作
	if err := app.DB.Model(model.Application{}).Delete(&application).
		Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  model.DelSuccess,
	})
}

// 7.4查询系统
func GetApplication(c *gin.Context) {
	var applist = new([]map[string]interface{})
	// 判断非必传值 查询条件是否存在
	num, appTypeIdOK := c.GetPostForm("app_type_id")
	appTypeId, _ := strconv.Atoi(num)
	appName, appNameOK := c.GetPostForm("app_name")
	appVer, appVerOK := c.GetPostForm("app_version")
	modelName, modelNameOK := c.GetPostForm("model_name")
	TheApp, TheAppOK := c.GetPostForm("the_app")
	TheModel, TheModelOK := c.GetPostForm("the_model")

	// 判断值是否为空
	appTypeIdOK = utils.IsStringEmpty(num, appTypeIdOK)
	appNameOK = utils.IsStringEmpty(appName, appNameOK)
	appVerOK = utils.IsStringEmpty(appVer, appVerOK)
	modelNameOK = utils.IsStringEmpty(modelName, modelNameOK)
	TheAppOK = utils.IsStringEmpty(TheApp, TheAppOK)
	TheModelOK = utils.IsStringEmpty(TheModel, TheModelOK)

	// 将page size offNum封装成工具方法
	page, size, offNum := utils.GetFromDataPageSizeOffNum(c)
	if page+size+offNum == -3 {
		return
	}
	var total int64
	var sqlString bytes.Buffer
	// 查询条件存在与否 导致 sql的where条件变更
	sqlString.WriteString("application.deleted_at is NULL")
	if appTypeIdOK {
		sqlString.WriteString(" AND application.app_type_id = ")
		sqlString.WriteString(strconv.Itoa(appTypeId))
	}
	if appNameOK {
		sqlString.WriteString(" AND application.app_name like '")
		sqlString.WriteString("%" + appName + "%")
		sqlString.WriteString("'")
	}
	if appVerOK {
		sqlString.WriteString(" AND application.app_version like '")
		sqlString.WriteString("%" + appVer + "%")
		sqlString.WriteString("'")
	}
	if modelNameOK {
		sqlString.WriteString(" AND application.model_name like '")
		sqlString.WriteString("%" + modelName + "%")
		sqlString.WriteString("'")
	}
	if TheAppOK {
		sqlString.WriteString(" AND application.the_app = ")
		sqlString.WriteString(`"` + TheApp + `"`)
	}
	if TheModelOK {
		sqlString.WriteString(" AND application.the_model = ")
		sqlString.WriteString(`"` + TheModel + `"`)
	}

	if err := app.DB.Table("application").
		Select("application_type.id as app_type_id,application.id,application.app_name," +
			"application_type.app_type,application.app_version,application.last_change_time," +
			"application.model_name,application.app_cn_name,application.model_cn_name,application.recommend_policy," +
			"application.the_app,application.the_model," +
			"jiagu_policy.name as recommend_policy_name").
		Joins("LEFT JOIN application_type ON application.app_type_id = application_type.id").
		Joins("LEFT JOIN jiagu_policy ON application.recommend_policy= jiagu_policy.id").
		Where(sqlString.String()).
		Count(&total).
		Offset(offNum).
		Limit(size).
		Order("created_at desc").
		Scan(applist).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": applist,
			"total":    total,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})

}

// 7.5获取所有的系统版本信息
func GetApplicationType(c *gin.Context) {
	var appTypeList = new([]model.ApplicationType)
	if err := app.DB.Model(model.ApplicationType{}).Find(&appTypeList).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": appTypeList,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})
}

// 7.6查询系统通过id
func GetApplicationShow(c *gin.Context) {
	type appData struct {
		ID              int    `json:"id"`
		AppName         string `json:"app_name"`         // 系统编号
		ModelName       string `json:"model_name"`       // 模板编号
		AppCnName       string `json:"app_cn_name"`      // 系统中文全称
		AppVersion      string `json:"app_version"`      // 系统英文简称
		ModelCnName     string `json:"model_cn_name"`    // 模板中文全称
		RecommendPolicy int    `json:"recommend_policy"` // 推荐策略
		AppTypeID       int    `json:"app_type_id"`      // 系统类型编号
		LastChangeTime  string `json:"last_change_time"` // 录入时间
		TheApp          string `json:"the_app"`          // 系统编号 + 系统中文全称
		TheModel        string `json:"the_model"`        //  模板编号 + 模板中文全称
	}
	var data []appData
	// 判断非必传值 查询条件是否存在
	appId, appIdOK := c.GetPostForm("app_id")
	// 加固应该类型
	appTypeId, appTypeIdOk := c.GetPostForm("app_type_id")
	// 系统
	TheApp, TheAppOK := c.GetPostForm("the_app")
	// 判断值是否为空
	appIdOK = utils.IsStringEmpty(appId, appIdOK)
	appTypeIdOk = utils.IsStringEmpty(appTypeId, appTypeIdOk)
	TheAppOK = utils.IsStringEmpty(TheApp, TheAppOK)
	if !appTypeIdOk {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}
	// 拼接sql
	var sqlString bytes.Buffer
	sqlString.WriteString("app_type_id = ")
	sqlString.WriteString(appTypeId)

	if appIdOK {
		sqlString.WriteString(" AND id = ")
		sqlString.WriteString(appId)

	}
	if TheAppOK {
		sqlString.WriteString(" AND application.the_app = ")
		sqlString.WriteString(`"` + TheApp + `"`)
	}

	if err := app.DB.Model(model.Application{}).
		Where(sqlString.String()).
		Find(&data).
		Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"info": gin.H{
			"datalist": data,
		},
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})

}

// 7.7通过特定字符获取对应的数据
func SearchList(c *gin.Context) {
	list := make([]string, 0)
	var sqlString bytes.Buffer
	keyName, keyNameOk := c.GetPostForm("key_name")
	keyVal, keyValOk := c.GetPostForm("key_val")

	keyNameOk = utils.IsStringEmpty(keyName, keyNameOk)
	keyValOk = utils.IsStringEmpty(keyName, keyValOk)

	if !(keyNameOk && keyValOk) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}

	sqlString.WriteString(keyName + " like ")
	sqlString.WriteString("'%" + keyVal + "%'")

	if err := app.DB.Model(model.Application{}).
		Distinct(keyName).
		Where(sqlString.String()).
		Find(&list).Error; err != nil {
		log.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": gin.H{
			"datalist": list,
		},
		"msg": model.ReqSuccess,
	})
}

// 7.8集中修改策略
func ModifyRecommendPolicy(c *gin.Context) {
	ids, idsOK := c.GetPostForm("ids")
	num, policyIdOK := c.GetPostForm("recommend_policy")
	policyId, _ := strconv.Atoi(num)
	num2, AppTypeOK := c.GetPostForm("app_type_id")
	AppType, _ := strconv.Atoi(num2)

	idsOK = utils.IsStringEmpty(ids, idsOK)
	policyIdOK = utils.IsStringEmpty(num, policyIdOK)
	AppTypeOK = utils.IsStringEmpty(num2, AppTypeOK)

	if !(idsOK && policyIdOK && AppTypeOK) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  model.ReqParameterMissing,
		})
		return
	}

	ids = strings.ReplaceAll(ids, "[", "")
	ids = strings.ReplaceAll(ids, "]", "")
	strArr := strings.Split(ids, ",")
	intArr := make([]int, len(strArr))
	for _, val := range strArr {
		iVal, _ := strconv.Atoi(val)
		intArr = append(intArr, iVal)
	}

	// 遍历id 更新策略id
	for _, iVal := range intArr {
		if err := app.DB.Model(model.Application{}).Where("id = ?", iVal).
			Updates(model.Application{RecommendPolicy: policyId, AppTypeID: AppType}).Error; err != nil {
			log.Error(err.Error())
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"err":  err.Error(),
			})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  model.ReqSuccess,
	})
}
