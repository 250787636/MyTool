package model

import (
	"gorm.io/gorm"
)

// 源码加固下载记录
type Source_Code_Reinforcement_Log struct {
	ID           uint           `gorm:"primarykey"`
	CreatedAt    FormatTime     `json:"created_at"`
	UpdatedAt    FormatTime     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserName     string         `json:"user_name"`
	DownloadTime string         `json:"download_time"`
	FileName     string         `json:"file_name"`
	IpAddress    string         `json:"ip_address"`
	BrowserName  string         `json:"browser_name"`
	SystemType   string         `json:"system_type"`
}
