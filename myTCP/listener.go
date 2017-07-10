package myTCP

import (
	"sync"
	"net"
	"time"
)

// Listener is a MyTCP network listener.
type Listener struct {
	addr       *Addr
	udpConn    *net.UDPConn
	newConn    chan *Packet
	conns      map[uint16]*ConnClient
	connsMutex *sync.RWMutex
}

// Create a new struct Listener.
func newListener(addr *Addr, udpConn *net.UDPConn) *Listener {
	return &Listener{
		addr:       addr,
		udpConn:    udpConn,
		newConn:    make(chan *Packet),
		conns:      make(map[uint16]*ConnClient),
		connsMutex: &sync.RWMutex{},
	}
}

// Wait for the next connection on a channel and return it to the listener.
func (l *Listener) Accept() (ConnClient, error) {
	var conn ConnClient
	// BUG test when two clients request at same time
	for packet := range l.newConn {
		if packet.header.syn {
			debug("Receiving SYN packet")
			conn = *newConnClient(packet.sourceAddr, l)

			debug("Sending SYN-ACK packet")
			packet = newSYNACKPacket(conn.ID, l.addr)
			_, err := l.udpConn.WriteToUDP(packet.compact(), conn.remoteAddr.udpAddr)
			if err != nil {
				return ConnClient{}, err
			}
		} else if packet.header.ackNum == seqNumInitialServer+1 {
			debug("Receiving ACK packet")

			l.saveConn(&conn)
			debug("Handshaking DONE")
			return conn, nil
		}
	}
	return ConnClient{}, nil
}

// Stop listening on the TCP address.
func (l *Listener) Close() {
	close(l.newConn)
	l.udpConn.Close()
}

// Search for a saved connection by the connID.
func (l *Listener) searchConn(connID uint16) (*ConnClient, bool) {
	l.connsMutex.RLock() // Mutual exclusion (reading)
	conn, exists := l.conns[connID]
	l.connsMutex.RUnlock()

	return conn, exists
}

// Save a connection on the list.
func (l *Listener) saveConn(conn *ConnClient) {
	l.connsMutex.Lock() // Mutual exclusion (writing)
	l.conns[conn.ID] = conn
	l.connsMutex.Unlock()
}

// Remove a saved connection, searching by the connID.
func (l *Listener) deleteConn(connID uint16) {
	l.connsMutex.Lock() // Mutual exclusion (writing)
	delete(l.conns, connID)
	l.connsMutex.Unlock()
}

// Differentiate the received packet and forwards it to the right place.
func (l *Listener) demultiplexer(packet *Packet) {
	if packet.header.syn || packet.header.ackNum == seqNumInitialServer+1 {
		l.newConn <- packet
	} else {
		conn, exists := l.searchConn(packet.header.connID)
		if exists {
			conn.timeoutInactive = time.After(10 * time.Second)
			conn.newPacket <- packet
		} else {
			// PACKET WITHOUT SYN WITHOUT KNOWN ID
			// IGNORE?
		}
	}
}
