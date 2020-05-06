package raft

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
	pb "transactiond/raft/raftpb"
)

// None is a placeholder node ID used when there is no leader.
const None uint64 = 0
const noLimit = math.MaxUint64

// Possible values for StateType.
const (
	StateFollower StateType = iota
	StateCandidate
	StateLeader
	StatePreCandidate
	numStates
)

type ReadOnlyOption int

const (
	// ReadOnlySafe guarantees the linearizability of the read only request by
	// communicating with the quorum. It is the default and suggested option.
	ReadOnlySafe ReadOnlyOption = iota
	// ReadOnlyLeaseBased ensures linearizability of the read only request by
	// relying on the leader lease. It can be affected by clock drift.
	// If the clock drift is unbounded, leader might keep the lease longer than it
	// should (clock can move backward/pause without any bound). ReadIndex is not safe
	// in that case.
	ReadOnlyLeaseBased
)

// Possible values for CampaignType
const (
	// campaignPreElection represents the first phase of a normal election when
	// Config.PreVote is true.
	campaignPreElection CampaignType = "CampaignPreElection"
	// campaignElection represents a normal (time-based) election (the second phase
	// of the election when Config.PreVote is true).
	campaignElection CampaignType = "CampaignElection"
	// campaignTransfer represents the type of leader transfer
	campaignTransfer CampaignType = "CampaignTransfer"
)

// lockedRand is a small wrapper around rand.Rand to provide
// synchronization. Only the methods needed by the code are exposed
// (e.g. Intn).
type lockedRand struct {
	mu   sync.Mutex
	rand *rand.Rand
}

func (r *lockedRand) Intn(n int) int {
	r.mu.Lock()
	v := r.rand.Intn(n)
	r.mu.Unlock()
	return v
}

var globalRand = &lockedRand{
	rand: rand.New(rand.NewSource(time.Now().UnixNano())),
}

// CampaignType represents the type of campaigning
// the reason we use the type of string instead of uint64
// is because it's simpler to compare and fill in raft entries
type CampaignType string

// StateType represents the role of a node in a cluster.
type StateType uint64

var stmap = [...]string{
	"StateFollower",
	"StateCandidate",
	"StateLeader",
	"StatePreCandidate",
}

func (st StateType) String() string {
	return stmap[uint64(st)]
}

// Config contains the parameters to start a raft.
type Config struct {
	// ID is the identity of the local raft. ID cannot be 0.
	// 当前节点id
	ID uint64

	// peers contains the IDs of all nodes (including self) in the raft cluster. It
	// should only be set when starting a new raft cluster. Restarting raft from
	// previous configuration will panic if peers is set. peer is private and only
	// used for testing right now.
	// 记录了集群中所有结点的id
	peers []uint64

	// ElectionTick is the number of Node.Tick invocations that must pass between
	// elections. That is, if a follower does not receive any message from the
	// leader of current term before ElectionTick has elapsed, it will become
	// candidate and start an election. ElectionTick must be greater than
	// HeartbeatTick. We suggest ElectionTick = 10 * HeartbeatTick to avoid
	// unnecessary leader switching.
	// 用于初始化raft.electionTimeout,即逻辑时钟连续推进多少次后，就会触发Follower节点的状态切换及新一轮的Learder选举
	ElectionTick int
	// HeartbeatTick is the number of Node.Tick invocations that must pass between
	// heartbeats. That is, a leader sends heartbeat messages to maintain its
	// leadership every HeartbeatTick ticks.
	// 用于初始化raft.heartbeatTimeout， 即逻辑时钟连续推进多
	// 少次后，就触发Leader节点发送心跳消息。
	HeartbeatTick int

	// Storage is the storage for raft. raft generates entries and states to be
	// stored in storage. raft reads the persisted entries and states out of
	// Storage when it needs. raft reads out the previous state and configuration
	// out of storage when restarting.
	// 当前节点保存raft日志记录使用的存储
	Storage Storage
	// Applied is the last applied index. It should only be set when restarting
	// raft. raft will not return entries to the application smaller or equal to
	// Applied. If Applied is unset when restarting, raft might return previous
	// applied entries. This is a very application dependent configuration.
	// 当前已经应用的记录位置(已应用的最后一条Entry记录的索引值)，
	// 该值在节点重启时需要设置，否则会重新应用已经应用过Entry记录。
	Applied uint64

	// MaxSizePerMsg limits the max size of each append message. Smaller value
	// lowers the raft recovery cost(initial probing and message lost during normal
	// operation). On the other side, it might affect the throughput during normal
	// replication. Note: math.MaxUint64 for unlimited, 0 for at most one entry per
	// message.
	// 用于初始化raft.maxMsgSize字段，每条消息的最大字节数，
	// 如果是math.MaxUint64则没有上限，如果是0则表示每条消息最多携带一条Entry。
	MaxSizePerMsg uint64
	// MaxInflightMsgs limits the max number of in-flight append messages during
	// optimistic replication phase. The application transportation layer usually
	// has its own sending buffer over TCP/UDP. Setting MaxInflightMsgs to avoid
	// overflowing that sending buffer. TODO (xiangli): feedback to application to
	// limit the proposal rate?
	// 用于初始化raft.maxInflight,即已经发送出去且未收到响应的最大消息个数。
	MaxInflightMsgs int

	// CheckQuorum specifies if the leader should check quorum activity. Leader
	// steps down when quorum is not active for an electionTimeout.

	// 是否开启CheckQuorum模式，用于初始化raf.checkQuorum字段
	CheckQuorum bool

	// PreVote enables the Pre-Vote algorithm described in raft thesis section
	// 9.6. This prevents disruption when a node that has been partitioned away
	// rejoins the cluster.
	// 是否开启PreVote模式，用于初始化raf.preVote字段
	PreVote bool

	// ReadOnlyOption specifies how the read only request is processed.
	//
	// ReadOnlySafe guarantees the linearizability of the read only request by
	// communicating with the quorum. It is the default and suggested option.
	//
	// ReadOnlyLeaseBased ensures linearizability of the read only request by
	// relying on the leader lease. It can be affected by clock drift.
	// If the clock drift is unbounded, leader might keep the lease longer than it
	// should (clock can move backward/pause without any bound). ReadIndex is not safe
	// in that case.
	// 与只读请求的处理相关
	ReadOnlyOption ReadOnlyOption

	// Logger is the logger used for raft log. For multinode which can host
	// multiple raft group, each raft group can have its own logger
	Logger Logger
}

func (c *Config) validate() error {
	if c.ID == None {
		return errors.New("cannot use none as id")
	}

	if c.HeartbeatTick <= 0 {
		return errors.New("heartbeat tick must be greater than 0")
	}

	if c.ElectionTick <= c.HeartbeatTick {
		return errors.New("election tick must be greater than heartbeat tick")
	}

	if c.Storage == nil {
		return errors.New("storage cannot be nil")
	}

	if c.MaxInflightMsgs <= 0 {
		return errors.New("max inflight messages must be greater than 0")
	}

	if c.Logger == nil {
		c.Logger = raftLogger
	}

	return nil
}

