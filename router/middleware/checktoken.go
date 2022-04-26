package middleware

import (
	"example.com/m/model"
	"example.com/m/pkg/app"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CheckToken 验证 Header 上的 token
func CheckToken(r *gin.Context) {
	var token string
	token = r.Request.Header.Get("token")
	if token == "" {
		token = r.Query("token")
		if token == "" {
			r.JSON(http.StatusBadRequest, gin.H{
				"code": http.StatusForbidden, // 没有访问权限
				"err":  model.TokenNotExist,
			})
			r.Abort()
			return
		}
	}
	var userinfo []model.UserInfo
	if err := app.DB.Where("user_token = ?", token).Find(&userinfo).Error; err != nil || len(userinfo) == 0 {
		r.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusForbidden, // 没有访问权限
			"err":  model.UserNotExist,
		})
		r.Abort()
		return
	}
	r.Set("userName", userinfo[0].UserName)
	r.Set("userId", userinfo[0].ID)
	r.Set("userType", userinfo[0].UserType)
}
