package main

import (
	"example.com/m/model/automigrate"
	"example.com/m/pkg/app"
	"example.com/m/pkg/log"
	"example.com/m/router"
	"example.com/m/router/api/ceping"
	"example.com/m/router/api/jiagu"
	"example.com/m/utils"
	"flag"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	flag.Parse()
	// 将请求放入控制台与run.log中
	//gin.DefaultWriter = io.MultiWriter(os.Stdout, log.FileWriter)
	r := gin.Default()
	// 进行建表操作
	automigrate.Run()
	// 进行服务注册操作
	router.Init(r)
	// 过5分钟进行加固策略更新
	go jiagu.GetPolicyListByType()
	// 过1个小时重新获取测评平台token
	go ceping.TimeToGetToken()
	go utils.TickerDel("datastorage/ak/upload/scan_upload", "data/scan_upload", 24*time.Hour)
	log.Info(`                               
		 /$$$$$$$   /$$$$$$  /$$$$$$ /$$$$$$$$ /$$$$$$ 
		| $$__  $$ /$$__  $$|_  $$_/|__  $$__//$$__  $$
		| $$  \ $$| $$  \__/  | $$     | $$  | $$  \__/
		| $$$$$$$/| $$        | $$     | $$  | $$      
		| $$____/ | $$        | $$     | $$  | $$      
		| $$      | $$    $$  | $$     | $$  | $$    $$
		| $$      |  $$$$$$/ /$$$$$$   | $$  |  $$$$$$/
		|__/       \______/ |______/   |__/   \______/
	`)
	log.Fatal(r.Run(":" + app.Conf.System.Port))
}
