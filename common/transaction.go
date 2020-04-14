package common

type Transaction struct {
	GlobalId   string           `json:"globalId"`   //全局事务编号
	UnitId     string           `json:"unitId"`     //分支事务编号
	CreateTime string           `json:"createTime"` //创建时间
	Model      TransactionModel `json:"model"`      //事务模式
	Source     string           `json:"source"`     //原服务
}
