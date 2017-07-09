package myTCP

import (
	"strconv"
	"time"
	"net"
	"math"
)

const windowSize uint32 = 512

type window struct {
	base       uint32
	nextSeqNum uint32
	size       uint32
	data       []byte
	timeout    <-chan time.Time
}

func (c *ConnServer) sendPacket() {
	// TODO get from channel, call writePacket
	go func() {

		for {
			select {
			case packet := <-c.outPacket:
				// store in window
			case <-c.timeoutAck:
				// retransmit
			default:
				// if is possible to send (nextseqnum < base+size)
			}
		}
	}()
}

func writePacket(udpConn *net.UDPConn, packet *Packet) (int, error) {
	return udpConn.Write(packet.compact())
}

// Write a packet to a connection
func (c ConnServer) Write(p []byte) (n int, err error) {
	debug("WRITING " + strconv.Itoa(len(p)) + " bytes : " + string(p))
	// TODO Write: coordinate sending of packets, break bytes on packets

	data := p

	// TODO receive local addr from client, or catch from OS
	localAddr, err := ResolveName("127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	for len(data) { // FIXME verify if > 0

		// FIXME remove
		//payload := make([]byte, len(p))
		//copy(payload[:], p)

		// FIXME remove
		//qtd_wrote, err := c.udpConn.Write(packet)
		//n += qtd_wrote - 12
		//if err != nil {
		//	return 0, err
		//}

		qtdWrite := int(math.Min(512, float64(len(data))))

		packet := newDataPacket(seqNumInitialClient, c.ID,
			data[:qtdWrite], localAddr)

		c.outPacket <- packet

		// FIXME remove
		if len(data) <= 512 {
			debug("LITTLE MESSAGE")
		} else {
			debug("BIG MESSAGE")
		}

		n += qtdWrite
		data = data[qtdWrite:]
	}

	return n, nil
}
