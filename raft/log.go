package raft

import (
	"fmt"
	"log"
)

type raftLog struct {
	// 实际上就是前面介绍的MemoryStorage实例 ，其中存储了快照数据及该快照之后的Entry记录
	storage Storage

	// unstable contains all unstable entries and snapshot.
	// they will be saved into storage.
	// 用于存储未写入Storage的快照数据及Entry记录
	unstable unstable

	// committed is the highest log position that is known to be in
	// stable storage on a quorum of nodes.
	// 已提交的位置，即已提交的Entry记录中最大的索引值
	committed uint64

	// applied is the highest log position that the application has
	// been instructed to apply to its state machine.
	// Invariant: applied <= committed
	// 已应用的位置，即已应用的Entry记录中最大的索引值。
	// 其中committed和applied之间始终满足applied<=committed这个不等式关系。
	applied uint64

	logger Logger
}

// newLog returns log using the given storage. It recovers the log to the state
// that it just commits and applies the latest snapshot.
func newLog(storage Storage, logger Logger) *raftLog {
	if storage == nil {
		log.Panic("storage must not be nil")
	}
	log := &raftLog{ //创建raftLog实例，并初始化storage字段
		storage: storage,
		logger:  logger,
	}
	//获取Storage中的第一条Entry和最后一条Entry的索引值
	firstIndex, err := storage.FirstIndex()
	if err != nil {
		panic(err) // TODO(bdarnell)
	}
	lastIndex, err := storage.LastIndex()
	if err != nil {
		panic(err) // TODO(bdarnell)
	}
	log.unstable.offset = lastIndex + 1 // 初始化unstable.offset
	log.unstable.logger = logger
	// Initialize our committed and applied pointers to the time of the last compaction.
	//初始化committed、applied字段
	log.committed = firstIndex - 1
	log.applied = firstIndex - 1

	return log
}

func (l *raftLog) String() string {
	return fmt.Sprintf("committed=%d, applied=%d, unstable.offset=%d, len(unstable.Entries)=%d", l.committed, l.applied, l.unstable.offset, len(l.unstable.entries))
}
