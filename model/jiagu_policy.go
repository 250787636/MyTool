package model

type JiaguPolicy struct {
	Id           int    `json:"id"`             //策略id
	Status       string `json:"status"`         // 策略状态
	Name         string `json:"name"`           // 策略名称
	Template     string `json:"Template"`       // 策略配置
	Type         string `json:"type"`           // 加固策略类型
	LicKeyHelper string `json:"lic_key_helper"` // 加固策略类型
}
