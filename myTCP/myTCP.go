package myTCP

import (
	"net"
	"sync"
)

// Addr is a wrapper for UDPAddr.
type Addr struct {
	addr *net.UDPAddr
}

// Listener is a MyTCP network listener.
type Listener struct {
	addr Addr
	conns map[string]Conn
	connsMutex sync.RWMutex
}

// Accept implements the Accept method in the net.Listener interface;
// it waits for the next client and returns a MyTCP Conn.
func (l *Listener) Accept() (Conn, error) {
//	TODO Accept: receive request from channel, create connection

	return nil, nil
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
	return &Addr{addr: addr}
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
	return a.addr.String()
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
	//ServerConn, err := net.ListenUDP("udp", addr.addr)
	//return newConn(ServerConn), err

	return nil, nil
}

// Connect connects to a server
func Connect(remoteAddr *Addr) (*Conn, error) {
	debug("Connecting to a server")
	localAddr, err := net.ResolveUDPAddr("udp", remoteAddr.addr.IP.String()+":0")
	checkError(err)

	conn, err := net.DialUDP("udp", localAddr, remoteAddr.addr)
	return newConn(conn), err
}
