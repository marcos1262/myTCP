package myTCP

import "net"

// Wrapper for UDPAddr
type Addr struct {
	addr *net.UDPAddr
}

// Wrapper for UDPConn
type Conn struct {
	conn *net.UDPConn
}

type Packet struct {
	header Header
	payload [512]byte
}

// Creates a new struct Addr
func newAddr(addr *net.UDPAddr) *Addr {
	return &Addr{addr: addr}
}

// Creates a new struct Conn
func newConn(conn *net.UDPConn) *Conn {
	return &Conn{conn: conn}
}

// Parses to string
func (a *Addr) String () string {
	return a.addr.String()
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
	return c.conn.Write(b)
}

// Parses a host name to IP/Port
func ResolveName(addr string) (*Addr, error) {
	debug("Resolving hostname")
	ServerAddr, err := net.ResolveUDPAddr("udp", addr)
	return newAddr(ServerAddr), err
}

// Listens to clients
func Listen(addr *Addr) (*Conn, error) {
	ServerConn, err := net.ListenUDP("udp", addr.addr)
	return newConn(ServerConn), err
}

// Connects to a server
func Connect(remoteAddr *Addr) (*Conn, error) {
	debug("Connecting to a server")
	debug("(REMOVE THIS) LOCAL ADDR: " + remoteAddr.addr.IP.String())
	localAddr, err := net.ResolveUDPAddr("udp", remoteAddr.addr.IP.String()+":0")
	checkError(err)

	conn, err := net.DialUDP("udp", localAddr, remoteAddr.addr)
	return newConn(conn), err
}
