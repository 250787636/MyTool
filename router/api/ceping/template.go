package ceping

import (
	"errors"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func Check(reponse map[string]interface{}, post []byte, handlerFunc gin.HandlerFunc, c *gin.Context) error {
	errMessage := ""
	if len(reponse) == 0 {
		errMessage = strings.Trim(string(post), `"`)
	} else if key, ok := reponse["state"].(float64); ok && key != 200 {
		errMessage = reponse["msg"].(string)
	}

	if errMessage == "签名验证失败" || errMessage == "token验证失败" {
		// 1.尝试是否可以获取到token
		_, _, err := app.GetCpToken(app.Conf.CePing.UserName, app.Conf.CePing.Password, app.Conf.CePing.Ip)
		if err != nil {
			// 如果获取不到就返回错误
			response.FailWithMessage("token获取失败，请检查配置", c)
			return errors.New("token获取失败，请检查配置")
		}
		// 2.获取到token便重新调用该方法
		app.Conf = app.LoadConfig()
		handlerFunc(c)
		return errors.New(errMessage)
	}
	if errMessage != "" {
		log.Error("err", errMessage)
		response.FailWithMessage("调用测评接口失败，错误信息:"+errMessage, c)
		return errors.New(errMessage)
	}
	return nil
}

type AddTemplateRequest struct {
	TemplateType   string `form:"template_type" binding:"required"` //模板类型 安卓ad 苹果 ios 小程序 mp
	TemplateName   string `form:"template_name" binding:"required"`
	ItemKeys       string `form:"item_keys" binding:"required"`
	IsOWASP        bool   `form:"is_owasp"`        // 是否是OWASP模板 true 是
	ReportLanguage string `form:"report_language"` //导出模板语言
}

// AddTemplate 新增模版
func AddTemplate(c *gin.Context) {

	req := AddTemplateRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}
	var tem model.Template
	if err := app.DB.Model(&model.Template{}).Where("template_name = ?", req.TemplateName).Where("template_type = ?", req.TemplateType).Find(&tem).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if tem.ID > 0 {
		response.FailWithMessage("该模板名称已存在", c)
		return
	}

	str1 := strings.ReplaceAll(req.ItemKeys, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")

	itemKeysArray := "["

	for key, value := range idArray {
		if key != len(idArray)-1 {
			itemKeysArray += "\"" + value + "\","

		} else {
			itemKeysArray += "\"" + value + "\""
		}
	}
	itemKeysArray += "]"
	fmt.Println("itemKeysArray", itemKeysArray)

	id, exist := c.Get("userId")
	if !exist {
		response.FailWithMessage("未获取到userid", c)
		return
	}
	userId, ok := id.(uint)
	if !ok {
		response.FailWithMessage("未获取到userid", c)
		return
	}
	//fmt.Println("userid", userId)
	var info model.Template
	info.CreatedID = int(userId)
	info.TemplateName = req.TemplateName
	info.TemplateType = req.TemplateType
	if req.IsOWASP == true {
		info.IsOwasp = 1
	} else {
		info.IsOwasp = 2
	}
	info.ReportLanguage = req.ReportLanguage

	info.Items = itemKeysArray

	if err := app.DB.Model(&model.Template{}).Create(&info).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("创建成功", c)
}

type GetTemplateRequest struct {
	PageSize     int    `form:"size"`
	PageNumber   int    `form:"page"`
	TemplateType string `form:"template_type" binding:"required"` //  安卓ad 苹果 ios 小程序 mp
	IsPage       int    `form:"is_page" binding:"required"`       // 是否需要分页 1 是  2 不是
}

// GetTemplate 获取模版列表
func GetTemplate(c *gin.Context) {

	req := GetTemplateRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	adTem := make(map[string]interface{})
	if err := app.DB.Debug().Model(&model.Template{}).Select("created_at,template_name,id,created_id").Where("id = ?", 1).Scan(&adTem).Error; err != nil {
		log.Error("获取安卓通用模板失败", err)
		response.FailWithMessage("获取安卓通用模板失败", c)
		return
	}

	iosTem := make(map[string]interface{})
	if err := app.DB.Debug().Model(&model.Template{}).Select("created_at,template_name,id,created_id").Where("id = ?", 2).Scan(&iosTem).Error; err != nil {
		log.Error("获取ios通用模板失败", err)
		response.FailWithMessage("获取ios通用模板失败", c)
		return
	}

	isSuper, _ := c.Get("superAdmin")
	isSuperAdmin, ok := isSuper.(bool)
	if !ok {
		log.Error("获取超级管理员标识错误")
		response.FailWithMessage("超级管理员标识错误", c)
		return
	}

	departId, _ := c.Get("departmentId")
	departmentId, ok := departId.(uint)
	if !ok {
		log.Error("获取该员工部门ID失败")
		response.FailWithMessage("获取该员工部门ID失败", c)
		return
	}

	isAdm, _ := c.Get("isAdmin")
	isAdmin, ok := isAdm.(bool)
	if !ok {
		log.Error("获取该员工是否为部门管理员失败")
		response.FailWithMessage("获取该员工是否为部门管理员失败", c)
		return
	}

	getUserId, _ := c.Get("userId")
	userId, ok := getUserId.(uint)
	if !ok {
		response.FailWithMessage("获取该员工ID失败", c)
		return
	}

	info := make([]map[string]interface{}, 0)
	var total int64
	if req.IsPage == 1 {
		// 1.如果你是超级管理员就获取所有
		if isSuperAdmin {

			app.DB.Model(&model.Template{}).Where("template_type = ?", req.TemplateType).
				Select("created_at,template_name,id,created_id").
				Count(&total).
				Offset((req.PageNumber - 1) * req.PageSize).
				Limit(req.PageSize).
				Order("created_at desc").
				Scan(&info)

			for _, v := range info {
				v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
			}
			response.OkWithList(info, int(total), req.PageNumber, req.PageSize, c)
			return
		}

		// 2.如果你是部门管理员的话，就获取该部门下的所有
		if isAdmin {

			app.DB.Model(&model.Template{}).Where("template.template_type = ?", req.TemplateType).
				Where("user.department_id = ? ", departmentId).
				Joins("inner join user on user.id = template.created_id").
				Select("template.created_at,template.template_name,template.id,template.created_id ").
				Count(&total).
				Offset((req.PageNumber - 1) * req.PageSize).
				Limit(req.PageSize).
				Order("template.created_at desc").
				Scan(&info)

			if req.TemplateType == "ad" {
				if len(adTem) != 0 {
					info = append(info, adTem)
				}
			}

			if req.TemplateType == "ios" {
				if len(iosTem) != 0 {
					info = append(info, iosTem)
				}
			}

			for _, v := range info {
				v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
			}
			response.OkWithList(info, int(total+1), req.PageNumber, req.PageSize, c)
			return
		}

		// 3.如果你是普通用户的话
		app.DB.Model(&model.Template{}).Where("template_type = ?", req.TemplateType).
			Where("created_id = ? ", userId).
			Select("created_at,template_name,id,created_id").
			Count(&total).
			Offset((req.PageNumber - 1) * req.PageSize).
			Limit(req.PageSize).
			Order("created_at desc").
			Scan(&info)

		if req.TemplateType == "ad" {
			if len(adTem) != 0 {
				info = append(info, adTem)
			}
		}

		if req.TemplateType == "ios" {
			if len(iosTem) != 0 {
				info = append(info, iosTem)
			}
		}

		for _, v := range info {
			v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
		}

		response.OkWithList(info, int(total+1), req.PageNumber, req.PageSize, c)
		return

	} else {

		// 1.如果你是超级管理员就获取所有
		if isSuperAdmin {

			app.DB.Model(&model.Template{}).Where("template_type = ?", req.TemplateType).
				Select("created_at,template_name,id ").
				Order("created_at desc").
				Scan(&info)

			for _, v := range info {
				v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
			}
			response.OkWithData(info, c)
			return
		}

		// 2.如果你是部门管理员的话，就获取该部门下的所有
		if isAdmin {

			app.DB.Model(&model.Template{}).Where("template.template_type = ?", req.TemplateType).
				Where("user.department_id = ? ", departmentId).
				Joins("inner join user on user.id = template.created_id").
				Select("template.created_at,template.template_name,template.id ").
				Order("template.created_at desc").
				Scan(&info)

			//
			if req.TemplateType == "ad" {
				if len(adTem) != 0 {
					info = append(info, adTem)
				}
			}

			if req.TemplateType == "ios" {
				if len(iosTem) != 0 {
					info = append(info, iosTem)
				}
			}

			for _, v := range info {
				v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
			}
			response.OkWithData(info, c)
			return
		}

		// 3.如果你是普通用户的话

		app.DB.Model(&model.Template{}).Where("template_type = ?", req.TemplateType).
			Where("created_id = ? ", userId).
			Select("created_at,template_name,id ").
			Order("created_at desc").
			Scan(&info)

		if req.TemplateType == "ad" {
			if len(adTem) != 0 {
				info = append(info, adTem)
			}
		}

		if req.TemplateType == "ios" {
			if len(iosTem) != 0 {
				info = append(info, iosTem)
			}
		}

		for _, v := range info {
			v["created_at"] = v["created_at"].(time.Time).Format("2006-01-02 15:04:05")
		}
		response.OkWithData(info, c)
		return
	}
}

type FixTemplateRequest struct {
	TemplateType   string `form:"template_type" binding:"required"` //  模板类型 1 android 2 ios 3 小程序
	TemplateId     int    `form:"template_id" binding:"required"`
	TemplateName   string `form:"template_name"`
	ItemKeys       string `form:"item_keys"` //测评项
	IsOWASP        bool   `form:"is_owasp"`
	ReportLanguage string `form:"report_language"` //导出模板语言
}

// FixTemplate 修改模版
func FixTemplate(c *gin.Context) {
	req := FixTemplateRequest{}
	err := c.ShouldBind(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	str1 := strings.ReplaceAll(req.ItemKeys, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")

	itemKeysArray := "["

	for key, value := range idArray {
		if key != len(idArray)-1 {
			itemKeysArray += "\"" + value + "\","

		} else {
			itemKeysArray += "\"" + value + "\""
		}
	}
	itemKeysArray += "]"
	var info model.Template
	info.ID = uint(req.TemplateId)
	info.TemplateName = req.TemplateName
	info.ReportLanguage = req.ReportLanguage
	info.Items = itemKeysArray

	//fmt.Println("info", info)
	if req.IsOWASP == true {
		info.IsOwasp = 1
	} else {
		info.IsOwasp = 2
	}

	if err := app.DB.Model(&model.Template{}).Where("id = ?", req.TemplateId).Updates(&info).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("修改成功", c)
}

type DelTemplateRequest struct {
	TemplateId int `form:"template_id" binding:"required"`
}

// DeleteTemplate 删除模版
func DeleteTemplate(c *gin.Context) {
	req := DelTemplateRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	var info model.Template
	if err := app.DB.Model(&model.Template{}).Where("id = ?", req.TemplateId).Delete(&info).Error; err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("删除成功", c)

}

type GetTemplateItemsRequest struct {
	TemplateType string `form:"template_type" binding:"required"`
	TemplateId   int    `form:"template_id" ` // 0 全部的
}

// GetTemplateItems 获取模版详细测评项
func GetTemplateItems(c *gin.Context) {
	req := GetTemplateItemsRequest{}
	valid, errs := app.BindAndValid(c, &req)
	if !valid {
		response.FailWithMessage(errs.Error(), c)
		return
	}

	//buff := &bytes.Buffer{}
	//writer := multipart.NewWriter(buff)
	//paramMap := make(map[string]interface{})
	//paramMap["token"] = app.Conf.CePing.Token
	//paramMap["signature"] = app.Conf.CePing.Signature
	//paramMap["type"] = req.TemplateType
	//paramMap["template_id"] = 0
	//
	//value, err := jsoniter.Marshal(paramMap)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//err = writer.WriteField("param", string(value))
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//err = writer.Close()
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//clientURL := IP + "/v4/template/view"
	//
	////fmt.Println(clientURL)
	////生成post请求
	//client := &http.Client{}
	//request, err := http.NewRequest("POST", clientURL, buff)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	////注意别忘了设置header
	//request.Header.Set("Content-Type", writer.FormDataContentType())
	//
	////Do方法发送请求
	//resp, err := client.Do(request)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//post, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//reponse := make(map[string]interface{})
	//err = jsoniter.Unmarshal(post, &reponse)
	//if err != nil {
	//	return
	//}
	//err = Check(reponse, post, GetTemplateItems, c)
	//if err != nil {
	//	return
	//}

	type TemplateItem struct {
		ItemKey   string `json:"item_key"`
		AuditName string `json:"audit_name"`
		IsDynamic bool   `json:"is_dynamic"`
		Status    int    `json:"status"`
	}
	type view struct {
		TemplateId       int                       `json:"template_id"`
		TemplateName     string                    `json:"template_name"`
		CreatorAccount   string                    `json:"creator_account"`
		CreateTime       string                    `json:"create_time"`
		IsOWASP          bool                      `json:"is_owasp"`
		ReportLanguage   string                    `json:"report_language"`
		Categories       []string                  `json:"categories"`
		CategorizedItems map[string][]TemplateItem `json:"categorized_items"`
	}

	//var getView view
	//_ = jsoniter.Unmarshal(post, &getView)
	//
	//var cat model.Category
	//cat.CePingType = "mp"
	//for _, value := range getView.Categories {
	//	if value == "第三方SDK检测" || value == "内容安全" {
	//		continue
	//	}
	//	cat.CategoryName = value
	//
	//	if err := app.DB.Model(&model.Category{}).Create(&cat).Error; err != nil {
	//		response.FailWithMessage(err.Error(), c)
	//		return
	//	}
	//	cat.Id = 0
	//}

	//var cats []model.Category
	//if err := app.DB.Model(&model.Category{}).Find(&cats).Error; err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//for key, tem := range getView.CategorizedItems {
	//	for _, cat := range cats {
	//		if cat.CategoryName == key {
	//			for _, value := range tem {
	//				var itemkey model.TemplateItem
	//				itemkey.CePingType = "mp"
	//				itemkey.ItemKey = value.ItemKey
	//				itemkey.AuditName = value.AuditName
	//				itemkey.IsDynamic = value.IsDynamic
	//				itemkey.Status = value.Status
	//				itemkey.CategoryId = int(cat.Id)
	//				if err := app.DB.Model(&model.TemplateItem{}).Create(&itemkey).Error; err != nil {
	//					response.FailWithMessage(err.Error(), c)
	//					return
	//				}
	//				itemkey.Id = 0
	//			}
	//		}
	//	}
	//}
	if req.TemplateId == 0 {

		toView, err := TemplateItemKeysToView(app.DB, req.TemplateType, nil)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		if req.TemplateType == "ad" {
			var newCataries []string
			for _, value := range toView.Categories {
				if value == "第三方SDK检测" || value == "内容安全" || value == "优化建议" {
					continue
				}
				newCataries = append(newCataries, value)
			}
			toView.Categories = newCataries
			response.OkWithData(toView, c)
			return
		}
		if req.TemplateType == "ios" {
			var newCataries []string
			for _, value := range toView.Categories {
				if value == "第三方SDK检测" || value == "内容安全" {
					continue
				}
				newCataries = append(newCataries, value)
			}
			toView.Categories = newCataries
			response.OkWithData(toView, c)
			return
		}
	}

	var record model.Template

	record.TemplateType = req.TemplateType
	if err := app.DB.Model(&model.Template{}).Where("id = ?", req.TemplateId).Where("template_type = ?", req.TemplateType).First(&record).Error; err != nil {
		response.FailWithMessage("未找到该模板", c)
		return
	}
	str1 := strings.ReplaceAll(record.Items, "[", "")
	str2 := strings.ReplaceAll(str1, "]", "")
	idArray := strings.Split(str2, ",")
	keysToView, err := TemplateItemKeysToView(app.DB, req.TemplateType, idArray)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	keysToView.TemplateId = int(record.ID)
	keysToView.TemplateName = record.TemplateName
	keysToView.CreatorAccount = strconv.Itoa(record.CreatedID)
	keysToView.CreateTime = record.CreatedAt.Format("2006-01-02 15:04:05")
	if record.IsOwasp == 1 {
		keysToView.IsOWASP = true
	} else {
		keysToView.IsOWASP = false
	}
	keysToView.ReportLanguage = record.ReportLanguage
	//var itemKeys []string
	itemKeySet := make(map[string]int, 0)
	item1 := strings.ReplaceAll(record.Items, "[", "")
	item2 := strings.ReplaceAll(item1, "]", "")
	item3 := strings.Split(item2, ",")

	for _, value := range item3 {
		str3 := strings.ReplaceAll(value, "\"", "")
		itemKeySet[str3] = 1
	}
	for _, items := range keysToView.CategorizedItems {
		for i := range items {
			itemKey := items[i].ItemKey
			if _, ok := itemKeySet[itemKey]; ok {
				items[i].Status = 1
			}
		}
	}

	if req.TemplateType == "ad" {
		var newCataries []string
		for _, value := range keysToView.Categories {
			if value == "第三方SDK检测" || value == "内容安全" || value == "优化建议" {
				continue
			}
			newCataries = append(newCataries, value)
		}
		keysToView.Categories = newCataries
		response.OkWithData(keysToView, c)
		return
	}
	if req.TemplateType == "ios" {
		var newCataries []string
		for _, value := range keysToView.Categories {
			if value == "第三方SDK检测" || value == "内容安全" {
				continue
			}
			newCataries = append(newCataries, value)
		}
		keysToView.Categories = newCataries
		response.OkWithData(keysToView, c)
		return
	}

	//response.OkWithData(keysToView, c)
	//return

	//// 1。默认值
	//var viewInfo view
	//viewInfo.TemplateId = 0
	//viewInfo.TemplateName = ""
	//viewInfo.CreatorAccount = "cd_admin"
	//viewInfo.CreateTime = "0001-01-01 00:00:00"
	//viewInfo.IsOWASP = false
	//viewInfo.ReportLanguage = "zh_cn"
	//viewInfo.CategorizedItems = make(map[string][]TemplateItem)
	//// 2.复制里面
	//var tem []model.Category
	//if err := app.DB.Where("ce_ping_type = ?", req.TemplateType).Find(&tem).Error; err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//var temsItem []model.TemplateItem
	//if err := app.DB.Where("ce_ping_type = ?", req.TemplateType).Find(&temsItem).Error; err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//
	//for _, value := range tem {
	//	viewInfo.Categories = append(viewInfo.Categories, value.CategoryName)
	//	for _, templateItem := range temsItem {
	//		if templateItem.CategoryId == int(value.Id) {
	//			var templateItemInfo TemplateItem
	//			templateItemInfo.ItemKey = templateItem.ItemKey
	//			templateItemInfo.AuditName = templateItem.AuditName
	//			templateItemInfo.IsDynamic = templateItem.IsDynamic
	//			templateItemInfo.Status = templateItem.Status
	//			viewInfo.CategorizedItems[value.CategoryName] = append(viewInfo.CategorizedItems[value.CategoryName], templateItemInfo)
	//		}
	//	}
	//}
	//
	//var record model.Template
	//
	//record.TemplateType = req.TemplateType
	//if err := app.DB.Model(&model.Template{}).Where("id = ?", req.TemplateId).Where("template_type = ?", req.TemplateType).First(&record).Error; err != nil {
	//	response.FailWithMessage("未找到该模板", c)
	//	return
	//}
	////fmt.Println("查找到的记录", record)
	//viewInfo.TemplateId = int(record.ID)
	//viewInfo.TemplateName = record.TemplateName
	//viewInfo.CreateTime = record.CreatedAt.Format("2006-01-02 15:04:05")
	//if record.IsOwasp == 1 {
	//	viewInfo.IsOWASP = true
	//} else {
	//	viewInfo.IsOWASP = false
	//}
	//viewInfo.ReportLanguage = record.ReportLanguage
	//
	////var itemKeys []string
	//itemKeySet := make(map[string]int, 0)
	//str1 := strings.ReplaceAll(record.Items, "[", "")
	//str2 := strings.ReplaceAll(str1, "]", "")
	//idArrayStr := strings.Split(str2, ",")
	//
	//for _, value := range idArrayStr {
	//	str3 := strings.ReplaceAll(value, "\"", "")
	//	itemKeySet[str3] = 1
	//}
	//for _, items := range viewInfo.CategorizedItems {
	//	for i := range items {
	//		itemKey := items[i].ItemKey
	//		if _, ok := itemKeySet[itemKey]; ok {
	//			items[i].Status = 1
	//		}
	//	}
	//}
	//
	//response.OkWithData(viewInfo, c)

}
