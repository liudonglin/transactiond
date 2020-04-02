package tcc

import "transactiond/common"

type TransactionTcc struct {
	common.TransactionModel
	Status       TransactionStatus      `json:"status"` //事务状态
	RetriedCount int32                  `json:"status"` //重试次数
	Type         common.TransactionType `json:"type"`   //事务类型
}
