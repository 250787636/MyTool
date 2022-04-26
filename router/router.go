package router

import (
	"example.com/m/pkg/log"
	"example.com/m/router/api"
	"example.com/m/router/api/ceping"
	"example.com/m/router/api/jiagu"
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
	// 如果不存在则创建账号
	r.POST("/api/check_account", api.CheckAccount)
	// 加入中间件验证token
	r.Use(middleware.CheckToken)
	// android 测评部分
	android := r.Group("/api/ceping/v4")
	//android.Use(middleware.Translations())
	{
		android.POST("/bin_check_apk", ceping.BinCheckApk)
		android.POST("/search_one_progress", ceping.SearchOneProgress)
		android.GET("/apk_report", ceping.ApkReport)
		android.POST("/batch_statistics_result", ceping.BatchStatisticsResult)
		android.POST("/batch_file_delete", ceping.BatchFileDelete)
		android.POST("/search_one_detail", ceping.SearchOneDetail)
		android.POST("/get_all", ceping.GetAllInfo) // 获取ad测评列表
		//android.POST("/batch_statistics_original_result", ceping.BatchStatisticsOriginalResult)
		//android.POST("/set_items", ceping.SetItems)
		//android.POST("/get_items", ceping.GetItems)
		//android.POST("/get_all_items", ceping.GetAllItems)
		//android.POST("/url_check_apk", ceping.URLCheckApk)
		//android.POST("/local_check_apk", ceping.LocalCheckApk)

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
		ios.GET("/ipa_report", ceping.IosIpaReport)
		ios.POST("/batch_file_delete", ceping.IosBatcFileDelete)
		//ios.POST("/batch_statistics_result", ceping.IosBatchStatisticsResult)
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
		other.POST("/ext/get_details", ceping.GetDetails)  // 获取测评平台数据统计
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
		//jiaguAndroid.POST("/get_android_ver", jiagu.WebBoxV5ReinforceVer)
		//jiaguAndroid.POST("/get_so_ver", jiagu.WebBoxV5SoVer)
		//jiaguAndroid.POST("/get_policy_user", jiagu.WebBoxV5PolicyUser)
		//jiaguAndroid.POST("/policy_add", jiagu.WebBoxV5PolicyAdd)
		//jiaguAndroid.POST("/policy_modify", jiagu.WebBoxV5PolicyModify)
		//jiaguAndroid.POST("/policy_delete", jiagu.WebBoxV5PolicyDelete)
		jiaguAndroid.POST("/get_all", jiagu.WebBoxV5AllTask)
		jiaguAndroid.Use(middleware.AdminOnly)
		jiaguAndroid.GET("/exporting", jiagu.Exporting)
	}

	// h5加固部分
	jiaguHtml5 := r.Group("/api/jiagu/h5")
	{
		jiaguHtml5.POST("/upload", jiagu.WebBoxH5Upload)
		jiaguHtml5.POST("/get", jiagu.WebBoxH5GetState)
		jiaguHtml5.GET("/download", jiagu.WebBoxH5Download)
		jiaguHtml5.POST("/delete", jiagu.WebBoxH5Delete)
		jiaguHtml5.GET("/download_log", jiagu.WebBoxH5DownloadLog)
		//jiaguHtml5.POST("/getH5Ver", jiagu.WebBoxH5Ver)
		//jiaguHtml5.POST("/get_policy_user", jiagu.WebBoxH5PolicyUser)
		//jiaguHtml5.POST("/policy_add", jiagu.WebBoxH5PolicyAdd)
		//jiaguHtml5.POST("/policy_modify", jiagu.WebBoxH5PolicyModify)
		//jiaguHtml5.POST("/policy_delete", jiagu.WebBoxH5PolicyDelete)
		jiaguHtml5.POST("/get_all", jiagu.WebBoxH5TaskAllTask)
		jiaguHtml5.Use(middleware.AdminOnly)
		jiaguHtml5.GET("/exporting", jiagu.Exporting)
	}

	// 加固策略部分
	jiaguPolicy := r.Group("api/jiagu/policy")
	{
		jiaguPolicy.POST("/policy_get", jiagu.JiaguPolicyFind)
		jiaguPolicy.POST("/policy_get_with_page", jiagu.JiaguPolicyFindWithPage)
	}

	// 加固数据统计部分
	jiaguCount := r.Group("/api/jiagu/count")
	{
		jiaguCount.POST("/get_by_department", jiagu.JiaGuDepartmentCount)
		jiaguCount.POST("/get_by_application", jiagu.JiaGuApplicationCount)
	}

	// 获取账户信息部分
	accountMethod := r.Group("/api/account")
	accountMethod.Use(middleware.AdminOnly)
	{
		accountMethod.POST("/get", api.GetListAccount)
		accountMethod.POST("/modify", api.ModifyAccount)
		accountMethod.POST("/get_department", api.GetDepartment)
	}

	// 应用管理部分
	applicationManagement := r.Group("/api/application")
	{
		applicationManagement.POST("/get", jiagu.GetApplication)
		applicationManagement.POST("/get_by_id", jiagu.GetApplicationShow)
		applicationManagement.POST("/get_type", jiagu.GetApplicationType)
		applicationManagement.POST("/search", jiagu.SearchList)
		applicationManagement.Use(middleware.AdminOnly)
		{
			applicationManagement.POST("/create", jiagu.CreateApplication)
			applicationManagement.POST("/edit", jiagu.EditApplication)
			applicationManagement.POST("/delete", jiagu.DelApplication)
			applicationManagement.POST("/edit_policy", jiagu.ModifyRecommendPolicy)
		}
	}
	//操作手册
	operationManual := r.Group("/api/handbook")
	{
		operationManual.POST("/get_all", jiagu.HandBookGetAll)
		operationManual.GET("/download", jiagu.HandBookDownland)
		operationManual.Use(middleware.AdminOnly)
		{
			operationManual.POST("/get_service", jiagu.HandBookGetServiceType)
			operationManual.POST("/create", jiagu.HandBookCreate)
			operationManual.POST("/delete", jiagu.HandBookDelete)
		}
	}

	// 获取文件信息接口
	file := r.Group("/api/file")
	//file.Use(middleware.Translations())
	{
		file.POST("/get", api.GetFile)
		file.GET("/download", api.DownloadFile) // 获取源文件
		file.POST("/import/excel", api.NewExcel)
	}
}
