package myTCP

import (
	"sync"
	"net"
)

// Listener is a MyTCP network listener.
type Listener struct {
	addr       *Addr
	udpConn    *net.UDPConn
	newConn    chan *Packet
	conns      map[uint16]*Conn
	connsMutex *sync.RWMutex
}

// newAddr creates a new struct Addr.
func newListener(addr *Addr, udpConn *net.UDPConn) *Listener {
	return &Listener{
		addr:    addr,
		udpConn: udpConn,
		newConn: make(chan *Packet),
		conns:   make(map[uint16]*Conn),
	}
}

// Accept implements the Accept method in the net.Listener interface;
// Accept waits for the next connection on a channel and returns it to the listener.
func (l *Listener) Accept() (*Conn, error) {
	for packet := range l.newConn {
		//	TODO Accept: handshake

	}
	return nil, nil
}

// Close stops listening on the TCP address.
// Already Accepted connections are not closed.
func (l *Listener) Close() {
	close(l.newConn)
}

// searchConn searches for a saved connection by the connID.
func (l *Listener) searchConn(connID uint16) (*Conn, bool) {
	l.connsMutex.RLock() // Mutual exclusion (reading)
	conn, exists := l.conns[connID]
	l.connsMutex.RUnlock()

	return conn, exists
}

// deleteConn removes a saved connection, searching by the connID.
func (l *Listener) deleteConn(connID uint16) {
	l.connsMutex.Lock() // Exclusão mútua (escrita)
	delete(l.conns, connID)
	l.connsMutex.Unlock()
}

// listenPacket listens UDP packets.
func (l *Listener) listenPacket() {
	go func() {
		for {
			buffer := make([]byte, 524)
			n, addr, err := l.udpConn.ReadFromUDP(buffer)
			checkError(err)

			debug("Read a packet")

			var packetByte [524]byte
			copy(packetByte[:], buffer[:n])

			go l.receivePacket(packetByte, newAddr(addr))
		}
	}()
}

// receivePacket differentiates the received packet and forwards it to the right place.
func (l *Listener) receivePacket(packetByte [524]byte, addr *Addr) {
	packet := newPacket(packetByte, addr)

	if packet.header.syn {
		l.newConn <- packet
	} else {
		conn, exists := l.searchConn(packet.header.connID)
		if exists {
			conn.newPacket <- packet
		} else {
			// PACKET WITHOUT SYN WITHOUT KNOWN ID
			// IGNORE?
		}
	}
}