type raft struct {

	// 当前节点在集群中的ID。
	id uint64

	// 当前任期号
	Term uint64

	// 当前任期中当前节点将选票投给了哪个节点，未投票时，该字段为None。
	Vote uint64

	// 与只读请求相关
	readStates []ReadState

	raftLog *raftLog

	maxInflight int

	// 单条消息的最大字节数。
	maxMsgSize uint64

	// Leader节点会记录集群中其他节点的日志复制情况Nextlndex和Matchlndex)
	prs map[uint64]*Progress

	// 当前节点在集群中的角色，可选值分为StateFollower, StateCandidate,StateLeader 和StatePreCandidate四种状态。
	state StateType

	// 在选举过程，如果当前节点收到了来自某个节点的投票，则会将votes中对应的值设置为true,
	// 通过统计votes 这个map,就可以确定当前节点收到的投票是否超过半数。
	votes map[uint64]bool

	// 缓存了当前节点等待发送的消息。
	msgs []pb.Message

	// the leader id
	// 当前集群中Leader节点的ID。
	lead uint64
	// leadTransferee is id of the leader transfer target when its value is not zero.
	// Follow the procedure defined in raft thesis 3.10.
	// 用于集群中Leader节点的转移，leadTransferee记录了此次Leader角色转移的目标节点的ID。
	leadTransferee uint64
	// New configuration is ignored if there exists unapplied configuration.
	pendingConf bool

	// 与只读请求相关
	readOnly *readOnly

	// number of ticks since it reached last electionTimeout when it is leader
	// or candidate.
	// number of ticks since it reached last electionTimeout or received a
	// valid message from current leader when it is a follower.
	// 选举计时器的指针，其单位是逻辑时钟的刻度，逻辑时钟每推进一次，该字段值就会增加1。
	electionElapsed int

	// number of ticks since it reached last heartbeatTimeout.
	// only leader keeps heartbeatElapsed.
	// 心跳计时器的指针，其单位也是逻辑时钟的刻度，逻辑时钟每推进一次，该字段值就会增加1。
	heartbeatElapsed int

	// 每隔一段时间，Leader节点会尝试连接集群中的其他节点(发送心跳消息)，如果发现自己可以
	// 连接到节点个数没有超过半数(即没有收到足够的心跳响应)，则主动切换成Follower状态。
	checkQuorum bool

	preVote bool

	// 心跳超时时间
	heartbeatTimeout int

	// 选举超时时间
	electionTimeout int

	// randomizedElectionTimeout is a random number between
	// [electiontimeout, 2 * electiontimeout - 1]. It gets reset
	// when raft changes its state to follower or candidate.
	randomizedElectionTimeout int

	tick func()

	step stepFunc

	logger Logger
}

func newRaft(c *Config) *raft {
	//检测参数Config中各字段的合法性，异常处理
	if err := c.validate(); err != nil {
		panic(err.Error())
	}
	//创建raftLog实例，用于记录Entry记录
	raftlog := newLog(c.Storage, c.Logger)
	// 获取raftLog.storage的初始状态(HardState和ConfState)
	// Storage的初始状态是通过本地Entry记录回放得到的
	hs, cs, err := c.Storage.InitialState()
	if err != nil {
		panic(err) // TODO(bdarnell)
	}
	//根据快照中记录的集群信息和Config的配置信息初始化集群的节点信息
	peers := c.peers
	if len(cs.Nodes) > 0 {
		if len(peers) > 0 {
			// TODO(bdarnell): the peers argument is always nil except in
			// tests; the argument should be removed and these tests should be
			// updated to specify their nodes through a snapshot.
			panic("cannot specify both newRaft(peers) and ConfState.Nodes)")
		}
		peers = cs.Nodes
	}
	//创建raft实例
	r := &raft{
		id:      c.ID,    //当前节点的ID
		lead:    None,    //当前集群中Leader节点的ID，初始化时先被设立成0
		raftLog: raftlog, //负责管理Entry记录的raftLog实例
		//每条消息的最大字节数，如果是math.MaxUint64则没有上限，如果是0则表示每条消息最多携带一条Entry
		maxMsgSize:       c.MaxSizePerMsg,
		maxInflight:      c.MaxInflightMsgs, // 已经发送出去且未收到响应的最大消息个数
		prs:              make(map[uint64]*Progress),
		electionTimeout:  c.ElectionTick,                //选举超时时间
		heartbeatTimeout: c.HeartbeatTick,               //心跳超时时间
		logger:           c.Logger,                      //普通日志输出
		checkQuorum:      c.CheckQuorum,                 //是否开启CheckQuorum模式
		preVote:          c.PreVote,                     //是否开启PreVote模式
		readOnly:         newReadOnly(c.ReadOnlyOption), //只读请求的相关配置
	}

	for _, p := range peers {
		r.prs[p] = &Progress{Next: 1, ins: newInflights(r.maxInflight)}
	}
	//根据从Storage中获取的HardState，初始化raftLog.committed字段，以及raft.Term和Vote字段
	if !isHardStateEqual(hs, emptyState) {
		r.loadState(hs)
	}

	if c.Applied > 0 {
		raftlog.appliedTo(c.Applied)
	}

	//当前节点切换成Follower状态
	r.becomeFollower(r.Term, None)

	var nodesStrs []string
	for _, n := range r.nodes() {
		nodesStrs = append(nodesStrs, fmt.Sprintf("%x", n))
	}

	r.logger.Infof("newRaft %x [peers: [%s], term: %d, commit: %d, applied: %d, lastindex: %d, lastterm: %d]",
		r.id, strings.Join(nodesStrs, ","), r.Term, r.raftLog.committed, r.raftLog.applied, r.raftLog.lastIndex(), r.raftLog.lastTerm())
	return r
}

func (r *raft) hasLeader() bool { return r.lead != None }

func (r *raft) softState() *SoftState { return &SoftState{Lead: r.lead, RaftState: r.state} }

func (r *raft) hardState() pb.HardState {
	return pb.HardState{
		Term:   r.Term,
		Vote:   r.Vote,
		Commit: r.raftLog.committed,
	}
}

func (r *raft) quorum() int { return len(r.prs)/2 + 1 }

func (r *raft) nodes() []uint64 {
	nodes := make([]uint64, 0, len(r.prs))
	for id := range r.prs {
		nodes = append(nodes, id)
	}
	sort.Sort(uint64Slice(nodes))
	return nodes
}

