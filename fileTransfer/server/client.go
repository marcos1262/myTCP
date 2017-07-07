package main

import (
	"bufio"
	"net"
	//"myTCP/myTCP"
	"strconv"
	"os"
	"io"
)

// FIXME remove me
var nextID uint16 = 1

func generateID() uint16 {
	nextID++
	return nextID - 1
}

type Client struct {
	ID uint16
	//conn myTCP.Conn
	conn net.Conn
	in   *bufio.Reader
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
	io.Copy(file, c.in)

	c.in = nil // Closing network stream
}

//// Send data to client
//func (c *Client) Write() {
//	sendData(c.out, )
//}

// Create a new Client.
//func newClient(conn myTCP.Conn) *Client {
func newClient(conn net.Conn) *Client {
	client := &Client{
		//ID: conn.ID,
		ID:   generateID(),
		conn: conn,
		in:   bufio.NewReader(conn),
		//out:  bufio.NewWriter(conn),
	}

	go client.Read()
	//go client.Write()

	return client
}
