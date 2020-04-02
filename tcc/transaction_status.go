package tcc

type TransactionStatus int32

const (
	TransactionStatus_Trying     TransactionStatus = 1 //尝试中状态
	TransactionStatus_Confirming TransactionStatus = 2 //确认中状态
	TransactionStatus_Cancelling TransactionStatus = 3 //取消中状态
)