// send persists state to stable storage and then sends to its mailbox.
func (r *raft) send(m pb.Message) {
	m.From = r.id //设置消息的发送节点ID,即当前节点ID
	if m.Type == pb.MsgVote || m.Type == pb.MsgPreVote {
		// 对MsgVote和 MsgPreVote消息的Term字段进行检查
		if m.Term == 0 { //对其他类型消息的Term字段位进行设置
			// PreVote RPCs are sent at a term other than our actual term, so the code
			// that sends these messages is responsible for setting the term.
			panic(fmt.Sprintf("term should be set when sending %s", m.Type))
		}
	} else {
		if m.Term != 0 {
			panic(fmt.Sprintf("term should not be set when sending %s (was %d)", m.Type, m.Term))
		}
		// do not attach term to MsgProp, MsgReadIndex
		// proposals are a way to forward to the leader and
		// should be treated as local message.
		// MsgReadIndex is also forwarded to leader.
		// 除了MsgProp和MsgReadindex两类消息(这两类消息的Term值为0即为本地消息)之外
		// 其他类型消息的Term字段值在这里统一设置
		if m.Type != pb.MsgProp && m.Type != pb.MsgReadIndex {
			m.Term = r.Term
		}
	}
	r.msgs = append(r.msgs, m) //将消息添加到r.msgs切片中等待发送
}

// sendAppend sends RPC, with entries to the given peer.
func (r *raft) sendAppend(to uint64) {
	pr := r.prs[to]
	if pr.IsPaused() {
		return
	}
	m := pb.Message{}
	m.To = to

	term, errt := r.raftLog.term(pr.Next - 1)
	ents, erre := r.raftLog.entries(pr.Next, r.maxMsgSize)

	if errt != nil || erre != nil { // send snapshot if we failed to get term or entries
		if !pr.RecentActive {
			r.logger.Debugf("ignore sending snapshot to %x since it is not recently active", to)
			return
		}

		m.Type = pb.MsgSnap
		snapshot, err := r.raftLog.snapshot()
		if err != nil {
			if err == ErrSnapshotTemporarilyUnavailable {
				r.logger.Debugf("%x failed to send snapshot to %x because snapshot is temporarily unavailable", r.id, to)
				return
			}
			panic(err) // TODO(bdarnell)
		}
		if IsEmptySnap(snapshot) {
			panic("need non-empty snapshot")
		}
		m.Snapshot = snapshot
		sindex, sterm := snapshot.Metadata.Index, snapshot.Metadata.Term
		r.logger.Debugf("%x [firstindex: %d, commit: %d] sent snapshot[index: %d, term: %d] to %x [%s]",
			r.id, r.raftLog.firstIndex(), r.raftLog.committed, sindex, sterm, to, pr)
		pr.becomeSnapshot(sindex)
		r.logger.Debugf("%x paused sending replication messages to %x [%s]", r.id, to, pr)
	} else {
		m.Type = pb.MsgApp
		m.Index = pr.Next - 1
		m.LogTerm = term
		m.Entries = ents
		m.Commit = r.raftLog.committed
		if n := len(m.Entries); n != 0 {
			switch pr.State {
			// optimistically increase the next when in ProgressStateReplicate
			case ProgressStateReplicate:
				last := m.Entries[n-1].Index
				pr.optimisticUpdate(last)
				pr.ins.add(last)
			case ProgressStateProbe:
				pr.pause()
			default:
				r.logger.Panicf("%x is sending append in unhandled state %s", r.id, pr.State)
			}
		}
	}
	r.send(m)
}

// sendHeartbeat sends an empty MsgApp
func (r *raft) sendHeartbeat(to uint64, ctx []byte) {
	// Attach the commit as min(to.matched, r.committed).
	// When the leader sends out heartbeat message,
	// the receiver(follower) might not be matched with the leader
	// or it might not have all the committed entries.
	// The leader MUST NOT forward the follower's commit to
	// an unmatched index.
	commit := min(r.prs[to].Match, r.raftLog.committed)
	m := pb.Message{
		To:      to,
		Type:    pb.MsgHeartbeat,
		Commit:  commit,
		Context: ctx,
	}

	r.send(m)
}

// bcastAppend sends RPC, with entries to all peers that are not up-to-date
// according to the progress recorded in r.prs.
func (r *raft) bcastAppend() {
	for id := range r.prs {
		if id == r.id {
			continue
		}
		r.sendAppend(id)
	}
}

// bcastHeartbeat sends RPC, without entries to all the peers.
func (r *raft) bcastHeartbeat() {
	lastCtx := r.readOnly.lastPendingRequestCtx()
	if len(lastCtx) == 0 {
		r.bcastHeartbeatWithCtx(nil)
	} else {
		r.bcastHeartbeatWithCtx([]byte(lastCtx))
	}
}

func (r *raft) bcastHeartbeatWithCtx(ctx []byte) {
	for id := range r.prs {
		if id == r.id {
			continue
		}
		r.sendHeartbeat(id, ctx)
	}
}

// maybeCommit attempts to advance the commit index. Returns true if
// the commit index changed (in which case the caller should call
// r.bcastAppend).
func (r *raft) maybeCommit() bool {
	// TODO(bmizerany): optimize.. Currently naive
	//将集群中所有节点对应的Progress.Match字段复制到mis切片中
	mis := make(uint64Slice, 0, len(r.prs))
	for id := range r.prs {
		mis = append(mis, r.prs[id].Match)
	}
	sort.Sort(sort.Reverse(mis)) // 对这些Match位进行排序
	mci := mis[r.quorum()-1]
	return r.raftLog.maybeCommit(mci, r.Term) //更新raftLog.committed字段,完成提交
}

func (r *raft) reset(term uint64) {
	if r.Term != term {
		r.Term = term //重置Term字段
		r.Vote = None //重置Vote字段
	}
	r.lead = None //清空lead字段

	//重置选举计时器和心跳计时器
	r.electionElapsed = 0
	r.heartbeatElapsed = 0
	//重置选举计时器的过期时间(随机值)
	r.resetRandomizedElectionTimeout()

	r.abortLeaderTransfer() //清空leadTransferee

	r.votes = make(map[uint64]bool) //重置votes字段

	//重置prs，其中每个Progress中的Next设置为raftLog.lastindex
	for id := range r.prs {
		r.prs[id] = &Progress{Next: r.raftLog.lastIndex() + 1, ins: newInflights(r.maxInflight)}
		if id == r.id { //将当前节点对应的prs.Match设置成lastIndex
			r.prs[id].Match = r.raftLog.lastIndex()
		}
	}
	r.pendingConf = false                       //清空pendingConf字段
	r.readOnly = newReadOnly(r.readOnly.option) //只读请求的相关设置
}

// 1.设置待追加的Entry记录的Term值和Index值。
// 2.向当前节点的raftLog中追加Entry记录。
// 3.更新当前节点对应的Progress实例。
// 4.尝试提交Entry记录，即修改raftLog.committed字段的值。
func (r *raft) appendEntry(es ...pb.Entry) {
	//获取raftLog中最后一条记录的索引值
	li := r.raftLog.lastIndex()
	for i := range es { //更新待追加记录的Term值和索引值
		es[i].Term = r.Term              //Entry记录的Term指定为当前Leader节点的任期号
		es[i].Index = li + 1 + uint64(i) //为日志记录指定index
	}
	r.raftLog.append(es...) //向raftLog中追加记录
	//更新当前节点对应的Progress，主要是更新Next和Match
	r.prs[r.id].maybeUpdate(r.raftLog.lastIndex())
	// Regardless of maybeCommit's return, our caller will call bcastAppend.
	r.maybeCommit() //尝试提交Entry记录
}

