package main

import (
	"fmt"
	"myTCP/myTCP"
)

func main() {
	port := "10001"
	serverAddr, err := myTCP.ResolveName(":"+port)
	checkError(err)

	socket, err := myTCP.Listen(serverAddr)
	checkError(err)
	defer socket.Close()

	buf := make([]byte, 524)

	for {
		debug("Esperando clientes na porta " + port + "...")
		for {
			conn, err := socket.Accept()
			if err != nil {
				debug("Algum erro ocorreu ao aceitar conexão de cliente")
				continue
			}

			debug("Um cliente se conectou")
			go ReceiveConn(conn)
		}
	}
}
