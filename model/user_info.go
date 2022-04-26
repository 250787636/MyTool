package model

import (
	"gorm.io/gorm"
)

type UserInfo struct {
	ID             uint           `gorm:"primary_key" json:"id"`
	CreatedAt      FormatTime     `json:"created_at"`
	UpdatedAt      FormatTime     `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	UserName       string         `json:"username" gorm:"unique;column:username;type:varchar(100);NOT NULL;default:'';comment:用户名"`               // 用户名
	UserType       string         `json:"user_type" gorm:"column:user_type;type:varchar(100);NOT NULL;default:0;comment:用户类型 0-普通  1-管理员 2-安全监测"` // 用户类型
	Email          string         `json:"email" gorm:"column:email;type:varchar(100);NOT NULL;default:'';comment:邮箱"`                             // 邮箱
	Phone          string         `json:"phone" gorm:"column:phone;type:varchar(100);NOT NULL;default:'';comment:手机号"`                            // 手机号
	IsActive       int            `json:"is_active" gorm:"column:is_active;type:tinyint(1);NOT NULL;default:1;comment:是否激活"`                      // 是否激活  0 false-禁用 1 true-启用
	ExpireTime     FormatTime     `json:"expire_time" gorm:"column:expire_time;type:datetime;comment:过期时间"`                                       // 过期时间
	RealName       string         `json:"realname" gorm:"column:realname;type:varchar(100);NOT NULL;default:'';comment:真实姓名"`                     // 真实姓名
	Company        string         `json:"company" gorm:"column:company;type:varchar(100);NOT NULL;default:'';comment:公司"`                         // 公司
	CompanyId      string         `json:"company_id" gorm:"column:company_id;type:varchar(100);NOT NULL;default:'';comment:公司id"`                 // 公司id
	UserPermission string         `json:"user_permission" gorm:"column:user_permission;type:varchar(100);NOT NULL;default:'';comment:用户权限"`       // 用户权限
	UserToken      string         `json:"user_token" gorm:"column:user_token;type:varchar(100);NOT NULL;default:'';comment:用户token"`              // 用户token
}