// tickElection is run by followers and candidates after r.electionTimeout.
func (r *raft) tickElection() {
	r.electionElapsed++ //递增electionElapsed

	//promotable()方法会检查prs字段中是否还存在当前节点对应的Progress实例 ，这是为了检测当前节点是否被从集群中移除了
	//pastElectionTimeout()方法是检测当前的选举计时器是否超时(r.electionElapsed>=r.randomizedElectionTimeout)
	if r.promotable() && r.pastElectionTimeout() {
		r.electionElapsed = 0                           //重置electionElapsed
		r.Step(pb.Message{From: r.id, Type: pb.MsgHup}) //发起选举
	}
}

// tickHeartbeat is run by leaders to send a MsgBeat after r.heartbeatTimeout.
func (r *raft) tickHeartbeat() {
	r.heartbeatElapsed++
	r.electionElapsed++

	if r.electionElapsed >= r.electionTimeout {
		r.electionElapsed = 0
		if r.checkQuorum {
			r.Step(pb.Message{From: r.id, Type: pb.MsgCheckQuorum})
		}
		// If current leader cannot transfer leadership in electionTimeout, it becomes leader again.
		if r.state == StateLeader && r.leadTransferee != None {
			r.abortLeaderTransfer()
		}
	}

	if r.state != StateLeader {
		return
	}

	if r.heartbeatElapsed >= r.heartbeatTimeout {
		r.heartbeatElapsed = 0
		r.Step(pb.Message{From: r.id, Type: pb.MsgBeat})
	}
}

func (r *raft) becomeFollower(term uint64, lead uint64) {
	//将step字段设置成stepFollower, stepFollower()函数中封装了Follower节点处理消息的行为
	r.step = stepFollower
	r.reset(term)           //重置raft实例Term,Vote等字段
	r.tick = r.tickElection //将tick字段设立成tickElection函数
	r.lead = lead           //设置当前集群的Leader节点
	r.state = StateFollower //设立当前节点的角色
	r.logger.Infof("%x became follower at term %d", r.id, r.Term)
}

// 当节点可以连接到集群中半数以上的节点时，会调用becomeCandidate()方法切换到Candidate状态，
// becomeCandidate()方法的具体实现如下:
func (r *raft) becomeCandidate() {
	//检测当前节点的状态，禁止直接从Leader状态切换到PreCandidate状态
	// TODO(xiangli) remove the panic when the raft implementation is stable
	if r.state == StateLeader {
		panic("invalid transition [leader -> candidate]")
	}
	//将step字段设置成stepCandidate, stepCandidate()函数中封装了Candidate节点处理消息的处理行为
	r.step = stepCandidate
	//重置raft实例Term、Vote等字段，注意这里的Term已经递增
	r.reset(r.Term + 1)
	r.tick = r.tickElection
	r.Vote = r.id            //在此次选举中，Candidate节点会将选票投给自己
	r.state = StateCandidate //修改当前节点的角色
	r.logger.Infof("%x became candidate at term %d", r.id, r.Term)
}

// 如果当前集群开启了PreVote模式，当Follower节点的选举计时器超时时，会先调用becomePreCandidate()方法切换到PreCandidate状态，
func (r *raft) becomePreCandidate() {
	//检测当前节点的状态，禁止直接从Leader状态切换到PreCandidate状态
	// TODO(xiangli) remove the panic when the raft implementation is stable
	if r.state == StateLeader {
		panic("invalid transition [leader -> pre-candidate]")
	}
	// Becoming a pre-candidate changes our step functions and state,
	// but doesn't change anything else. In particular it does not increase
	// r.Term or change r.Vote.
	//将step字段设置成stepCandidate, stepCandidate()函数中封装了PreCandidate节点处理消息的处理行为
	r.step = stepCandidate
	r.tick = r.tickElection
	r.state = StatePreCandidate
	r.logger.Infof("%x became pre-candidate at term %d", r.id, r.Term)
}

// 当Candidate节点得到集群中半数以上节点的选票时，会调用becomeLeader()方法切换成Leader状态，
func (r *raft) becomeLeader() {
	//检测当前节点的状态，禁止从Follower状态切换成Leader状态
	// TODO(xiangli) remove the panic when the raft implementation is stable
	if r.state == StateFollower {
		panic("invalid transition [follower -> leader]")
	}
	// 将step字段设置成stepLeader, stepLeader()函数中封装了Leader节点处理消息的行为
	r.step = stepLeader
	r.reset(r.Term)
	r.tick = r.tickHeartbeat
	r.lead = r.id         //将leader字段当前节点的id
	r.state = StateLeader //更新当前节点的角色
	// 获取当前节点中所有未提交的Entry记录
	ents, err := r.raftLog.entries(r.raftLog.committed+1, noLimit)
	if err != nil {
		r.logger.Panicf("unexpected error getting uncommitted entries (%v)", err)
	}

	// 检测未提交的记录中是否存在多条集群配置变更的Entry记录(即EntryConfChange类型的Entry记录)
	// 如果存在多条 EntryConfChange类型记录，则异常关闭
	nconf := numOfPendingConf(ents)
	if nconf > 1 {
		panic("unexpected multiple uncommitted config entry")
	}

	if nconf == 1 {
		r.pendingConf = true
	}

	// 向当前节点的raftLog中追加一条空的Entry记录
	r.appendEntry(pb.Entry{Data: nil})
	r.logger.Infof("%x became leader at term %d", r.id, r.Term)
}

