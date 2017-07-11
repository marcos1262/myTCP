package myTCP

import (
	"net"
	"io"
	"strconv"
	"time"
	"errors"
)

type resultRead struct {
	n int
	err error
}

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

	c.getData <- p // send slice to be filled with data from network

	res := <-c.resultRead

	n = res.n
	err = res.err

	debug("READ ALL")

	// TODO catch some error
	return n, err
}

func (c *ConnClient) receivePacket() {
	go func() {
		expectedSeqNum := initialSeqNum

	waiting:
		for {
			select {
			case p := <-c.getData:
				var n int
			reading:
				for {
					select {
					case packet := <-c.newPacket: // packet arrived from network
						go func() {
							<-c.timeoutInactive // receive old timeout
						}()
						c.timeoutInactive = time.After(10 * time.Second)

						if packet.header.seqNum != expectedSeqNum {
							debug("E: " + strconv.Itoa(int(expectedSeqNum)) +
								", G: " + strconv.Itoa(int(packet.header.seqNum)))
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
						//debug("SENDING ACK ackNum: " + strconv.Itoa(int(expectedSeqNum)))
						ack := newACKPacket(packet.header.ackNum, expectedSeqNum, c.ID, c.listener.addr)
						_, err := writePacketToAddr(c.listener.udpConn, c.remoteAddr, ack)
						if err != nil {
							c.resultRead <- resultRead{n,err}
							break reading
						}

						// TODO keep payload rest, if it doesnt fit in p buffer

						if len(payload) < 512 {
							debug("LITTLE MESSAGE")
							c.resultRead <- resultRead{n,io.EOF}
							break reading
						}

						if n >= len(p) {
							debug("Buffer Completed")
							c.resultRead <- resultRead{n,nil}
							break reading
						}
					case <-c.timeoutInactive: // client has been inactive for a long time
						debug("TIMEOUT INACTIVE connID:" + strconv.Itoa(int(c.ID)))
						c.Close()
						c.resultRead <- resultRead{n,errors.New("Client inactive for 10s")}
						break reading
					}
				}

			case <-c.timeoutInactive: // client has been inactive for a long time
				debug("TIMEOUT INACTIVE connID:" + strconv.Itoa(int(c.ID)))
				c.Close()
				break waiting
			}
		}
	}()
}
