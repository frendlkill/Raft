package Raft

import (
	"log"
	"testing"
	"time"
)

func TestRaft(t *testing.T) {

	addrs := []string{"127.0.0.1:1500", "127.0.0.1:1501", "127.0.0.1:1502"}

	for i := 0; i < len(addrs); i++ {
		go func(idx int) {
			rf, _ := Make(addrs, idx)
			ret := 1
			s := T1{}
			err := rf.peers[1].Call("Raft.GetState", s, &ret)

			log.Println(err, ret)
		}(i)
	}
	time.Sleep(1 * time.Hour)
}
