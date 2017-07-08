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
	payload    []byte
	sourceAddr *Addr
}

// compact compresses a packet into a byte array
func (p *Packet) compact() []byte {
	var packet [524]byte
	header := p.header.compact()

	copy(packet[:], append(header[:], p.payload[:]...))
	return packet[:12+len(p.payload)]
}

// newAddr creates a new struct Addr
func newAddr(addr *net.UDPAddr) *Addr {
	return &Addr{udpAddr: addr}
}

func newPacket(header *Header, payload []byte, sourceAddr *Addr) *Packet {
	return &Packet{
		header:     header,
		payload:    payload,
		sourceAddr: sourceAddr,
	}
}

// decompactPacket decompresses a byte array into a Packet struct
func decompactPacket(packet []byte, sourceAddr *Addr) *Packet {
	var header [12]byte
	copy(header[:], packet[0:12])

	var payload = make([]byte, 524)
	if len(packet) > 12 {
		copy(payload, packet[12:])
	}

	return &Packet{
		header:     decompactHeader(&header),
		payload:    payload[:len(packet)-12],
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
func Connect(remoteAddr *Addr) (*ConnServer, error) {
	debug("Connecting to a server")
	addr, err := net.ResolveUDPAddr("udp", remoteAddr.udpAddr.IP.String()+":0")
	if err != nil {
		return nil, err
	}
	localAddr := newAddr(addr)

	udpConn, err := net.DialUDP("udp", addr, remoteAddr.udpAddr)
	if err != nil {
		return nil, err
	}

	debug("Sending SYN packet")
	header := newHeader(12345, 0, 0,
		false, true, false)
	packet := newPacket(header, nil, localAddr)

	_, err = udpConn.Write(packet.compact())
	if err != nil {
		return nil, err
	}

	for !packet.header.syn || !packet.header.ack {
		debug("Waiting SYN-ACK packet")
		var response [524]byte
		n, addr, err := udpConn.ReadFromUDP(response[:])
		if err != nil {
			return nil, err
		}

		packet = decompactPacket(response[:n], newAddr(addr))
	}

	debug("Sending ACK packet")
	header = newHeader(
		packet.header.ackNum,
		packet.header.seqNum+1,
		packet.header.connID,
		true, false, false,
	)
	packet = newPacket(header, nil, localAddr)

	_, err = udpConn.Write(packet.compact()[:])
	if err != nil {
		return nil, err
	}
	debug("Handshaking DONE")

	conn := newConnServer(udpConn, remoteAddr, packet.header.connID)

	return conn, nil
}
