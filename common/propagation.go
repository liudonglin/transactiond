package common

/*
Spring事务传播级别
不能要事务的
PROPAGATION_NEVER:			没有就非事务执行，有就抛出异常
PROPAGATION_NOT_SUPPORTED:	没有就非事务执行，有就直接挂起，然后非事务执行

可有可无的
PROPAGATION_SUPPORTS:		如果没有，就以非事务方式执行；如果有，就使用当前事务

必须有事务的
PROPAGATION_REQUIRES_NEW: 	有没有都新建事务，如果原来有，就将原来的挂起
PROPAGATION_NESTED: 		如果没有，就新建一个事务;如果有,就在当前事务中嵌套其他事务(会新建事务嵌套进上级事务)
PROPAGATION_REQUIRED:		如果没有,就新建一个事务; 如果有，就加入当前事务
PROPAGATION_MANDATORY: 		如果没有，就抛出异常;如果有，就使用当前事务
*/

type PropagationType int32

const (
	PropagationRequired    PropagationType = 0 //如果没有,就新建一个事务; 如果有，就加入当前事务
	PropagationSupports    PropagationType = 1 //如果没有，就以非事务方式执行；如果有，就使用当前事务
	PropagationMandatory   PropagationType = 2 //如果没有，就抛出异常;如果有，就使用当前事务。
	PropagationRequiresNew PropagationType = 3 //有没有都新建事务，如果原来有，就将原来的挂起。
)

func (p PropagationType) String() string {
	switch p {
	case PropagationRequired:
		return "Required"
	case PropagationSupports:
		return "Supports"
	case PropagationMandatory:
		return "Mandatory"
	case PropagationRequiresNew:
		return "RequiresNew"
	default:
		return ""
	}
}
