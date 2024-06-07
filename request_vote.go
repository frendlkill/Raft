package Raft

import (
	"log"
	"math/rand"
	"time"
)

type RequestVoteArgs struct {
	Term        int
	CandidateId int
	//LastLogIndex int
	//LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

func (rf *Raft) VoteRpc(args *RequestVoteArgs, reply *RequestVoteReply) (err error) {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	reply.VoteGranted = false
	reply.Term = rf.currentTerm
	if args.Term < rf.currentTerm {
		return
	}
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.votedFor = -1
		rf.role = 2
		rf.leaderId = -1
	}
	if rf.votedFor == -1 {
		rf.votedFor = args.CandidateId
		reply.VoteGranted = true
		rf.lastActiveTime = time.Now()
	}
	log.Println(rf.me, "voteFor", rf.votedFor, "term", rf.currentTerm)
	return
}
func (rf *Raft) state() (err error) {
	for !rf.killed() {
		time.Sleep(1000 * time.Millisecond)
		func() {
			rf.mu.Lock()

			defer rf.mu.Unlock()
			log.Println("", rf.me, rf.currentTerm, rf.role, rf.lastActiveTime)
		}()

	}
	return nil
}
func (rf *Raft) election() (err error) {
	for !rf.killed() {

		time.Sleep(100 * time.Millisecond)

		func() {
			rf.mu.Lock()
			defer rf.mu.Unlock()
			now := time.Now()
			D := time.Duration(150+rand.Intn(200)) * time.Millisecond
			elapses := now.Sub(rf.lastActiveTime)
			if rf.role == 2 && elapses > D {
				rf.role = 3
			}
			if rf.role == 3 && elapses > D {
				log.Println(rf.me, "第", rf.currentTerm, "选举开始")
				rf.lastActiveTime = now
				rf.currentTerm += 1
				rf.votedFor = rf.me
				args := RequestVoteArgs{
					Term:        rf.currentTerm,
					CandidateId: rf.me,
				}
				type VoteResult struct {
					peerId int
					resp   *RequestVoteReply
				}
				voteCount := 1 // 收到投票个数（先给自己投1票）
				finishCount := 1

				voteChan := make(chan *VoteResult, len(rf.peers))
				for peerId := 0; peerId < len(rf.peers); peerId++ {
					go func(id int) {
						if peerId == rf.me {
							return
						}
						resp := RequestVoteReply{}
						ok := rf.sendRequestVote(peerId, &args, &resp)
						if ok {
							voteChan <- &VoteResult{peerId, &resp}
						} else {
							voteChan <- &VoteResult{peerId, nil}
						}

					}(peerId)
				}
				maxTerm := 0
				for {
					select {
					case Vote := <-voteChan:
						finishCount++
						if Vote.resp != nil {
							if Vote.resp.Term > maxTerm {
								maxTerm = Vote.resp.Term
							}
							if Vote.resp.VoteGranted {
								voteCount += 1
							}
						}
						if finishCount == len(rf.peers) || voteCount > len(rf.peers)/2 {
							goto VOTE_END
						}
					}
					break
				}
			VOTE_END:
				log.Println(rf.me, "选举结束得票", voteCount)
				if rf.role != 3 {
					return
				}

				if maxTerm > rf.currentTerm {
					rf.role = 2
					rf.leaderId = -1
					rf.currentTerm = maxTerm
					rf.votedFor = -1
					return
				}
				if voteCount > len(rf.peers)/2 {
					log.Println(rf.me, "成为了第", rf.currentTerm, "任村长")
					rf.role = 1
					rf.leaderId = rf.me
					rf.lastBroadcastTime = time.Unix(0, 0) // 令appendEntries广播立即执行
					return
				}

			}
		}()

	}
	return nil
}
func (rf *Raft) sendRequestVote(id int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[id].Call("Raft.VoteRpc", args, reply)
	return ok
}
