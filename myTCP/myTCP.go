package myTCP

import (
	"net"
)

// Addr is a wrapper for UDPAddr.
type Addr struct {
	udpAddr *net.UDPAddr
}

// newAddr creates a new struct Addr
func newAddr(addr *net.UDPAddr) *Addr {
	return &Addr{udpAddr: addr}
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

	listenPacket(l.udpConn, l.demultiplexer)

	return l, nil
}

// Connect connects to a server
func Connect(remoteAddr *Addr) (*ConnServer, error) {
	debug("Connecting to a server")

	localAddr, err := ResolveName(remoteAddr.udpAddr.IP.String() + ":0")
	if err != nil {
		return nil, err
	}

	udpConn, err := net.DialUDP("udp", localAddr.udpAddr, remoteAddr.udpAddr)
	if err != nil {
		return nil, err
	}

	debug("Sending SYN packet")
	_, err = writePacket(udpConn, newSYNPacket(localAddr))
	if err != nil {
		return nil, err
	}

	var packet *Packet
	for !packet.header.syn || !packet.header.ack { // FIXME verify packet == nil
		debug("Waiting SYN-ACK packet")

		// FIXME remove
		//var response [524]byte
		//n, addr, err := udpConn.ReadFromUDP(response[:])
		//if err != nil {
		//	return nil, err
		//}
		//packet = decompactPacket(response[:n], newAddr(addr))

		packet, err = readPacket(udpConn)
		if err != nil {
			return nil, err
		}
	}

	debug("Sending ACK packet")
	packet = newACKPacket(packet.header.ackNum, packet.header.seqNum+1,
		packet.header.connID, localAddr)

	_, err = writePacket(udpConn, packet)
	if err != nil {
		return nil, err
	}
	debug("Handshaking DONE")

	conn := newConnServer(udpConn, remoteAddr, packet.header.connID)
	conn.sendPacket()

	return conn, nil
}
