package myTCP_test

import (
	"testing"
	"net"
	"myTCP/myTCP"
)

func TestAddr_String(t *testing.T) {
	serverAddr, _ := myTCP.ResolveName("127.0.0.1:12345")
	if serverAddr.String() != "127.0.0.1:12345" {
		t.Error("String returning different data")
	}
}

func TestResolveName(t *testing.T) {
	serverAddr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:12345")
	serverAddr2, _ := myTCP.ResolveName("127.0.0.1:12345")

	if serverAddr.String() != serverAddr2.String() {
		t.Error("ResolveName resulting different from ResolveUDPAddr")
	}
}

func TestListen(t *testing.T) {
	serverAddr, _ := myTCP.ResolveName("127.0.0.1:12345")
	_, err := myTCP.Listen(serverAddr)

	if err != nil {
		t.Error("Cannot initialize server using Listen (" + err.Error() + ")")
	}
}

func TestConnect(t *testing.T) {
	serverAddr, _ := myTCP.ResolveName("127.0.0.1:12345")
	myTCP.Listen(serverAddr)

	_, err := myTCP.Connect(serverAddr)
	if err != nil {
		t.Error("Cannot connect to a server (" + err.Error() + ")")
	}
}

func TestConn_Read(t *testing.T) {
	//serverAddr, _ := myTCP.ResolveName("127.0.0.1:12345")
	//conn, _ := myTCP.Listen(serverAddr)
	//
	//buf := make([]byte, 2)
	//go func() {
	//	_, _, err := conn.Read(buf)
	//	if err != nil {
	//		t.Error("Cannot read from another host (" + err.Error() + ")")
	//	}
	//	if string(buf) != "OK" {
	//		t.Error("Read returning different data")
	//	}
	//	return
	//}()

	// TODO discover how to create another process to connect with conn
	//conn2, _ := Connect(serverAddr)
	//conn2.Write([]byte("OK"))
}

func TestConn_Write(t *testing.T) {
	//serverAddr, _ := myTCP.ResolveName("127.0.0.1:12345")
	//conn, _ := myTCP.Listen(serverAddr)
	//
	//buf := make([]byte, 2)
	//go func() {
	//	conn.Read(buf)
	//	if string(buf) != "OK" {
	//		t.Error("Write not sending all the data")
	//	}
	//	return
	//}()

	// TODO same fix as above
	//conn2, _ := Connect(serverAddr)
	//_, err := conn2.Write([]byte("OK"))
	//if err != nil {
	//	t.Error("Cannot write to another host (" + err.Error() + ")")
	//}
}

func TestConn_Close(t *testing.T) {
	serverAddr, _ := myTCP.ResolveName("127.0.0.1:12345")
	myTCP.Listen(serverAddr)
	conn, _ := myTCP.Connect(serverAddr)

	conn.Close()
	_, err := conn.Write([]byte(""))
	if err == nil {
		t.Error("Connection not closing properly")
	}
}
