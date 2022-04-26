package model

import (
	"gorm.io/gorm"
	"time"
)

type JiaGuTask struct {
	gorm.Model
	UserID     uint       `json:"user_id"`     // 用户id
	TaskID     int        `json:"task_id"`     // 任务id
	PolicyID   int        `json:"policy_id"`   // 应用策略id
	ApkName    string     `json:"apk_name"`    // 应用名称
	ApkSize    string     `json:"apk_size"`    // 文件大小
	Filename   string     `json:"filename"`    // apk文件名
	Version    string     `json:"version"`     // 应用版本
	TaskStatus string     `json:"task_status"` // 加固任务状态
	FinishTime *time.Time `json:"finish_time"` // 完成时间
}
