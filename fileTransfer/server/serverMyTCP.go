package main

import (
	"net"
	"fmt"
)

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	checkError(err)

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	checkError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
