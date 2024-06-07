package Raft

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Client struct {
	addr string
}

func (clientEnd *Client) Call(serviceMethod string, args interface{}, reply interface{}) bool {
	var client *rpc.Client
	var err error

	if client, err = rpc.DialHTTP("tcp", clientEnd.addr); err != nil {
		log.Fatal(err)
		return false
	}
	defer client.Close()

	if err = client.Call(serviceMethod, args, reply); err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

func (rf *Raft) initRpcServer() {
	server := rpc.NewServer()
	server.Register(rf)

	var err error
	var listener net.Listener
	if listener, err = net.Listen("tcp", rf.peers[rf.me].addr); err != nil {
		log.Fatal(err)
	}
	if err = http.Serve(listener, server); err != nil {
		log.Fatal(err)
	}
	return
}
