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
	payload    *[512]byte
	sourceAddr *Addr
}

// compact compresses a packet into a byte array
func (p *Packet) compact() *[524]byte {
	var packet [524]byte
	header := p.header.compact()

	copy(packet[:], append((*header)[:], (*p.payload)[:]...))
	return &packet
}

// newAddr creates a new struct Addr
func newAddr(addr *net.UDPAddr) *Addr {
	return &Addr{udpAddr: addr}
}

func newPacket(header *Header, payload *[512]byte, sourceAddr *Addr) *Packet {
	return &Packet{
		header:     header,
		payload:    payload,
		sourceAddr: sourceAddr,
	}
}

// decompactPacket decompresses a byte array into a Packet struct
func decompactPacket(packet *[524]byte, sourceAddr *Addr) *Packet {
	var header [12]byte
	copy(header[:], (*packet)[0:12])

	var payload [512]byte
	copy(payload[:], (*packet)[12:524])

	return &Packet{
		header:     decompactHeader(&header),
		payload:    &payload,
		sourceAddr: sourceAddr,
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
	addr, err := net.ResolveUDPAddr("udp", remoteAddr.udpAddr.IP.String()+":0")
	if err != nil {
		return nil, err
	}
	localAddr := newAddr(addr)

	conn, err := net.DialUDP("udp", addr, remoteAddr.udpAddr)
	if err != nil {
		return nil, err
	}

	debug("Sending SYN packet")
	header := newHeader(12345,
		0,
		0,
		false, true, false)
	var emptyPayload [512]byte = [512]byte{}
	packet := newPacket(header, &emptyPayload, localAddr)

	_, err = conn.Write((*packet.compact())[:])
	if err != nil {
		return nil, err
	}

	for !packet.header.syn || !packet.header.ack {
		debug("Waiting SYN-ACK packet")
		var response [524]byte
		_, addr, err = conn.ReadFromUDP(response[:])
		if err != nil {
			return nil, err
		}

		packet = decompactPacket(&response, newAddr(addr))
	}

	debug("Sending ACK packet")
	header = newHeader(
		packet.header.ackNum,
		packet.header.seqNum+1,
		packet.header.connID,
		true, false, false,
	)
	packet = newPacket(header, &emptyPayload, localAddr)

	_, err = conn.Write((*packet.compact())[:])
	if err != nil {
		return nil, err
	}
	debug("Handshaking DONE")

	return newConn(conn, remoteAddr), nil
}
