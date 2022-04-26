package automigrate

import "example.com/m/model"

func setTable(migrate *AutoMigrate) {
	// 数据表多用于 c u d
	migrate.DataTables = []TableVersion{
		{model.JiaGuTask{}, ""},
		{model.CePingUserTask{}, ""},
		{model.JiaguPolicy{}, ""},
		{model.UserInfo{}, ""},
		{model.Source_Code_Reinforcement_Log{}, ""},
	}
	//工具表多用于 r
	migrate.ToolTables = []TableVersion{
		// 模板分类表
		{model.Template{}, ""},
		{model.Category{}, ""},
		{model.TemplateItem{}, ""},
		{model.CepingAdAuditItem{}, ""},
		{model.CepingIosAuditItem{}, ""},
		{model.CepingAuditCategory{}, ""},
		// IOS客户端数据
		{model.IosClientDownloadPath{}, ""},
		// 预设置用户等级
		{model.UserType{}, ""},
	}
}
