package main

import (
	"strconv"
	"os"
	"io"
	"myTCP/myTCP"
)

//// for TCP version
//var nextID uint16 = 1
//
//func generateID() uint16 {
//	nextID++
//	return nextID - 1
//}

type Client struct {
	ID   uint16
	conn myTCP.ConnClient
	//conn net.Conn
	//in *bufio.Reader
	//out  *bufio.Writer
}

// Get data from network and save to the client file
func (c *Client) Read() {
	//for {
	//	data, err := receiveData(c.in, 4096)
	//	if err != nil { // Client has disconnected
	//		debug("Client " + strconv.Itoa(int(c.ID)) + " disconnected")
	//		break
	//	}
	//
	//	// TODO save data to file
	//}

	file, err := os.Create(dir + "/" + strconv.Itoa(int(c.ID)) + ".file")
	checkError(err)
	defer file.Close()

	n, err := io.Copy(file, c.conn)
	checkError(err)

	debug("Received " + strconv.Itoa(int(n)) + " bytes from client " + strconv.Itoa(int(c.ID)))

	//c.in = nil // Closing network stream
}

//// Send data to client
//func (c *Client) Write() {
//	sendData(c.out, )
//}

// Create a new Client.
func newClient(connClient myTCP.ConnClient) *Client {
	//func newClient(connClient net.Conn) *Client {
	client := &Client{
		ID: connClient.ID,
		//ID:   generateID(),
		conn: connClient,
		//in:   bufio.NewReader(connClient),
		//out:  bufio.NewWriter(connClient),
	}

	go client.Read()
	//go client.Write()

	return client
}
