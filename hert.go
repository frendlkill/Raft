package Raft

import (
	"time"
)

type HertBag struct {
	LeaderId int
	Term     int
}
type R1 struct {
	ErrId    int
	LeaderId int
	Term     int
}

func (rf *Raft) Hert() {
	for !rf.killed() {
		time.Sleep(50 * time.Millisecond)
		func() {
			rf.mu.Lock()
			defer rf.mu.Unlock()
			if rf.role != 1 {
				return
			}
			h := HertBag{
				Term:     rf.currentTerm,
				LeaderId: rf.me,
			}
			ret := R1{
				ErrId:    -1,
				LeaderId: -1,
				Term:     -1,
			}
			rf.lastActiveTime = time.Now()
			for i := 0; i < len(rf.peers); i++ {
				go func() {
					if i == rf.me {
						return
					}
					err := rf.peers[i].Call("Raft.SendHert", &h, &ret)
					if ret.ErrId == 1 {
						rf.leaderId = ret.LeaderId
						rf.role = 2
						rf.currentTerm = ret.Term
					}
					if err {
						//log.Println(rf.me, "心跳传达", i)
					}
				}()
				if rf.role != 1 {
					break
				}
			}
		}()
	}
}

func (rf *Raft) SendHert(h *HertBag, ret *R1) error {
	//rf.mu.Lock()
	//defer rf.mu.Unlock()
	//log.Println(rf.me, "收到来自", h.LeaderId, "的心跳")
	if rf.currentTerm < h.Term {
		rf.currentTerm = h.Term
		rf.role = 2
		rf.leaderId = h.LeaderId
		rf.lastActiveTime = time.Now()

	} else if rf.currentTerm > h.Term {
		ret.ErrId = 1
		ret.Term = rf.currentTerm
		ret.LeaderId = rf.leaderId
	} else {
		rf.lastActiveTime = time.Now()
	}
	return nil
}

//func (rf *Raft)Res bool {
//
//}
