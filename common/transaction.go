package common

type Transaction struct {
	GlobalId   string           `json:"globalId"`   //全局事务编号
	BranchId   string           `json:"branchId"`   //分支事务编号
	CreateTime string           `json:"createTime"` //创建时间
	UpdateTime string           `json:"updateTime"` //最后更新时间
	Model      TransactionModel `json:"model"`      //事务模式
	Source     string           `json:"source"`     //原服务
	SourceIp   string           `json:"sourceIp"`   //原服务地址
}
