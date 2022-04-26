package app

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"os"
)

var Conf = LoadConfig()

type config struct {
	System struct {
		Version  string `ini:"version"`
		RootInit string `ini:"root_init"`
		Port     string `ini:"port"`
	} `ini:"system"`
	Mysql struct { // mysql配置项
		Host     string `ini:"server" json:"host"`
		Port     string `ini:"port" json:"port"`
		Username string `ini:"user" json:"-"`
		Password string `ini:"pwd" json:"-"`
		Database string `ini:"name" json:"database"`
		ResetDb  bool   `ini:"reset_db" json:"is_delete"`
	} `ini:"mysql"`
	Minio struct { // minio配置项
		Endpoint        string `ini:"endpoint"`
		AccessKeyID     string `ini:"access_key_id"`
		SecretAccessKey string `ini:"secret_access_key"`
		UseSSL          bool   `ini:"use_ssl"`
	} `ini:"minio"`
	Engine struct { // 引擎配置项
		RequestUrl  string `ini:"request_url"`
		CallbackUrl string `ini:"callback_url"`
		ClientId    string `ini:"client_id"`
		Key         string `ini:"key"`
	} `ini:"engine"`
	Platforms struct { // 3s平台传入用户名的固定ip
		IP string `ini:"ip"`
	} `ini:"platforms"`
	CePing struct { // 测评平台内置账户
		UserName  string `ini:"username" json:"username"`
		Password  string `ini:"password" json:"password"`
		Token     string `ini:"token" json:"token"`
		Signature string `ini:"signature" json:"signature"`
		Ip        string `ini:"ip" json:"ip"`
	} `ini:"ceping"`
	JiaGu struct { // 加固平台内置账户
		UserName  string `ini:"username" json:"username"`
		ApiKey    string `ini:"api_key" json:"api_key"`
		ApiSecret string `ini:"api_secret" json:"api_secret"`
		Ip        string `ini:"ip" json:"ip"`
		FileName  string `ini:"file_name" json:"file_name"`
	} `ini:"jiagu"`
	Sso struct { // 单点登录请求地址
		Ip string `ini:"ip" json:"ip"`
	} `ini:"sso"`
}

// LoadConfig 加载config文件
func LoadConfig() config {
	var (
		defaultConf  = RootDir + "/deploy.ini"
		customerConf = RootDir + "/customer.ini"
	)
	var conf config
	// 判断是否有costomer
	_, err := os.Open(customerConf)
	if err != nil {
		err = ini.MapTo(&conf, defaultConf)
	} else {
		err = ini.MapTo(&conf, defaultConf, customerConf)
	}
	if err != nil {
		panic(err)
	}

	token, signature, err := GetCpToken(conf.CePing.UserName, conf.CePing.Password, conf.CePing.Ip)
	if err != nil {
		log.Println("获取测评平台token失败,请检查customer.ini配置文件,错误为", err)
	}

	conf.CePing.Token = token
	fmt.Println("token", conf.CePing.Token)
	conf.CePing.Signature = signature
	fmt.Println("signature", conf.CePing.Signature)

	if token == "" || signature == "" {
		log.Println("获取测评平台token失败,请检查customer.ini配置文件")
	}

	return conf
}
