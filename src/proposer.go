package main

import (
	"fmt"
	"log"
	"time"
)

func CreateProposer(id int, val string, srvr *servers, acceptors ...int) *proposer {
	prop := proposer{id: id, proposeVal: val, seq: 1, srvr: srvr}
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
		fmt.Println("recvd promises ", p.getRecevPromiseCount())
		//for {
		fmt.Println("inloop")
		sendMsgs := p.prepare()
		for _, val := range sendMsgs {
			p.srvr.sendMessage(val)
			fmt.Println("propose sent")
		}
		fmt.Println("Sleeping")
		time.Sleep(10 * time.Millisecond)
		p.receiveMessage()
		fmt.Println("Done")

	}
	fmt.Println("Proposing")
	msgs := p.propose()
	fmt.Println(msgs)
	for _, msg := range msgs {

		p.srvr.sendMessage(msg)
		fmt.Println("Proposed value ", msg.val, " to ", msg.to)
	}
}
func (p *proposer) receiveMessage() {
	m := p.srvr.receiveMessage(0)
	if m == nil {
		log.Println("[Proposer: no msg... ")
		return
	}
	log.Println("[Proposer: recev", m)
	switch m.typ {
	case 2: //for promise
		log.Println(" proposer recev a promise from ", m.from)
		p.checkReceivePromise(*m)
	default:
		return
		//panic("Unsupport message.")
	}

}
func (p *proposer) majority() int {
	return len(p.acceptors)/2 + 1
}
func (p *proposer) getRecevPromiseCount() int {
	recvCount := 0
	for _, acepMsg := range p.acceptors {
		fmt.Println("Proposer has total ", len(p.acceptors), " acceptors  current Seq:", p.getNewSeq(), " Accepted message sequence:", p.proposeNum)
		if acepMsg.seq <= p.getNewSeq() {
			fmt.Println("Received ", recvCount)
			recvCount++
		}
	}
	fmt.Println("Current proposer recev promise count=", recvCount)
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

func (p *proposer) propose() []value {
	p.seq++

	sentCount := 0
	var msgs []value
	for acceptorId, acceptorVal := range p.acceptors {
		if acceptorVal.seq <= p.getNewSeq() {
			newMsg := value{from: 0, to: acceptorId, typ: 3, seq: p.getNewSeq(), val: p.proposeVal}
			msgs = append(msgs, newMsg)
			sentCount++

		}
	}
	return msgs
}

func (p *proposer) getNewSeq() int {
	p.proposeNum = p.seq<<4 | p.id
	return p.proposeNum
}

func (p *proposer) checkReceivePromise(v value) {
	previousPromise := p.acceptors[v.from]
	fmt.Println("New promise ", v)
	if previousPromise.seq < v.seq {
		p.acceptors[v.from] = v
		if v.seq > p.getNewSeq() {
			p.proposeNum = v.seq
			p.proposeVal = v.val

		}
	}

}
