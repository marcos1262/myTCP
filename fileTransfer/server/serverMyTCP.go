package main

import (
	"myTCP/myTCP"
)

func ReceiveConn(conn *myTCP.Conn) {
	debug("Giving client a coffee")

	// TODO receive conn
}

func main() {
	port := "10001"
	serverAddr, err := myTCP.ResolveName(":" + port)
	checkError(err)

	socket, err := myTCP.Listen(serverAddr)
	checkError(err)
	defer socket.Close()

	//buf := make([]byte, 524)

	for {
		conn, err := socket.Accept()
		if err != nil {
			debug("Something really bad happened while accepting client...")
			continue
		}

		debug("A client has connected")
		go ReceiveConn(conn)
	}
}
