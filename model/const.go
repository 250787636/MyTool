package model

// 常量声明部分

const (
	ReqSuccess = "请求成功"
	CreSuccess = "创建成功"
	ModSuccess = "修改成功"
	DelSuccess = "删除成功"
)

const (
	ReqParameterMissing = "必传参数缺失"
	ReqParameterError   = "必传参数错误"
	ReqParameterModFail = "未传入修改参数，未进行账户修改"
)
const (
	FileAlreadyDelete = "文件已删除"
	FileNotExist      = "文件不存在"
)

const (
	TokenNotExist          = "token不存在"
	UserNotExist           = "用户不存在"
	UserNotExistInDataBase = "用户未同步到数据库中"
)

const (
	LoginPlease = "请先进行登录"
)

const (
	PolicyNotExist   = "策略不存在"
	PolicyCountError = "策略使用次数添加失败"
)
