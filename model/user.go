package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName      string `json:"name"`            // 用户名
	DepartmentID  uint   `json:"department_id"`   // 部门id
	AccountLevel  string `json:"account_level"`   // 账户等级
	JobTitle      string `json:"job_title"`       // 职位名称
	LastLoginTime string `json:"last_login_time"` // 最近一次登录

	Token   string `json:"token"`    // 生成验证token
	IsAdmin bool   `json:"is_admin"` // 是否是管理员
}
