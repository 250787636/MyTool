package user

import (
	"crypto/md5"
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/router/middleware"
	"example.com/m/utils"
	"example.com/m/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type Info struct {
	UserName   string `json:"username" binding:"required"`    // 用户名
	UserType   string `json:"user_type" binding:"required"`   // 用户类型
	Email      string `json:"email"`                          // 邮箱
	Phone      string `json:"phone"`                          // 手机号
	IsActive   bool   `json:"is_active" binding:"required"`   // 是否激活  0 false-禁用 1 true-启用
	ExpireTime string `json:"expire_time" binding:"required"` // 过期时间
	RealName   string `json:"realname" binding:"required"`    // 真实姓名
	Company    string `json:"company" binding:"required"`     // 公司
	CompanyId  string `json:"company_id" binding:"required"`  // 公司id
}

// 1.1 AddUser 添加用户
func AddUser(c *gin.Context) {
	var userInfo Info
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	apiKey := c.GetHeader("api_key")
	sign := c.GetHeader("sign")

	// 接口验证
	err := middleware.ApiAuth(apiKey, sign, userInfo)
	if err != nil {
		log.Error("添加用户接口验证失败", err.Error())
		response.FailWithMessage("接口验证失败", c)
		return
	}
	expiredTime, _ := time.ParseInLocation("2006-01-02 15:04:05", userInfo.ExpireTime, time.Local)

	// 三元表达式判断是否是true
	activeBool := TPoint(userInfo.IsActive == true, 1, 0)

	md5data := md5.Sum([]byte(userInfo.UserName))
	userToken := fmt.Sprintf("%x", md5data)

	// 创建用户
	user := model.UserInfo{
		UserName:       userInfo.UserName,
		UserType:       userInfo.UserType,
		Email:          userInfo.Email,
		Phone:          userInfo.Phone,
		IsActive:       activeBool.(int),
		ExpireTime:     model.FormatTime(expiredTime),
		RealName:       userInfo.RealName,
		Company:        userInfo.Company,
		CompanyId:      userInfo.CompanyId,
		UserPermission: `["应用加固","源码加固","安全测评"]`,
		UserToken:      userToken,
	}

	if err := app.DB.Create(&user).Error; err != nil {
		log.Error("创建用户失败:", err.Error())
		response.FailWithMessage("创建用户失败", c)
		return
	}

	/*type dataList struct {
		code int
		msg  string
		info map[interface{}]interface{}
	}
	urlString := fmt.Sprintf("%s/v5/user/add?api_key=%s&contact=%s&name=%s&password=%s&username=%s",
		jiagu.JIAGUIP, jiagu.APIKEY, user.Phone, user.RealName, jiagu.PASSWORD, user.UserName)
	m := map[string]interface{}{
		"api_key":  jiagu.APIKEY,
		"contact":  user.Phone,
		"name":     user.RealName,
		"password": jiagu.PASSWORD,
		"username": user.UserName,
	}
	res, err := jiagu.PostWithoutFile(urlString, m)
	if err != nil {
		log.Error("创建用户失败:", err.Error())
		response.FailWithMessage("创建用户失败", c)
		return
	}
	_ = jsoniter.Unmarshal(res, &dataList{})*/

	response.OkWithMessage("添加用户成功", c)
	return

}

// 1.2 FixUser 修改用户信息
func FixUser(c *gin.Context) {
	var userInfo Info
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	apiKey := c.GetHeader("api_key")
	sign := c.GetHeader("sign")

	// 接口验证
	err := middleware.ApiAuth(apiKey, sign, userInfo)
	if err != nil {
		log.Error("修改用户信息接口验证失败", err.Error())
		response.FailWithMessage("接口验证失败", c)
		return
	}
	expiredTime, _ := time.ParseInLocation("2006-01-02 15:04:05", userInfo.ExpireTime, time.Local)

	// 三元表达式判断是否是true
	activeBool := TPoint(userInfo.IsActive == true, 1, 0)

	md5data := md5.Sum([]byte(userInfo.UserName))
	userToken := fmt.Sprintf("%x", md5data)

	// 创建用户
	user := model.UserInfo{
		UserName:   userInfo.UserName,
		UserType:   userInfo.UserType,
		Email:      userInfo.Email,
		Phone:      userInfo.Phone,
		IsActive:   activeBool.(int),
		ExpireTime: model.FormatTime(expiredTime),
		RealName:   userInfo.RealName,
		Company:    userInfo.Company,
		CompanyId:  userInfo.CompanyId,
		UserToken:  userToken,
	}
	exist := model.UserInfo{}
	if err := app.DB.Model(model.UserInfo{}).Where("username = ?", user.UserName).First(&exist).Error; err != nil {
		log.Error("修改用户信息失败:", err.Error())
		response.FailWithMessage("修改用户信息失败", c)
		return
	} else {
		if err := app.DB.Model(model.UserInfo{}).Where("username = ?", user.UserName).Updates(user).Error; err != nil {
			log.Error("修改用户信息失败:", err.Error())
			response.FailWithMessage("修改用户信息失败", c)
			return
		}
	}
	response.OkWithMessage("修改用户信息成功", c)
	return
}