// 在raft.campaign()方法中,除了完成状态切换,还会向集群中的其他节点发送相应类型的消
// 例如,如果当前Follower节点要切换成PreCandidate状态,则会发送MsgPreVote消息。
// raft.campaign()方法的具体实现如下 :
func (r *raft) campaign(t CampaignType) {
	//在该方法的最后,会发送一条消息,这里term和voteMsg分别用于确定该消息的Term值和类型
	var term uint64
	var voteMsg pb.MessageType
	if t == campaignPreElection { //切才是的目标状态是PreCandidate
		r.becomePreCandidate()  //将当前节点切换成PreCandidate状态
		voteMsg = pb.MsgPreVote //确定最后发送的消息是MsgPreVote类型
		// PreVote RPCs are sent for the next term before we've incremented r.Term.
		//确定最后发送消息的Term值，这里只是增加了消息的Term值，并未增加raft.Term字段的值
		term = r.Term + 1
	} else { //切换的目标状态是Candidate
		// 将当前节点切换成Candidate状态，becomeCandidate()方法中会
		// 增加raft.Term字段的值，并将当前节点的选票投给自身
		r.becomeCandidate()
		voteMsg = pb.MsgVote //确定最后发送的消息是MsgPreVote类型
		term = r.Term        //确定最后发送消息的Term值
	}
	//统计当前节点收到的选票，并统计其得票数是否超过半数，这次检测主要是为单节点设置的
	if r.quorum() == r.poll(r.id, voteRespMsgType(voteMsg), true) {
		// We won the election after voting for ourselves (which must mean that
		// this is a single-node cluster). Advance to the next state.
		// 当得到足够的选票时，则将PreCandidate状态的节点切换成Candidate状态
		// Candidate状态的节点则切换成Leader状态
		if t == campaignPreElection {
			r.campaign(campaignElection)
		} else {
			r.becomeLeader()
		}
		return
	}
	for id := range r.prs { //状态切换完成之后，当前节点会向集群中所有节点发送指定类型的消息
		if id == r.id { //跳过当前节点自身
			continue
		}
		r.logger.Infof("%x [logterm: %d, index: %d] sent %s request to %x at term %d",
			r.id, r.raftLog.lastTerm(), r.raftLog.lastIndex(), voteMsg, id, r.Term)

		var ctx []byte
		//在进行Leader节点转移时，MsgPreVote或MsgVote消息会在Context字段中设立该特殊值
		if t == campaignTransfer {
			ctx = []byte(t)
		}
		// 发送指定类型的消息，其中Index和LogTerm分别是当前节点的raftLog中
		// 最后一条消息的Index值和Term值
		r.send(pb.Message{Term: term, To: id, Type: voteMsg, Index: r.raftLog.lastIndex(), LogTerm: r.raftLog.lastTerm(), Context: ctx})
	}
}

func (r *raft) poll(id uint64, t pb.MessageType, v bool) (granted int) {
	if v {
		r.logger.Infof("%x received %s from %x at term %d", r.id, t, id, r.Term)
	} else {
		r.logger.Infof("%x received %s rejection from %x at term %d", r.id, t, id, r.Term)
	}
	if _, ok := r.votes[id]; !ok {
		r.votes[id] = v
	}
	for _, vv := range r.votes {
		if vv {
			granted++
		}
	}
	return granted
}

// raft.Step()方法是etcd-raft模块处理各类消息的入口
// raft.Step()方法主要分为两部分:
// 第一部分是根据Term值对消息进行分类处理
// 第二部分是根据消息的类型进行分类处理
func (r *raft) Step(m pb.Message) error {
	// Handle the message term, which may result in our stepping down to a follower.
	//首先根据消息的Term值进行分类处理
	switch {
	case m.Term == 0:
		// 对本地消息并没有做什么处理，这里介绍的MsgHup消息Term值为0,就是本地消息的一种;
		// 后面介绍的MsgProp消息和MsgReadindex消息也是本地消息，其Term值也是0
		// local message
	case m.Term > r.Term:
		lead := m.From                                       //记录消息携带的Term值
		if m.Type == pb.MsgVote || m.Type == pb.MsgPreVote { // 这里只对MsgVote和MsgPreVote两种类型消息进行处理
			//根据消息的Context字段判断收到的MsgPreVote(或MsgVote)消息是否为Leader
			//节点转移场景下产生的，如果是，则强制当前节点参与本次预选(或选举)
			force := bytes.Equal(m.Context, []byte(campaignTransfer))
			// 下面通过一系列条件判断当前节点是否参与此次选举(或预选)，其中主要检测集群是否开启
			// 了CheckQuorum模式、当前节点是否有已知的Lead节点，以及其选举计时器的时间
			inLease := r.checkQuorum && r.lead != None && r.electionElapsed < r.electionTimeout
			if !force && inLease {
				// If a server receives a RequestVote request within the minimum election timeout
				// of hearing from a current leader, it does not update its term or grant its vote
				r.logger.Infof("%x [logterm: %d, index: %d, vote: %x] ignored %s from %x [logterm: %d, index: %d] at term %d: lease is not expired (remaining ticks: %d)",
					r.id, r.raftLog.lastTerm(), r.raftLog.lastIndex(), r.Vote, m.Type, m.From, m.LogTerm, m.Index, r.Term, r.electionTimeout-r.electionElapsed)
				return nil //当前节点不参与此次选举
			}
			lead = None
		}
		switch { //在这个switch中，当前节点会根据消息类型决定是否切换状态
		case m.Type == pb.MsgPreVote: //收到MsgPreVote消息时，不会引起当前节点的状态切换
			// Never change our term in response to a PreVote
		case m.Type == pb.MsgPreVoteResp && !m.Reject:
			// We send pre-vote requests with a term in our future. If the
			// pre-vote is granted, we will increment our term when we get a
			// quorum. If it is not, the term comes from the node that
			// rejected our vote so we should become a follower at the new
			// term.
		default:
			r.logger.Infof("%x [term: %d] received a %s message with higher term from %x [term: %d]",
				r.id, r.Term, m.Type, m.From, m.Term)
			r.becomeFollower(m.Term, lead)
		}

	case m.Term < r.Term:
		if r.checkQuorum && (m.Type == pb.MsgHeartbeat || m.Type == pb.MsgApp) {
			// We have received messages from a leader at a lower term. It is possible
			// that these messages were simply delayed in the network, but this could
			// also mean that this node has advanced its term number during a network
			// partition, and it is now unable to either win an election or to rejoin
			// the majority on the old term. If checkQuorum is false, this will be
			// handled by incrementing term numbers in response to MsgVote with a
			// higher term, but if checkQuorum is true we may not advance the term on
			// MsgVote and must generate other messages to advance the term. The net
			// result of these two features is to minimize the disruption caused by
			// nodes that have been removed from the cluster's configuration: a
			// removed node will send MsgVotes (or MsgPreVotes) which will be ignored,
			// but it will not receive MsgApp or MsgHeartbeat, so it will not create
			// disruptive term increases
			r.send(pb.Message{To: m.From, Type: pb.MsgAppResp})
		} else {
			// ignore other cases
			r.logger.Infof("%x [term: %d] ignored a %s message with lower term from %x [term: %d]",
				r.id, r.Term, m.Type, m.From, m.Term)
		}
		return nil
	}

	switch m.Type { //根据Message的Type进行分类处理
	case pb.MsgHup: //这里针对MsgHup类型的消息进行处理
		if r.state != StateLeader { //只有非Leader状态的节点才会处理MsgHup消息
			//获取raftLog中已提交但未应用(即applied~committed)的Entry记录
			ents, err := r.raftLog.slice(r.raftLog.applied+1, r.raftLog.committed+1, noLimit)
			if err != nil {
				r.logger.Panicf("unexpected error getting unapplied entries (%v)", err)
			}
			//检测是否有未应用的EntryConfChange记录，如果有就放弃发起选举的机会
			if n := numOfPendingConf(ents); n != 0 && r.raftLog.committed > r.raftLog.applied {
				r.logger.Warningf("%x cannot campaign at term %d since there are still %d pending configuration changes to apply", r.id, r.Term, n)
				return nil
			}

			r.logger.Infof("%x is starting a new election at term %d", r.id, r.Term)
			if r.preVote {
				// 检测当前集群是否开启了PreVote模式，如果开启了，则允切换到
				// 调用raft.campaign()方法切换当前节点的角色，发起PreVote
				r.campaign(campaignPreElection)
			} else {
				r.campaign(campaignElection)
			}
		} else { //如果当前节点已经是Leader状态，则仅仅输出一条Debug日志
			r.logger.Debugf("%x ignoring MsgHup because already leader", r.id)
		}

	case pb.MsgVote, pb.MsgPreVote: //对MsgVote和MsgPreVote消息的处理
		// The m.Term > r.Term clause is for MsgPreVote. For MsgVote m.Term should
		// always equal r.Term.
		// 当前节点在参与预选时，会综合下面几个条件决定是否投票(在Raft协议的介绍中也提到过)
		// 1、当前节点是否已经投过票
		// 2、MsgPreVote消息发送者的任期号是否更大
		// 3、当前节点投票给了对方节点
		// 4、MsgPreVote消息发送者的raftLog中是否包含当前节点的全部Entry记录
		if (r.Vote == None || m.Term > r.Term || r.Vote == m.From) && r.raftLog.isUpToDate(m.Index, m.LogTerm) {
			r.logger.Infof("%x [logterm: %d, index: %d, vote: %x] cast %s for %x [logterm: %d, index: %d] at term %d",
				r.id, r.raftLog.lastTerm(), r.raftLog.lastIndex(), r.Vote, m.Type, m.From, m.LogTerm, m.Index, r.Term)
			// 将票投给MsgPreVote消息的发送节点
			r.send(pb.Message{To: m.From, Type: voteRespMsgType(m.Type)})
			// 针对MsgVote消息的处理
			if m.Type == pb.MsgVote {
				// Only record real votes.
				r.electionElapsed = 0
				r.Vote = m.From
			}
		} else {
			r.logger.Infof("%x [logterm: %d, index: %d, vote: %x] rejected %s from %x [logterm: %d, index: %d] at term %d",
				r.id, r.raftLog.lastTerm(), r.raftLog.lastIndex(), r.Vote, m.Type, m.From, m.LogTerm, m.Index, r.Term)
			// 不满足上述投赞同票条件时，当前节点会返回拒绝票(即响应消息中的Reject字段会设立成true)
			r.send(pb.Message{To: m.From, Type: voteRespMsgType(m.Type), Reject: true})
		}

	default:
		r.step(r, m)
	}
	return nil
}

