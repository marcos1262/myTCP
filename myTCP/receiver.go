package myTCP

import (
	"net"
	"io"
	"strconv"
	"time"
	"errors"
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

	expectedSeqNum := initialSeqNum

reading:
	for {
		select {
		case packet := <-c.newPacket:
			go func() {
				<-c.timeoutInactive // receive old timeout
			}()
			c.timeoutInactive = time.After(10 * time.Second)

			if packet.header.seqNum != expectedSeqNum {
				debug("E: " + strconv.Itoa(int(expectedSeqNum)) + ", G: " + strconv.Itoa(int(packet.header.seqNum)))
				debug("Discarding packet out of order...")
				break
			}

			payload := packet.payload

			// TODO detect fin and return EOF
			//if len(payload) == 0{
			//	break
			//}

			copy(p[n:n+len(payload)], payload)

			n += len(payload)

			expectedSeqNum = packet.header.seqNum + uint32(len(payload))

			//debug("seqNum received: " + strconv.Itoa(int(packet.header.seqNum)))
			//debug("SENDING ACK ackNum:" + strconv.Itoa(int(expectedSeqNum)))
			ack := newACKPacket(packet.header.ackNum, expectedSeqNum, c.ID, c.listener.addr)
			_, err := writePacketToAddr(c.listener.udpConn, c.remoteAddr, ack)
			if err != nil {
				return 0, err
			}

			// TODO keep payload rest, if it doesnt fit in p buffer

			if len(payload) < 512 {
				debug("LITTLE MESSAGE")
				err = io.EOF
				return n, err
			}

			if n >= len(p) {
				debug("RECEIVED ALL")
				break reading
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
		case <-c.timeoutInactive:
			debug("TIMEOUT INACTIVE connID:" + strconv.Itoa(int(c.ID)))
			c.Close()
			err = errors.New("Client inactive for 10s")
			break reading
		}
	}

	// TODO verifies when Copy func stops reading
	//c.Close()

	// TODO catch some error
	return n, err
}
