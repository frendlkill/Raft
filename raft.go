package Raft

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type LogEntry struct {
	Command interface{}
	Term    int
}

type Raft struct {
	mu    sync.Mutex
	peers []*Client //记录ip，然后调用内部的CALL方法调用RAFT的函数
	me    int
	dead  int32

	currentTerm int // 见过的最大任期
	votedFor    int // 记录在currentTerm任期投票给谁了

	log               []LogEntry
	lastIncludedIndex int
	lastIncludedTerm  int

	commitIndex int // 已知的最大已提交索引
	lastApplied int // 当前应用到状态机的索引

	role              int       // 1:leader 2:follower 3:candidates
	leaderId          int       // leader的id
	lastActiveTime    time.Time // 上次活跃时间（刷新时机：收到leader心跳、给其他candidates投票、请求其他节点投票）
	lastBroadcastTime time.Time // 作为leader，上次的广播时间
}

func (rf *Raft) electionTimeout() time.Duration {
	return time.Duration((rand.Int31n(150) + 150)) * time.Millisecond
}

func Make(addrs []string, me int) (rf *Raft, err error) {
	rf = &Raft{}
	rf.me = me
	rf.role = 2
	rf.log = make([]LogEntry, 0)
	rf.lastIncludedIndex = 0
	rf.lastIncludedTerm = -1
	rf.votedFor = -1
	rf.lastActiveTime = time.Now()
	rf.initConn(addrs)
	go rf.initRpcServer()
	go rf.election()
	//go rf.state()
	go rf.Hert()
	return
}

type T1 struct{}

func (rf *Raft) GetState(p T1, ret *int) error {
	*ret = rf.me
	return nil
}

func (rf *Raft) initConn(addrs []string) {
	peers := make([]*Client, 0)
	for _, addr := range addrs {
		peers = append(peers, &Client{addr: addr})
	}
	rf.peers = peers
	return
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}
