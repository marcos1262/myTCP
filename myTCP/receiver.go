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
	conns      map[uint16]*ConnClient
	connsMutex *sync.RWMutex
}

// newAddr creates a new struct Addr.
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
	for packet := range l.newConn {
		if packet.header.syn {
			debug("Receiving SYN packet")
			conn = *newConnClient(packet.sourceAddr)

			debug("Sending SYN-ACK packet")
			header := newHeader(
				4321,
				packet.header.seqNum+1,
				conn.ID,
				true, true, false,
			)
			packet = newPacket(header, nil, l.addr)

			_, err := l.udpConn.WriteToUDP(packet.compact(), conn.remoteAddr.udpAddr)
			if err != nil {
				return ConnClient{}, err
			}
		} else if packet.header.ackNum == 4321+1 {
			debug("Receiving ACK packet")

			l.saveConn(&conn)
			debug("Handshaking DONE")
			return conn, nil
		}
	}
	return ConnClient{}, nil
}

// Close stops listening on the TCP address.
func (l *Listener) Close() {
	close(l.newConn)
	l.udpConn.Close()
}

// searchConn searches for a saved connection by the connID.
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

// deleteConn removes a saved connection, searching by the connID.
func (l *Listener) deleteConn(connID uint16) {
	l.connsMutex.Lock() // Mutual exclusion (writing)
	delete(l.conns, connID)
	l.connsMutex.Unlock()
}

// listenPacket listens UDP packets.
func (l *Listener) listenPacket() {
	go func() {
		for {
			packetByte, addr := l.readPacket()

			go l.demultiplexer(packetByte, addr)
		}
	}()
}

// Get packet from network.
func (l *Listener) readPacket() ([]byte, *Addr) {
	buffer := make([]byte, 524)
	n, addr, err := l.udpConn.ReadFromUDP(buffer)
	checkError(err)

	return buffer[:n], newAddr(addr)
}

// Differentiate the received packet and forwards it to the right place.
func (l *Listener) demultiplexer(packetByte []byte, addr *Addr) {
	packet := decompactPacket(packetByte, addr)

	if packet.header.syn || packet.header.ackNum == 4321+1 {
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