type stepFunc func(r *raft, m pb.Message)

func stepLeader(r *raft, m pb.Message) {
	// These message types do not require any progress for m.From.
	switch m.Type {
	case pb.MsgBeat:
		r.bcastHeartbeat()
		return
	case pb.MsgCheckQuorum:
		if !r.checkQuorumActive() {
			r.logger.Warningf("%x stepped down to follower since quorum is not active", r.id)
			r.becomeFollower(r.Term, None)
		}
		return
	case pb.MsgProp:
		if len(m.Entries) == 0 {
			r.logger.Panicf("%x stepped empty MsgProp", r.id)
		}
		if _, ok := r.prs[r.id]; !ok {
			// If we are not currently a member of the range (i.e. this node
			// was removed from the configuration while serving as leader),
			// drop any new proposals.
			return
		}
		if r.leadTransferee != None {
			r.logger.Debugf("%x [term %d] transfer leadership to %x is in progress; dropping proposal", r.id, r.Term, r.leadTransferee)
			return
		}

		for i, e := range m.Entries {
			if e.Type == pb.EntryConfChange {
				if r.pendingConf {
					r.logger.Infof("propose conf %s ignored since pending unapplied configuration", e.String())
					m.Entries[i] = pb.Entry{Type: pb.EntryNormal}
				}
				r.pendingConf = true
			}
		}
		r.appendEntry(m.Entries...)
		r.bcastAppend()
		return
	case pb.MsgReadIndex:
		if r.quorum() > 1 {
			if r.raftLog.zeroTermOnErrCompacted(r.raftLog.term(r.raftLog.committed)) != r.Term {
				// Reject read only request when this leader has not committed any log entry at its term.
				return
			}

			// thinking: use an interally defined context instead of the user given context.
			// We can express this in terms of the term and index instead of a user-supplied value.
			// This would allow multiple reads to piggyback on the same message.
			switch r.readOnly.option {
			case ReadOnlySafe:
				r.readOnly.addRequest(r.raftLog.committed, m)
				r.bcastHeartbeatWithCtx(m.Entries[0].Data)
			case ReadOnlyLeaseBased:
				var ri uint64
				if r.checkQuorum {
					ri = r.raftLog.committed
				}
				if m.From == None || m.From == r.id { // from local member
					r.readStates = append(r.readStates, ReadState{Index: r.raftLog.committed, RequestCtx: m.Entries[0].Data})
				} else {
					r.send(pb.Message{To: m.From, Type: pb.MsgReadIndexResp, Index: ri, Entries: m.Entries})
				}
			}
		} else {
			r.readStates = append(r.readStates, ReadState{Index: r.raftLog.committed, RequestCtx: m.Entries[0].Data})
		}

		return
	}

	// All other message types require a progress for m.From (pr).
	pr, prOk := r.prs[m.From]
	if !prOk {
		r.logger.Debugf("%x no progress available for %x", r.id, m.From)
		return
	}
	switch m.Type {
	case pb.MsgAppResp:
		pr.RecentActive = true

		if m.Reject {
			r.logger.Debugf("%x received msgApp rejection(lastindex: %d) from %x for index %d",
				r.id, m.RejectHint, m.From, m.Index)
			if pr.maybeDecrTo(m.Index, m.RejectHint) {
				r.logger.Debugf("%x decreased progress of %x to [%s]", r.id, m.From, pr)
				if pr.State == ProgressStateReplicate {
					pr.becomeProbe()
				}
				r.sendAppend(m.From)
			}
		} else {
			oldPaused := pr.IsPaused()
			if pr.maybeUpdate(m.Index) {
				switch {
				case pr.State == ProgressStateProbe:
					pr.becomeReplicate()
				case pr.State == ProgressStateSnapshot && pr.needSnapshotAbort():
					r.logger.Debugf("%x snapshot aborted, resumed sending replication messages to %x [%s]", r.id, m.From, pr)
					pr.becomeProbe()
				case pr.State == ProgressStateReplicate:
					pr.ins.freeTo(m.Index)
				}

				if r.maybeCommit() {
					r.bcastAppend()
				} else if oldPaused {
					// update() reset the wait state on this node. If we had delayed sending
					// an update before, send it now.
					r.sendAppend(m.From)
				}
				// Transfer leadership is in progress.
				if m.From == r.leadTransferee && pr.Match == r.raftLog.lastIndex() {
					r.logger.Infof("%x sent MsgTimeoutNow to %x after received MsgAppResp", r.id, m.From)
					r.sendTimeoutNow(m.From)
				}
			}
		}
	case pb.MsgHeartbeatResp:
		pr.RecentActive = true
		pr.resume()

		// free one slot for the full inflights window to allow progress.
		if pr.State == ProgressStateReplicate && pr.ins.full() {
			pr.ins.freeFirstOne()
		}
		if pr.Match < r.raftLog.lastIndex() {
			r.sendAppend(m.From)
		}

		if r.readOnly.option != ReadOnlySafe || len(m.Context) == 0 {
			return
		}

		ackCount := r.readOnly.recvAck(m)
		if ackCount < r.quorum() {
			return
		}

		rss := r.readOnly.advance(m)
		for _, rs := range rss {
			req := rs.req
			if req.From == None || req.From == r.id { // from local member
				r.readStates = append(r.readStates, ReadState{Index: rs.index, RequestCtx: req.Entries[0].Data})
			} else {
				r.send(pb.Message{To: req.From, Type: pb.MsgReadIndexResp, Index: rs.index, Entries: req.Entries})
			}
		}
	case pb.MsgSnapStatus:
		if pr.State != ProgressStateSnapshot {
			return
		}
		if !m.Reject {
			pr.becomeProbe()
			r.logger.Debugf("%x snapshot succeeded, resumed sending replication messages to %x [%s]", r.id, m.From, pr)
		} else {
			pr.snapshotFailure()
			pr.becomeProbe()
			r.logger.Debugf("%x snapshot failed, resumed sending replication messages to %x [%s]", r.id, m.From, pr)
		}
		// If snapshot finish, wait for the msgAppResp from the remote node before sending
		// out the next msgApp.
		// If snapshot failure, wait for a heartbeat interval before next try
		pr.pause()
	case pb.MsgUnreachable:
		// During optimistic replication, if the remote becomes unreachable,
		// there is huge probability that a MsgApp is lost.
		if pr.State == ProgressStateReplicate {
			pr.becomeProbe()
		}
		r.logger.Debugf("%x failed to send message to %x because it is unreachable [%s]", r.id, m.From, pr)
	case pb.MsgTransferLeader:
		leadTransferee := m.From
		lastLeadTransferee := r.leadTransferee
		if lastLeadTransferee != None {
			if lastLeadTransferee == leadTransferee {
				r.logger.Infof("%x [term %d] transfer leadership to %x is in progress, ignores request to same node %x",
					r.id, r.Term, leadTransferee, leadTransferee)
				return
			}
			r.abortLeaderTransfer()
			r.logger.Infof("%x [term %d] abort previous transferring leadership to %x", r.id, r.Term, lastLeadTransferee)
		}
		if leadTransferee == r.id {
			r.logger.Debugf("%x is already leader. Ignored transferring leadership to self", r.id)
			return
		}
		// Transfer leadership to third party.
		r.logger.Infof("%x [term %d] starts to transfer leadership to %x", r.id, r.Term, leadTransferee)
		// Transfer leadership should be finished in one electionTimeout, so reset r.electionElapsed.
		r.electionElapsed = 0
		r.leadTransferee = leadTransferee
		if pr.Match == r.raftLog.lastIndex() {
			r.sendTimeoutNow(leadTransferee)
			r.logger.Infof("%x sends MsgTimeoutNow to %x immediately as %x already has up-to-date log", r.id, leadTransferee, leadTransferee)
		} else {
			r.sendAppend(leadTransferee)
		}
	}
}

