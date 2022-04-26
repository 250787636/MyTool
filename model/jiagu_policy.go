package model

type JiaguPolicy struct {
	Id          int    `json:"id"`            //策略id
	Status      int    `json:"status"`        // 策略状态
	Name        string `json:"name"`          // 策略名称
	ConfigJson  string `json:"config_json"`   // 策略配置
	CustId      int    `json:"cust_id"`       // 所属客户id
	UserIds     string `json:"user_ids"`      // 关联的用户id
	NumberOfUse int    `json:"number_of_use"` // 策略使用次数
	//AppIds       string `json:"app_ids"`    // 关联的应用id
}
