package myTCP

import (
	"net"
	"sync"
)

// Addr is a wrapper for UDPAddr.
type Addr struct {
	udpAddr *net.UDPAddr
}

// Listener is a MyTCP network listener.
type Listener struct {
	addr       *Addr
	udpConn    *net.UDPConn
	newConn    chan *Packet
	conns      map[string]*Conn
	connsMutex *sync.RWMutex
}

// newAddr creates a new struct Addr
func newListener(addr *Addr, udpConn *net.UDPConn) *Listener {
	return &Listener{
		addr: addr,
		udpConn:udpConn,
		newConn: make(chan *Packet),
		conns: make(map[string]*Conn),
	}
}

// Accept implements the Accept method in the net.Listener interface;
// Accept waits for and returns the next connection to the listener
func (l *Listener) Accept() (Conn, error) {
	// TODO Accept: receive request from channel, create connection
	return nil, nil
}

// Close stops listening on the TCP address.
// Already Accepted connections are not closed.
func (l *Listener) Close() {
	close(l.newConn)
}

// Packet represents structure
type Packet struct {
	header  *Header
	payload [512]byte
}

// compact compresses a packet into a byte array
func (p *Packet) compact() [524]byte {
	var packet [524]byte
	header := p.header.compact()
	copy(packet[:], append(header[:], p.payload[:]...))
	return packet
}

// newAddr creates a new struct Addr
func newAddr(addr *net.UDPAddr) *Addr {
	return &Addr{udpAddr: addr}
}

// newPacket decompresses a byte array into a Packet struct
func newPacket(packet [524]byte) *Packet {
	var header [12]byte
	copy(header[:], packet[0:12])

	var payload [512]byte
	copy(payload[:], packet[12:524])

	return &Packet{
		header:  newHeader(header),
		payload: payload,
	}
}

// FIXME prototype func
func newDataPacket(payload [512]byte) *Packet {
	header := newDataHeader(1, 1, 1);
	return &Packet{
		header:  newHeader(header.compact()),
		payload: payload,
	}
}

// String parses to string
func (a *Addr) String() string {
	return a.udpAddr.String()
}

// ResolveName parses a host name to IP/Port
func ResolveName(addr string) (*Addr, error) {
	debug("Resolving hostname")
	ServerAddr, err := net.ResolveUDPAddr("udp", addr)
	return newAddr(ServerAddr), err
}

// Listen listens to clients
func Listen(addr *Addr) (*Listener, error) {
	// TODO start receivePacket, initialize Listener
	udpConn, err := net.ListenUDP("udp", addr.udpAddr)
	if err != nil {
		return nil, err
	}

	l := newListener(addr, udpConn)
	l.receivePacket()

	return l, nil
}

// Connect connects to a server
func Connect(remoteAddr *Addr) (*Conn, error) {
	debug("Connecting to a server")
	localAddr, err := net.ResolveUDPAddr("udp", remoteAddr.addr.IP.String()+":0")
	checkError(err)

	conn, err := net.DialUDP("udp", localAddr, remoteAddr.addr)
	return newConn(conn), err
}

// receivePacket listens UDP packets and differentiates
func (l *Listener) receivePacket() {
	go func() {
		debug("Reading a packet")
		buffer := make([]byte, 524)
		n, addr, err := .ReadFromUDP(b)
	}();
}