// stepCandidate is shared by StateCandidate and StatePreCandidate; the difference is
// whether they respond to MsgVoteResp or MsgPreVoteResp.
func stepCandidate(r *raft, m pb.Message) {
	// Only handle vote responses corresponding to our candidacy (while in
	// StateCandidate, we may get stale MsgPreVoteResp messages in this term from
	// our pre-candidate state).
	var myVoteRespType pb.MessageType
	if r.state == StatePreCandidate {
		myVoteRespType = pb.MsgPreVoteResp
	} else {
		myVoteRespType = pb.MsgVoteResp
	}
	switch m.Type {
	case pb.MsgProp:
		r.logger.Infof("%x no leader at term %d; dropping proposal", r.id, r.Term)
		return
	case pb.MsgApp:
		r.becomeFollower(r.Term, m.From)
		r.handleAppendEntries(m)
	case pb.MsgHeartbeat:
		r.becomeFollower(r.Term, m.From)
		r.handleHeartbeat(m)
	case pb.MsgSnap:
		r.becomeFollower(m.Term, m.From)
		r.handleSnapshot(m)
	case myVoteRespType:
		gr := r.poll(m.From, m.Type, !m.Reject)
		r.logger.Infof("%x [quorum:%d] has received %d %s votes and %d vote rejections", r.id, r.quorum(), gr, m.Type, len(r.votes)-gr)
		switch r.quorum() {
		case gr:
			if r.state == StatePreCandidate {
				r.campaign(campaignElection)
			} else {
				r.becomeLeader()
				r.bcastAppend()
			}
		case len(r.votes) - gr:
			r.becomeFollower(r.Term, None)
		}
	case pb.MsgTimeoutNow:
		r.logger.Debugf("%x [term %d state %v] ignored MsgTimeoutNow from %x", r.id, r.Term, r.state, m.From)
	}
}

func stepFollower(r *raft, m pb.Message) {
	switch m.Type {
	case pb.MsgProp:
		if r.lead == None {
			r.logger.Infof("%x no leader at term %d; dropping proposal", r.id, r.Term)
			return
		}
		m.To = r.lead
		r.send(m)
	case pb.MsgApp:
		r.electionElapsed = 0
		r.lead = m.From
		r.handleAppendEntries(m)
	case pb.MsgHeartbeat:
		r.electionElapsed = 0
		r.lead = m.From
		r.handleHeartbeat(m)
	case pb.MsgSnap:
		r.electionElapsed = 0
		r.lead = m.From
		r.handleSnapshot(m)
	case pb.MsgTransferLeader:
		if r.lead == None {
			r.logger.Infof("%x no leader at term %d; dropping leader transfer msg", r.id, r.Term)
			return
		}
		m.To = r.lead
		r.send(m)
	case pb.MsgTimeoutNow:
		if r.promotable() {
			r.logger.Infof("%x [term %d] received MsgTimeoutNow from %x and starts an election to get leadership.", r.id, r.Term, m.From)
			// Leadership transfers never use pre-vote even if r.preVote is true; we
			// know we are not recovering from a partition so there is no need for the
			// extra round trip.
			r.campaign(campaignTransfer)
		} else {
			r.logger.Infof("%x received MsgTimeoutNow from %x but is not promotable", r.id, m.From)
		}
	case pb.MsgReadIndex:
		if r.lead == None {
			r.logger.Infof("%x no leader at term %d; dropping index reading msg", r.id, r.Term)
			return
		}
		m.To = r.lead
		r.send(m)
	case pb.MsgReadIndexResp:
		if len(m.Entries) != 1 {
			r.logger.Errorf("%x invalid format of MsgReadIndexResp from %x, entries count: %d", r.id, m.From, len(m.Entries))
			return
		}
		r.readStates = append(r.readStates, ReadState{Index: m.Index, RequestCtx: m.Entries[0].Data})
	}
}

