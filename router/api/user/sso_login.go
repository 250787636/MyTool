package user

import (
	"example.com/m/model"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/utils"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

func SsoLogin(r *gin.Context) {
	ip := app.Conf.Sso.Ip
	stringUrl := ip + "/dmcwebapi/api/dmc/ReqToken"

	token := r.Query("token")
	if token == "" {
		r.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusForbidden, // 没有访问权限
			"err":  "地址栏token不存在",
		})
		return
	}
	// 请求体参数
	m := map[string]interface{}{
		"type":    "sgtoken",
		"imtoken": token,
	}
	req, err := utils.NewJsonRequest(stringUrl, m)
	if err == nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Error(err.Error())
		}
		var TokenResponse struct {
			Status     string `json:"status"`
			FailureMsg string `json:"failure_msg"`
			Uid        string `json:"uid"`
		}
		err = jsoniter.Unmarshal(body, &TokenResponse)
		if err != nil {
			log.Error(err.Error())
		}

		if TokenResponse.Status != "1" {
			r.JSON(http.StatusOK, gin.H{
				"code": http.StatusForbidden, // 没有访问权限
				"info": gin.H{"status": TokenResponse.Status},
				"err":  TokenResponse.FailureMsg,
			})
		} else {
			// 获取到唯一标识后 在数据库中查询是否有当前账户
			var userinfo []model.UserInfo
			if err := app.DB.Where("username = ?", TokenResponse.Uid).Find(&userinfo).Error; err != nil || len(userinfo) == 0 {
				r.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusForbidden, // 没有访问权限
					"info": gin.H{"status": TokenResponse.Status},
					"err":  model.UserNotExistInDataBase,
				})
			} else if userinfo[0].IsActive == 0 { //IsActive为0时账户停用
				r.JSON(http.StatusBadRequest, gin.H{
					"code": http.StatusForbidden, // 没有访问权限
					"info": gin.H{"status": "0"},
					"err":  "该账户已停用",
				})
			} else {
				r.JSON(http.StatusOK, gin.H{
					"code": http.StatusOK,
					"info": gin.H{
						"status":          TokenResponse.Status,
						"uid":             TokenResponse.Uid, // 获取用户信息成功
						"user_type":       userinfo[0].UserType,
						"user_permission": userinfo[0].UserPermission,
					},
					"msg": model.ReqSuccess,
				})
			}
		}
	}
}
