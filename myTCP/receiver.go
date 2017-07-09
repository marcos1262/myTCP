package myTCP

import (
	"net"
	"io"
	"strconv"
)

// FIXME remove
//// Listen UDP packets.
//func (l *Listener) listenPacket() {
//	go func() {
//		for {
//			packet := readPacket(l.udpConn)
//
//			go l.demultiplexer(packet)
//		}
//	}()
//}

// Listen UDP packets.
func listenPacket(udpConn *net.UDPConn, demux func(*Packet)) {
	go func() {
		for {
			packet, err := readPacket(udpConn)
			checkError(err)

			go demux(packet)
		}
	}()
}

// Get packet from network.
func readPacket(udpConn *net.UDPConn) (*Packet, error) {
	buffer := make([]byte, 524)
	n, addr, err := udpConn.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}

	packet := decompactPacket(buffer[:n], newAddr(addr))

	return packet, nil
}

// Take a packet from a connection, copying the payload into p
func (c ConnClient) Read(p []byte) (n int, err error) {
	debug("READING " + strconv.Itoa(len(p)) + " bytes")
	//	TODO Read: coordinate receiving of packets,
	// 				create packet/payload BUFFER,
	// 				return only until the requested size

	for {
		select {
		case packet := <-c.newPacket:
			payload := packet.payload

			// TODO detect fin and return EOF
			//if len(payload) == 0{
			//	break
			//}

			copy(p[n:n+len(payload)], payload)

			debug("A: " + string(payload))
			debug("B: " + string(p[n:n+len(payload)]))

			n += len(payload)

			// TODO send ACK

			// TODO keep payload rest, if it doesnt fit in p buffer

			if len(payload) < 512 {
				debug("LITTLE MESSAGE")
				err = io.EOF
				break
			}

			if n >= len(p) {
				break
			}
			//if len(p) <= 512 {
			//	debug("LITTLE MESSAGE")
			//
			//} else {
			//	debug("BIG MESSAGE")
			//
			//
			//
			//	panic("IIIIIIIIIII")
			//}
		case <-c.timeoutInative:
			c.Close()
		}
	}

	debug("OK")

	// TODO verifies when Copy func stops reading
	//c.Close()

	// TODO catch some error
	return n, err
}
