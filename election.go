package Raft

type RequestVote struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

type RespondVote struct {
	trem int
	Vote bool
}

func (rf *Raft) Vote(Requset RequestVote, Respond *RespondVote) error {

	Respond.trem = rf.currentTerm
	Respond.Vote = false

	return nil
}
