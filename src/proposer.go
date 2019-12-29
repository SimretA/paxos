package main

type proposer struct {
	id         int
	seq        int
	proposeNum int
	proposeVal string
	promiseCount int
	acceptors  map[int]value
	srvr         *servers
}

func (p *proposer) run() {

}

