package myTCP

import (
	"net"
	"time"
)

var nextID uint16 = 1

// Conn represents a single client connection
type Conn struct {
	ID         uint16
	newPacket  chan *Packet
	remoteAddr *Addr
	udpConn    *net.UDPConn // Only for MyTCP Client
	timeout    chan bool    // Only for MyTCP Server
}

// newConn creates a new struct Conn
func newConn(conn *net.UDPConn, remoteAddr *Addr) *Conn {
	newConn := &Conn{
		ID:         generateID(),
		newPacket:  make(chan *Packet),
		remoteAddr: remoteAddr,
		udpConn:    conn,
		timeout:    make(chan bool, 1),
	}

	return newConn
}

// Close closes a connection, checking for errors
func (c *Conn) Close() error {
	debug("Closing connection")
	if c.udpConn != nil {
		return c.udpConn.Close()
	}

	//	 TODO Conn: close connection, remove from listener if it is a server conn

	return nil
}

// Read takes a packet from a connection, copying the payload into p
func (c *Conn) Read(p []byte) (n int, err error) {
	//for packet := range c.newPacket {
	//	//	//	TODO Read: coordinate receiving of packets, create packet/payload BUFFER, return only until the requested size
	//	//
	//}
	//n, addr, err := c.conn.ReadFromUDP(b)
	//return n, newAddr(addr), err
	return 0, nil
}

// Writes a packet to a connection
func (c *Conn) Write(b []byte) (int, error) {
	var packet [524]byte

	// TODO Write: coordinate sending of packets, break bytes on packets

	if len(b) < 512 {
		debug("LITTLE MESSAGE")

		var data [512]byte
		copy(data[:], b)

		//packet = newDataPacket(data).compact()
	} else {
		debug("BIG MESSAGE")

		var data [512]byte
		copy(data[:], b)

		//packet = newDataPacket(data).compact()
	}

	return c.udpConn.Write(packet[:])
}

func generateID() uint16 {
	nextID++
	return nextID - 1
}

func (c *Conn) connTimeout(timeSeconds time.Duration) {
	time.Sleep(timeSeconds * time.Second)
	c.timeout <- true
}
