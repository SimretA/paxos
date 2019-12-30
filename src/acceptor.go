package main

import "fmt"

func CreateAcceptor(id int, srvr *servers) acceptor {
	return acceptor{id: id, srvr: srvr}

}

type acceptor struct {
	id         int
	acceptMsg  value
	promiseMsg value
	srvr       *servers
}

func (accptr *acceptor) run() {
	for {
		val := accptr.srvr.receiveMessage(accptr.id)
		if val == nil {
			continue
		}
		switch val.typ {
		case 1: //for prepare
			PromiseMsg := accptr.receivePrepare(*val)
			accptr.srvr.sendMessage(*PromiseMsg)
			fmt.Println("Acceptor ", accptr.id, " received prepare from ", val.from)
			continue
		case 3: //for propose
			accpted := accptr.receivePropose(*val)
			if accpted {
				fmt.Println("Value ", val.val, " accepted by ", accptr.id)
				//TODO broadcast to learners
			} else {
				fmt.Println("Value not accepted by ", accptr.id)
			}
		default:
			fmt.Println()

		}
	}
}

func (accptr *acceptor) receivePrepare(v value) *value {
	if accptr.promiseMsg.seq >= v.seq {
		fmt.Println("Acceptor ", accptr.id, " has Already promised")
		return nil //ignore
	}
	fmt.Println("Acceptor ", accptr.id, " promise")
	v.to = v.from
	v.from = accptr.id
	v.typ = 2 //for promise
	accptr.acceptMsg = v
	return &v
}

func (accptr *acceptor) receivePropose(v value) bool {
	if accptr.acceptMsg.seq > v.seq {
		return false
	}
	return true
}
