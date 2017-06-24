package myTCP

import "net"

var nextID uint16 = 1

// Wrapper for UDPConn
type Conn struct {
	conn *net.UDPConn
	ID   uint16
}

// Creates a new struct Conn
func newConn(conn *net.UDPConn) *Conn {
	return &Conn{
		conn: conn,
		ID:   generateID(),
	}
}

// Closes a connection, checking for errors
func (c *Conn) Close() {
	debug("Closing connection")
	checkError(c.conn.Close())
}

// Reads a packet from a connection, copying the payload into b
func (c *Conn) Read(b []byte) (int, *Addr, error) {
	debug("Reading a packet")
	n, addr, err := c.conn.ReadFromUDP(b)
	return n, newAddr(addr), err
}

// Writes a packet to a connection
func (c *Conn) Write(b []byte) (int, error) {
	var packet [524]byte

	if len(b) < 512 {
		debug("LITTLE MESSAGE")

		var data [512]byte
		copy(data[:], b)

		packet = newDataPacket(data).compact()
	} else {
		debug("BIG MESSAGE")

		var data [512]byte
		copy(data[:], b)

		packet = newDataPacket(data).compact()
	}

	return c.conn.Write(packet[:])
}

func generateID() uint16 {
	nextID++
	return nextID - 1
}
