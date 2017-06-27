package myTCP

import "net"

var nextID uint16 = 1

// Conn represents a single client connection
type Conn struct {
	ID         uint16
	newPacket  chan *Packet
	remoteAddr *Addr
	udpConn    *net.UDPConn
}

// newConn creates a new struct Conn
func newConn(conn *net.UDPConn) *Conn {
	return &Conn{
		ID:        generateID(),
		newPacket: make(chan *Packet),
	}
}

// Closes a connection, checking for errors
func (c *Conn) Close() {
	debug("Closing connection")
	//	 TODO Conn: close connection
}

// Reads a packet from a connection, copying the payload into b
func (c *Conn) Read(b []byte) (int, *Addr, error) {
	for packet := range c.newPacket {
		//	TODO Read: coordinate receiving of packets, create packet BUFFER, return only until the requested size

	}
	//n, addr, err := c.conn.ReadFromUDP(b)
	//return n, newAddr(addr), err
	return 0, nil, nil
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
