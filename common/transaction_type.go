package common

type TransactionType int32

const (
	TransactionType_Root   TransactionType = 1 //tcc模式
	TransactionType_Branch TransactionType = 2 //lcn模式
)
