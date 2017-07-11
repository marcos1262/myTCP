package myTCP

import (
	"net"
	"time"
)

var nextID uint16 = 1

// FIXME remove
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
	udpConn     *net.UDPConn // connection with a server
	localAddr   *Addr        // local address
	window      window       // window for sending data
	resultWrite chan error   // to wait for the writing to end
	outPacket   chan []byte  // to send data for the sender routine
	ackPacket   chan *Packet // to receive ack from the client listener
}

type ConnClient struct {
	*conn
	listener        *Listener        // listener that keep the conn
	newPacket       chan *Packet     // to receive data from network
	getData         chan []byte      // to get data from client
	resultRead      chan resultRead       // to wait for the reading to end
	timeoutInactive <-chan time.Time // for when the client stay inactive for a long time
}

// Create a new struct ConnServer.
func newConnServer(udpConn *net.UDPConn, localAddr *Addr, remoteAddr *Addr, id uint16) *ConnServer {
	return &ConnServer{
		conn: &conn{
			ID:         id,
			remoteAddr: remoteAddr,
		},
		udpConn:     udpConn,
		localAddr:   localAddr,
		resultWrite: make(chan error),
		outPacket:   make(chan []byte),
		ackPacket:   make(chan *Packet),
	}
}

// Create a new struct ConnClient.
func newConnClient(remoteAddr *Addr, listener *Listener) *ConnClient {
	conn := &ConnClient{
		conn: &conn{
			ID:         generateID(),
			remoteAddr: remoteAddr,
		},
		listener:  listener,
		newPacket: make(chan *Packet),
		resultRead: make(chan resultRead),
		getData:   make(chan []byte),
	}

	conn.receivePacket()

	return conn
}

func (c conn) RemoteAddr() *Addr {
	return c.remoteAddr
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

// FIXME remove
//func (c *ConnClient) connTimeout(timeSeconds time.Duration) {
//	time.Sleep(timeSeconds * time.Second)
//	c.timeoutInactive <- true
//}

// TODO implement ConnServer Read
func (c ConnServer) Read(p []byte) (n int, err error) {
	return 0, nil
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