func (r *raft) handleAppendEntries(m pb.Message) {
	if m.Index < r.raftLog.committed {
		r.send(pb.Message{To: m.From, Type: pb.MsgAppResp, Index: r.raftLog.committed})
		return
	}

	if mlastIndex, ok := r.raftLog.maybeAppend(m.Index, m.LogTerm, m.Commit, m.Entries...); ok {
		r.send(pb.Message{To: m.From, Type: pb.MsgAppResp, Index: mlastIndex})
	} else {
		r.logger.Debugf("%x [logterm: %d, index: %d] rejected msgApp [logterm: %d, index: %d] from %x",
			r.id, r.raftLog.zeroTermOnErrCompacted(r.raftLog.term(m.Index)), m.Index, m.LogTerm, m.Index, m.From)
		r.send(pb.Message{To: m.From, Type: pb.MsgAppResp, Index: m.Index, Reject: true, RejectHint: r.raftLog.lastIndex()})
	}
}

func (r *raft) handleHeartbeat(m pb.Message) {
	r.raftLog.commitTo(m.Commit)
	r.send(pb.Message{To: m.From, Type: pb.MsgHeartbeatResp, Context: m.Context})
}

func (r *raft) handleSnapshot(m pb.Message) {
	sindex, sterm := m.Snapshot.Metadata.Index, m.Snapshot.Metadata.Term
	if r.restore(m.Snapshot) {
		r.logger.Infof("%x [commit: %d] restored snapshot [index: %d, term: %d]",
			r.id, r.raftLog.committed, sindex, sterm)
		r.send(pb.Message{To: m.From, Type: pb.MsgAppResp, Index: r.raftLog.lastIndex()})
	} else {
		r.logger.Infof("%x [commit: %d] ignored snapshot [index: %d, term: %d]",
			r.id, r.raftLog.committed, sindex, sterm)
		r.send(pb.Message{To: m.From, Type: pb.MsgAppResp, Index: r.raftLog.committed})
	}
}

// restore recovers the state machine from a snapshot. It restores the log and the
// configuration of state machine.
func (r *raft) restore(s pb.Snapshot) bool {
	if s.Metadata.Index <= r.raftLog.committed {
		return false
	}
	if r.raftLog.matchTerm(s.Metadata.Index, s.Metadata.Term) {
		r.logger.Infof("%x [commit: %d, lastindex: %d, lastterm: %d] fast-forwarded commit to snapshot [index: %d, term: %d]",
			r.id, r.raftLog.committed, r.raftLog.lastIndex(), r.raftLog.lastTerm(), s.Metadata.Index, s.Metadata.Term)
		r.raftLog.commitTo(s.Metadata.Index)
		return false
	}

	r.logger.Infof("%x [commit: %d, lastindex: %d, lastterm: %d] starts to restore snapshot [index: %d, term: %d]",
		r.id, r.raftLog.committed, r.raftLog.lastIndex(), r.raftLog.lastTerm(), s.Metadata.Index, s.Metadata.Term)

	r.raftLog.restore(s)
	r.prs = make(map[uint64]*Progress)
	for _, n := range s.Metadata.ConfState.Nodes {
		match, next := uint64(0), r.raftLog.lastIndex()+1
		if n == r.id {
			match = next - 1
		}
		r.setProgress(n, match, next)
		r.logger.Infof("%x restored progress of %x [%s]", r.id, n, r.prs[n])
	}
	return true
}

// promotable indicates whether state machine can be promoted to leader,
// which is true when its own id is in progress list.
func (r *raft) promotable() bool {
	_, ok := r.prs[r.id]
	return ok
}

func (r *raft) addNode(id uint64) {
	r.pendingConf = false
	if _, ok := r.prs[id]; ok {
		// Ignore any redundant addNode calls (which can happen because the
		// initial bootstrapping entries are applied twice).
		return
	}

	r.setProgress(id, 0, r.raftLog.lastIndex()+1)
	// When a node is first added, we should mark it as recently active.
	// Otherwise, CheckQuorum may cause us to step down if it is invoked
	// before the added node has a chance to communicate with us.
	r.prs[id].RecentActive = true
}

func (r *raft) removeNode(id uint64) {
	r.delProgress(id)
	r.pendingConf = false

	// do not try to commit or abort transferring if there is no nodes in the cluster.
	if len(r.prs) == 0 {
		return
	}

	// The quorum size is now smaller, so see if any pending entries can
	// be committed.
	if r.maybeCommit() {
		r.bcastAppend()
	}
	// If the removed node is the leadTransferee, then abort the leadership transferring.
	if r.state == StateLeader && r.leadTransferee == id {
		r.abortLeaderTransfer()
	}
}

func (r *raft) resetPendingConf() { r.pendingConf = false }

func (r *raft) setProgress(id, match, next uint64) {
	r.prs[id] = &Progress{Next: next, Match: match, ins: newInflights(r.maxInflight)}
}

func (r *raft) delProgress(id uint64) {
	delete(r.prs, id)
}

func (r *raft) loadState(state pb.HardState) {
	if state.Commit < r.raftLog.committed || state.Commit > r.raftLog.lastIndex() {
		r.logger.Panicf("%x state.commit %d is out of range [%d, %d]", r.id, state.Commit, r.raftLog.committed, r.raftLog.lastIndex())
	}
	r.raftLog.committed = state.Commit
	r.Term = state.Term
	r.Vote = state.Vote
}

// pastElectionTimeout returns true iff r.electionElapsed is greater
// than or equal to the randomized election timeout in
// [electiontimeout, 2 * electiontimeout - 1].
func (r *raft) pastElectionTimeout() bool {
	return r.electionElapsed >= r.randomizedElectionTimeout
}

func (r *raft) resetRandomizedElectionTimeout() {
	r.randomizedElectionTimeout = r.electionTimeout + globalRand.Intn(r.electionTimeout)
}

// checkQuorumActive returns true if the quorum is active from
// the view of the local raft state machine. Otherwise, it returns
// false.
// checkQuorumActive also resets all RecentActive to false.
func (r *raft) checkQuorumActive() bool {
	var act int

	for id := range r.prs {
		if id == r.id { // self is always active
			act++
			continue
		}

		if r.prs[id].RecentActive {
			act++
		}

		r.prs[id].RecentActive = false
	}

	return act >= r.quorum()
}

func (r *raft) sendTimeoutNow(to uint64) {
	r.send(pb.Message{To: to, Type: pb.MsgTimeoutNow})
}

func (r *raft) abortLeaderTransfer() {
	r.leadTransferee = None
}

func numOfPendingConf(ents []pb.Entry) int {
	n := 0
	for i := range ents {
		if ents[i].Type == pb.EntryConfChange {
			n++
		}
	}
	return n
}
