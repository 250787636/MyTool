package model

type UserType struct {
	Id    string `gorm:"column:id;type:varchar(100);NOT NULL;default:0;comment:用户类型 0-普通用户  1-安控管理员 2-安全监测员 3-沙箱管理员 4-沙箱审计员"` // 用户类型
	UType string `gorm:"column:u_type;type:varchar(100);NOT NULL;default:0;comment:用户类型 普通用户  安控管理员 安全监测员 沙箱管理员 沙箱审计员"`       // 用户类型
}
