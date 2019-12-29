package main

import (
	"fmt"
	"log"
)

func CreateProposer(id int, val string, srvr *servers, acceptors ...int) *proposer {
	prop := proposer{id: id, proposeVal: val, seq: 0, srvr: srvr}
	prop.acceptors = make(map[int]value, len(acceptors))

	for _, acceptor := range acceptors {
		prop.acceptors[acceptor] = value{}
	}
	return &prop
}

type proposer struct {
	id         int
	seq        int
	proposeNum int
	proposeVal string
	acceptors  map[int]value
	srvr       *servers
}

func (p *proposer) run() {
	for !(p.getRecevPromiseCount() > p.majority()) {
		sendMsgs := p.prepare()
		for _, val := range sendMsgs {
			p.srvr.sendMessage(val)
			fmt.Println("propose sent")
		}

		m := p.srvr.receiveMessage(p.id)
		if m == nil {
			log.Println("[Proposer: no msg... ")
			continue
		}
		log.Println("[Proposer: recev", m)
		switch m.typ {
		case 2: //for promise
			log.Println(" proposer recev a promise from ", m.from)
			p.checkReceivePromise(*m)
		default:
			panic("Unsupport message.")
		}

	}
}

func (p *proposer) majority() int {
	return len(p.acceptors)/2 + 1
}
func (p *proposer) getRecevPromiseCount() int {
	recvCount := 0
	for _, acepMsg := range p.acceptors {
		log.Println(" proposer has total ", len(p.acceptors), " acceptor ", acepMsg, " current Num:", p.getNewSeq(), " msgNum:", acepMsg.seq)
		if acepMsg.seq == p.getNewSeq() {
			log.Println("recv ++", recvCount)
			recvCount++
		}
	}
	log.Println("Current proposer recev promise count=", recvCount)
	return recvCount
}

func (p *proposer) prepare() []value {
	p.seq++
	msgCount := 0
	var msgList []value
	for acepId, _ := range p.acceptors {
		val := value{from: p.id, to: acepId, typ: 1, seq: p.getNewSeq(), val: p.proposeVal}
		msgList = append(msgList, val)
		msgCount++
		if msgCount > p.majority() {
			break
		}
	}
	return msgList
}

func (p *proposer) getNewSeq() int {
	p.proposeNum = p.seq<<4 | p.id
	return p.proposeNum
}

func (p *proposer) checkReceivePromise(v value) {

}
