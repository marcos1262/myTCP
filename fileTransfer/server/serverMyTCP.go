package main

import (
	//"myTCP/myTCP"
	"net"
	"sync"
	"os"
)

var clients = make(map[uint16]*Client)
var clientsMutex = sync.RWMutex{}

var dir = "save"

// Search for a saved client by the ID.
func searchClient(ID uint16) (*Client, bool) {
	clientsMutex.RLock() // Mutual exclusion (reading)
	client, exists := clients[ID]
	clientsMutex.RUnlock()

	return client, exists
}

// Save a client on the list.
func saveClient(client *Client) {
	clientsMutex.Lock() // Mutual exclusion (writing)
	clients[client.ID] = client
	clientsMutex.Unlock()
}

// Remove a saved client, searching by the ID.
func deleteClient(ID uint16) {
	clientsMutex.Lock() // Mutual exclusion (writing)
	delete(clients, ID)
	clientsMutex.Unlock()
}

// Receive client on the server.
//func ReceiveConn(conn myTCP.Conn) {
func ReceiveConn(conn net.Conn) {
	debug("Giving client a coffee")

	saveClient(newClient(conn))
}

func main() {
	port := "10101"
	//flag.StringVar(&port, "p", "", "Give a port with -p=#####")
	//flag.StringVar(&dir, "d", "", "Give a dir name with -d=#####")
	//flag.Parse()
	//
	//if port == "" || dir == "" {
	//	flag.PrintDefaults()
	//	os.Exit(1)
	//}

	//serverAddr, err := myTCP.ResolveName(":" + port)
	serverAddr, err := net.ResolveTCPAddr("tcp", ":"+port)

	checkError(err)

	//socket, err := myTCP.Listen(serverAddr)
	socket, err := net.ListenTCP("tcp", serverAddr)
	checkError(err)
	defer socket.Close()

	os.Mkdir(dir, 0777)

	debug("Listening clients at port " + port + "...")
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
