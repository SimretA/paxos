package main

//var wg sync.WaitGroup

func main0() {
	nodeServers := CreateServer(5)
	m1 := value{from: 3, to: 1, typ: 1, seq: 1, preSeq: 0, val: "m1"}
	nodeServers.sendMessage(m1)
	//func(){
	//	fmt.Println("Looping")
	//	nodeServers.receiveMessage(0)
	//	nodeServers.receiveMessage(1)
	//	//defer wg.Done()
	//}()
	for i := 0; i < 5; i++ {
		recv(nodeServers, i)
	}
	//recv(nodeServers, 0)
	//recv(nodeServers, 1)
	//wg.Wait()

}
func main() {
	nodeServers := CreateServer(5)
	var acceptors []acceptor
	aId := 1
	for aId < 5 {
		accptr := CreateAcceptor(aId, nodeServers)
		acceptors = append(acceptors, accptr)
		aId++
	}
	prop := CreateProposer(0, "value", nodeServers, 1, 2, 3, 4)
	prop.run()

	for index, _ := range acceptors {
		go acceptors[index].run()
	}
}
func recv(nodes *servers, node int) {
	nodes.receiveMessage(node)
}
