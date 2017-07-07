package myTCP

import (
	"net"
	"time"
	"strconv"
	"io"
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

// Create a new struct Conn.
func newConn(conn *net.UDPConn, remoteAddr *Addr, id uint16) *Conn {
	if id == 0 {
		if conn == nil && remoteAddr == nil {
			return &Conn{}
		}
		id = generateID()
	}

	return &Conn{
		ID:         id,
		newPacket:  make(chan *Packet),
		remoteAddr: remoteAddr,
		udpConn:    conn,
		timeout:    make(chan bool, 1),
	}
}

// Close a connection, checking for errors.
func (c *Conn) Close() error {
	debug("Closing connection")
	if c.udpConn != nil {
		return c.udpConn.Close()
	}

	//	 TODO Conn: close connection, remove from listener if it is a server conn

	return nil
}

// Take a packet from a connection, copying the payload into p
func (c Conn) Read(p []byte) (n int, err error) {
	debug("READING " + strconv.Itoa(len(p)) + " bytes")
	//	TODO Read: coordinate receiving of packets, create packet/payload BUFFER, return only until the requested size

	for packet := range c.newPacket {
		payload := packet.payload

		// TODO detect fin and return EOF
		//if len(payload) == 0{
		//	break
		//}

		copy(p[n:n+len(payload)], payload)

		debug("A: " + string(payload))
		debug("B: " + string(p[n:n+len(payload)]))

		n += len(payload)

		// TODO send ACK

		// TODO keep payload rest, if it doesnt fit in p buffer

		if len(payload) < 512 {
			debug("LITTLE MESSAGE")
			err = io.EOF
			break
		}

		if n >= len(p) {
			break
		}
		//if len(p) <= 512 {
		//	debug("LITTLE MESSAGE")
		//
		//} else {
		//	debug("BIG MESSAGE")
		//
		//
		//
		//	panic("IIIIIIIIIII")
		//}
	}

	debug("OK")

	// TODO verifies when Copy func stops reading
	//c.Close()

	// TODO catch some error
	return n, err
}

// Write a packet to a connection
func (c Conn) Write(p []byte) (int, error) {
	debug("WRITING " + strconv.Itoa(len(p)) + " bytes : " + string(p))
	// TODO Write: coordinate sending of packets, break bytes on packets

	var qtd_wrote int

	if len(p) <= 512 {
		debug("LITTLE MESSAGE")

		payload := make([]byte, len(p))
		copy(payload[:], p)

		header := newHeader(24524, 1231241, c.ID,
			false, false, false)
		addr, err := ResolveName("127.0.0.1:0")
		if err != nil {
			return 0, nil
		}
		packet := newPacket(header, payload, addr).compact()

		n, err := c.udpConn.Write(packet)
		qtd_wrote += n - 12
		if err != nil {
			return 0, nil
		}
	} else {
		debug("BIG MESSAGE")

		var data [512]byte
		copy(data[:], p[qtd_wrote:qtd_wrote+512])

		panic("AAAAAAAAAAAA")
	}

	return qtd_wrote, nil
}

func generateID() uint16 {
	nextID++
	return nextID - 1
}

func (c *Conn) connTimeout(timeSeconds time.Duration) {
	time.Sleep(timeSeconds * time.Second)
	c.timeout <- true
}
