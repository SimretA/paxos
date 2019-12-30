package main

//var wg sync.WaitGroup

func main() {
	nodeServers := CreateServer(5)
	var acceptors []acceptor
	aId := 1
	for aId < 5 {
		accptr := CreateAcceptor(aId, nodeServers)
		acceptors = append(acceptors, accptr)
		aId++
	}
	for index, _ := range acceptors {
		go acceptors[index].run()
	}
	prop := CreateProposer(0, "value", nodeServers, 1, 2, 3, 4)
	prop.run()

	for index, _ := range acceptors {
		acceptors[index].run()
	}

}
func recv(nodes *servers, node int) {
	nodes.receiveMessage(node)
}
