package myTCP

import (
	"net"
)

// Addr is a wrapper for UDPAddr.
type Addr struct {
	udpAddr *net.UDPAddr
}

// Packet represents structure
type Packet struct {
	header     *Header
	payload    [512]byte
	sourceAddr *Addr
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
func newPacket(packet [524]byte, sourceAddr *Addr) *Packet {
	var header [12]byte
	copy(header[:], packet[0:12])

	var payload [512]byte
	copy(payload[:], packet[12:524])

	return &Packet{
		header:     newHeader(header),
		payload:    payload,
		sourceAddr: sourceAddr,
	}
}

// FIXME prototype func
//func newDataPacket(payload [512]byte) *Packet {
//	header := newDataHeader(1, 1, 1);
//	return &Packet{
//		header:  newHeader(header.compact()),
//		payload: payload,
//	}
//}

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
	debug("Listening to packets on " + addr.String())

	udpConn, err := net.ListenUDP("udp", addr.udpAddr)
	if err != nil {
		return nil, err
	}

	l := newListener(addr, udpConn)
	l.listenPacket()

	return l, nil
}

// Connect connects to a server
func Connect(remoteAddr *Addr) (*Conn, error) {
	debug("Connecting to a server")
	localAddr, err := net.ResolveUDPAddr("udp", remoteAddr.udpAddr.IP.String()+":0")
	checkError(err)

	conn, err := net.DialUDP("udp", localAddr, remoteAddr.udpAddr)
	return newConn(conn), err
}