type DeleteUser struct {
	UserName string `json:"username" binding:"required"`
}

// 1.3 DelUser 删除用户
func DelUser(c *gin.Context) {

	var deleteInfo DeleteUser
	if err := c.ShouldBindJSON(&deleteInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	apiKey := c.GetHeader("api_key")
	sign := c.GetHeader("sign")

	// 接口验证
	err := middleware.ApiAuth(apiKey, sign, deleteInfo)
	if err != nil {
		log.Error("删除用户信息接口验证失败", err.Error())
		response.FailWithMessage("接口验证失败", c)
		return
	}

	var userInfo model.UserInfo

	userInfo.UserName = deleteInfo.UserName

	if err := app.DB.Where("username = ?", userInfo.UserName).Unscoped().Delete(&userInfo).Error; err != nil {
		log.Error("删除用户信息失败:", err.Error())
		response.FailWithMessage("删除用户失败", c)
		return
	}
	response.OkWithMessage("删除用户信息成功", c)
}

// 1.4 ListUser 陈列用户
func ListUser(c *gin.Context) {
	var userList []model.UserInfo
	var query string
	var args []interface{}
	username, flag := c.GetPostForm("username")
	isture := utils.IsStringEmpty(username, flag)
	if isture {
		query = "AND username like ? "
		args = append(args, fmt.Sprintf("%%%s%%", utils.SpecialString(username)))
	}
	realname, flag := c.GetPostForm("realname")
	isture = utils.IsStringEmpty(realname, flag)
	if isture {
		query += "AND realname like ? "
		args = append(args, fmt.Sprintf("%%%s%%", utils.SpecialString(realname)))
	}
	query = strings.TrimPrefix(query, "AND")
	// 将page size offnum封装成工具方法
	page, size, offNum := utils.GetFromDataPageSizeOffNum(c)
	if page+size+offNum == -3 {
		return
	}
	var total int64
	if err := app.DB.Model(model.UserInfo{}).
		Where(query, args...).
		Order("id DESC").
		Count(&total).
		Offset(offNum).
		Limit(size).
		Find(&userList).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": gin.H{
			"datalist": userList,
			"total":    total,
		},
		"msg": model.ReqSuccess,
	})
}

// 1.5 不分页查询用户
func ListUserWithOutPage(c *gin.Context) {
	var userList []model.UserInfo
	// 将page size offnum封装成工具方法

	if err := app.DB.Model(model.UserInfo{}).
		Find(&userList).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": gin.H{
			"datalist": userList,
		},
		"msg": model.ReqSuccess,
	})
}

// 仿写三元表达
func TPoint(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// UserPermission 用户权限
func UserPermission(c *gin.Context) {
	userName, flag := c.GetPostForm("username")
	if !utils.IsStringEmpty(userName, flag) {
		log.Error("username参数为空")
		response.FailWithMessage("username参数为空", c)
		return
	}
	permission, flag := c.GetPostForm("permission")
	if !utils.IsStringEmpty(permission, flag) {
		log.Error("permission参数为空")
		response.FailWithMessage("permission参数为空", c)
		return
	}
	if err := app.DB.Model(model.UserInfo{}).Where("username = ?", userName).Update("user_permission", permission).Error; err != nil {
		log.Error("权限添加失败:", err.Error())
		response.FailWithMessage("权限添加失败", c)
		return
	}
	response.OkWithMessage("权限添加成功", c)
}

// UserStatus 用户状态
func UserStatus(c *gin.Context) {
	userName, flag := c.GetPostForm("username")
	if !utils.IsStringEmpty(userName, flag) {
		log.Error("username参数为空")
		response.FailWithMessage("username参数为空", c)
		return
	}
	isActive, flag := c.GetPostForm("is_active")
	if !utils.IsStringEmpty(isActive, flag) {
		log.Error("username参数为空")
		response.FailWithMessage("username参数为空", c)
		return
	}
	if err := app.DB.Model(model.UserInfo{}).Where("username = ?", userName).Update("is_active", isActive).Error; err != nil {
		log.Error("修改状态失败:", err.Error())
		response.FailWithMessage("修改状态失败", c)
		return
	}
	response.OkWithMessage("修改状态成功", c)
}
