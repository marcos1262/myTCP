package main

import (
	"fmt"
	"myTCP/myTCP"
)

func main() {
	serverAddr, err := myTCP.ResolveName(":10001")
	checkError(err)

	serverConn, err := myTCP.Listen(serverAddr)
	checkError(err)
	defer serverConn.Close()

	buf := make([]byte, 524)

	for {
		n, addr, err := serverConn.Read(buf)
		checkError(err)

		fmt.Println("Received ", string(buf[:n]), " (", n, " bytes) from ", addr)
	}
}
