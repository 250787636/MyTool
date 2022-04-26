package router

import (
	"example.com/m/pkg/log"
	"example.com/m/router/api"
	"example.com/m/router/api/ceping"
	"example.com/m/router/api/jiagu"
	"example.com/m/router/api/user"
	"example.com/m/router/middleware"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {
	// 设置gin.Default().Recovery()中间件输出到日志文件中
	gin.DefaultErrorWriter = log.DefalutLogger.Writer()
	// recover panic
	r.Use(gin.Recovery())
	// 加入中间件 进行ip验证
	r.Use(middleware.CheckAddress)
	r.GET("/api/login", user.SsoLogin)
	userApi := r.Group("/mspapi/open_api")
	{
		userApi.POST("user/add_user", user.AddUser)
		userApi.POST("user/fix_user", user.FixUser)
		userApi.POST("user/del_user", user.DelUser)
	}
	// 加入中间件验证token
	r.Use(middleware.CheckToken)

	userInApi := r.Group("/api/user")
	{
		userInApi.Use(middleware.AdminOnly)
		userInApi.POST("/list_user", user.ListUser)
		userInApi.POST("/list_user_without_page", user.ListUserWithOutPage)
		userInApi.POST("/user_permission", user.UserPermission)
		userInApi.POST("/user_status", user.UserStatus)

	}
	// android 测评部分
	android := r.Group("/api/ceping/v4")
	//android.Use(middleware.Translations())
	{
		android.POST("/bin_check_apk", ceping.BinCheckApk)

		android.POST("/search_one_progress", ceping.SearchOneProgress)
		android.POST("/search_one_detail", ceping.SearchOneDetail)
		android.GET("/apk_report", ceping.ApkReport)
		android.POST("/batch_statistics_result", ceping.BatchStatisticsResult)
		android.POST("/batch_statistics_original_result", ceping.BatchStatisticsOriginalResult)
		android.POST("/batch_file_delete", ceping.BatchFileDelete)
		android.POST("/get_all", ceping.GetAllInfo) // 获取ad测评列表
		//android.POST("/get_items", ceping.GetItems)
		//android.POST("/get_all_items", ceping.GetAllItems)
		//android.POST("/local_check_apk", ceping.LocalCheckApk)
		//android.POST("/url_check_apk", ceping.URLCheckApk)
	}
	// 模板
	template := r.Group("/api/ceping/v4/template")
	//template.Use(middleware.Translations())
	{
		template.POST("/add_template", ceping.AddTemplate)
		template.POST("/get_template", ceping.GetTemplate)
		template.POST("/fix_template", ceping.FixTemplate)
		template.POST("/delete_template", ceping.DeleteTemplate)
		template.POST("/get_template_items", ceping.GetTemplateItems)
	}

	// ios 测评部分
	ios := r.Group("/api/ceping/v4/ios")
	//ios.Use(middleware.Translations())
	{
		ios.POST("/bin_check", ceping.IosBinCheck)
		ios.POST("/search_one_detail", ceping.IosSearchOneDetail)
		ios.POST("/batch_statistics_result", ceping.IosBatchStatisticsResult)
		ios.GET("/ipa_report", ceping.IosIpaReport)
		ios.POST("/batch_file_delete", ceping.IosBatcFileDelete)
	}
	// 小程序 测评部分
	mp := r.Group("/api/ceping/v4/mp")
	//mp.Use(middleware.Translations())
	{
		mp.POST("/bin_check", ceping.MpBinCheck)
		mp.POST("/search_mp", ceping.MpSearch)
		mp.GET("/mini_report", ceping.MpMiniReport)
		mp.POST("/re_bin_check", ceping.MpReBinCheck)
		mp.POST("/batch_file_delete", ceping.MpBatcFileDelete)
	}
	// 测评其他接口
	other := r.Group("/api/ceping/v4")
	//other.Use(middleware.Translations())
	{ // 测评平台数据统计
		other.GET("/batch_download", ceping.BatchDownload) // 获取报告
		//other.GET("/data_export", ceping.DataExport)
	}

	// android加固部分
	jiaguAndroid := r.Group("/api/jiagu/v5")
	{
		jiaguAndroid.POST("/upload", jiagu.WebBoxV5Upload)
		jiaguAndroid.POST("/get", jiagu.WebBoxV5GetState)
		jiaguAndroid.GET("/download", jiagu.WebBoxV5Download)
		jiaguAndroid.POST("/delete", jiagu.WebBoxV5Delete)
		jiaguAndroid.GET("/download_log", jiagu.WebBoxV5DownloadLog)
		jiaguAndroid.POST("/get_all", jiagu.WebBoxV5AllTask)
		jiaguAndroid.GET("/exporting", jiagu.Exporting)
	}
	// IOS加固部分
	jiaguIos := r.Group("/api/jiagu/i5")
	{
		//jiaguIos.POST("/get", jiagu.WebBoxFindIosEdition)
		//jiaguIos.POST("/modify", jiagu.WebBoxModifyIosEdition)
		jiaguIos.GET("/downland_ios_file", api.DownlandIosFile) // 获取客户端信息文件
		jiaguIos.Use(middleware.AdminOnly)
		jiaguIos.POST("/downland_ios_list", jiagu.IosDownlandList) // 获取客户端下载记录

	}

	// 加固策略部分
	jiaguPolicy := r.Group("api/jiagu/policy")
	{
		jiaguPolicy.POST("/policy_get", jiagu.JiaguPolicyFind)
		jiaguPolicy.Use(middleware.AdminOnly)
		jiaguPolicy.POST("/policy_related", jiagu.JiaguPolicyRelated)
		jiaguPolicy.POST("/policy_get_with_page", jiagu.JiaguPolicyFindWithPage)
	}

	// 获取文件信息接口
	file := r.Group("/api/file")
	//file.Use(middleware.Translations())
	{
		file.POST("/get", api.GetFile)
		file.GET("/download", api.DownloadFile) // 获取源文件
	}
}
