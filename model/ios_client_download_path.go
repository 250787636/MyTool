package model

// 源码加固客户端信息
type IosClientDownloadPath struct {
	Id                     uint   `json:"id" gorm:"comment:'源码加固客户端id'"`
	IosClientEdition       string `json:"ios_client_edition" gorm:"comment:'源码加固客户端版本'"`
	IosWindowsName         string `json:"ios_windowns_name" gorm:"comment:'源码加固客户端(windows)名'"`
	IosWindowsDownloadPath string `json:"ios_windowns_download_path" gorm:"comment:'源码加固客户端(windows)地址'"`
	IosMacName             string `json:"ios_mac_name" gorm:"comment:'源码加固客户端(Mac)名'"`
	IosMacDownloadPath     string `json:"ios_mac_download_path" gorm:"comment:'源码加固客户端(Mac)地址'"`
}
