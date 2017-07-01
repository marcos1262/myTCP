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
		addr:       addr,
		udpConn:    udpConn,
		newConn:    make(chan *Packet),
		conns:      make(map[uint16]*Conn),
		connsMutex: &sync.RWMutex{},
	}
}

// Accept implements the Accept method in the net.Listener interface;
// Accept waits for the next connection on a channel and returns it to the listener.
func (l *Listener) Accept() (*Conn, error) {
	var conn *Conn
	for packet := range l.newConn {
		if packet.header.syn {
			debug("Receiving SYN packet")
			var emptyPayload [512]byte = [512]byte{}
			conn := newConn(nil, packet.sourceAddr)

			debug("Sending SYN-ACK packet")
			header := newHeader(
				4321,
				packet.header.seqNum+1,
				conn.ID,
				true, true, false,
			)
			packet = newPacket(header, &emptyPayload, l.addr)

			_, err := l.udpConn.WriteToUDP((*packet.compact())[:], conn.remoteAddr.udpAddr)
			if err != nil {
				return nil, err
			}
		} else if packet.header.ackNum == 4321+1 {
			debug("Receiving ACK packet")

			l.saveConn(conn)
			debug("Handshaking DONE")

			break;
		}
	}
	return conn, nil
}

// Close stops listening on the TCP address.
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

//
func (l *Listener) saveConn(conn *Conn) {
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

			go l.receivePacket(packetByte, addr)
		}
	}()
}

func (l *Listener) readPacket() (*[524]byte, *Addr) {
	buffer := make([]byte, 524)
	n, addr, err := l.udpConn.ReadFromUDP(buffer)
	checkError(err)

	debug("Read a packet OK")

	var packetByte [524]byte
	copy(packetByte[:], buffer[:n])

	return &packetByte, newAddr(addr)
}

// receivePacket differentiates the received packet and forwards it to the right place.
func (l *Listener) receivePacket(packetByte *[524]byte, addr *Addr) {
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
