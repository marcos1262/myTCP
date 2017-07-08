package myTCP

import (
	"net"
	"time"
	"strconv"
	"io"
)

var nextID uint16 = 1

//type Conn interface {
//	Read(p []byte) (n int, err error)
//	Write(p []byte) (n int, err error)
//	Close() error
//}

// conn represents a single client connection
type conn struct {
	ID         uint16
	remoteAddr *Addr
}

type ConnServer struct {
	*conn
	udpConn *net.UDPConn
}

type ConnClient struct {
	*conn
	newPacket chan *Packet
	timeout   chan bool
}

// Create a new struct ConnServer.
func newConnServer(udpConn *net.UDPConn, remoteAddr *Addr, id uint16) *ConnServer {
	return &ConnServer{
		conn: &conn{
			ID:         id,
			remoteAddr: remoteAddr,
		},
		udpConn: udpConn,
	}
}

// Create a new struct ConnClient.
func newConnClient(remoteAddr *Addr) *ConnClient {
	return &ConnClient{
		conn: &conn{
			ID:         generateID(),
			remoteAddr: remoteAddr,
		},
		newPacket: make(chan *Packet),
		timeout:   make(chan bool, 1),
	}
}

func (c conn) RemoteAddr() *Addr {
	return c.remoteAddr
}

// Take a packet from a connection, copying the payload into p
func (c ConnClient) Read(p []byte) (n int, err error) {
	debug("READING " + strconv.Itoa(len(p)) + " bytes")
	//	TODO Read: coordinate receiving of packets,
	// 				create packet/payload BUFFER,
	// 				return only until the requested size

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

// TODO implement ConnClient Write
func (c ConnClient) Write(p []byte) (n int, err error) {
	return 0, nil
}

// Close a connection, checking for errors.
func (c ConnClient) Close() error {
	debug("Closing connection")

	// TODO conn: close connection, remove from listener if it is a server conn

	close(c.newPacket)
	return nil
}

func (c *ConnClient) connTimeout(timeSeconds time.Duration) {
	time.Sleep(timeSeconds * time.Second)
	c.timeout <- true
}

// TODO implement ConnServer Read
func (c ConnServer) Read(p []byte) (n int, err error) {
	return 0, nil
}

// Write a packet to a connection
func (c ConnServer) Write(p []byte) (n int, err error) {
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

// Close a connection, checking for errors.
func (c *ConnServer) Close() error {
	debug("Closing connection")

	// TODO conn: close connection, remove from listener if it is a server conn

	return c.udpConn.Close()
}

func generateID() uint16 {
	nextID++
	return nextID - 1
}
