package main

import (
	"fmt"
	"time"
)


func CreateServer(num int) *servers{
	srvrs := servers{recvQueue: make(map[int] chan value, 0)}

	for i := 0; i< num; i++{
		fmt.Println("creating server ", i)
		srvrs.recvQueue[i]= make(chan value, 1024)
	}
	return &srvrs
}

type servers struct {
	recvQueue map[int]chan value
}

func (s *servers) sendMessage(value2 value) {
	fmt.Println("Sending ", value2.typ, "to ", value2.to)
	s.recvQueue[value2.to] <- value2
}
func (s *servers) receiveMessage(id int) *value {

	select {
	case msg:= <-s.recvQueue[id]:
		//msg := <-s.recvQueue[id]
		fmt.Println("Received message ",msg.typ," Value ",msg.val)
		return &msg
	case <-time.After(1000):
		fmt.Println("Time Out")
		return nil
	}
}


/**
typ:
	1 for Prepare
	2 for Promise
	3 for Propose
	4 for Accept
**/


type value struct {
	from   int
	to     int
	typ    int
	seq    int
	preSeq int
	val    string
}
