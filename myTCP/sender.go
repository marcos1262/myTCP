package myTCP

import (
	"strconv"
	"time"
	"net"
	"math"
)

var cwnd uint32 = 10240
var timeoutInterval = 5 * time.Millisecond
var initialSeqNum = seqNumInitialClient + 1 + 1 // seqNum after sending ACK of SYNACK

type segment struct {
	*Packet
	acked bool
}

type window struct {
	base       uint32
	nextSeqNum uint32
	size       uint32
	data       []byte
	segments   []*segment
	timeoutAck <-chan time.Time
}

func (c *ConnServer) sendPacket() {
	// TODO get from channel, call writePacketToConn
	go func() {
		c.window = window{
			base:       initialSeqNum,
			nextSeqNum: initialSeqNum,
			size:       cwnd,
			data:       []byte{}, // create infinite size slice
		}
		w := &c.window

	sending:
		for {
			select {
			case packetByte := <-c.outPacket: // append new data
				c.window.data = append(c.window.data, packetByte...)
			case packet := <-c.ackPacket: // move window pointers
				//debug("									RECEIVED ACK ackNum:" + strconv.Itoa(int(packet.header.ackNum)))
				if packet.header.ackNum > w.base { // expected ack
					w.base = packet.header.ackNum
					if w.base < w.nextSeqNum || w.base < uint32(len(w.data))+initialSeqNum {
						// there are currently not-yet-acked segments
						w.timeoutAck = time.After(timeoutInterval)
					} else {
						w.timeoutAck = nil
						debug("ENDING SENDING")
						c.resultWrite <- nil
						break sending
					}
				} else { // duplicated ack

				}
			case <-w.timeoutAck: // retransmit
				debug("TIMEOUT END segment:" + strconv.Itoa(int(w.base)))
				indexNotAckedPacket := (w.base - initialSeqNum) / 512

				notAckedPacket := w.segments[indexNotAckedPacket].Packet
				debug("Retransmitting seqNum: " + strconv.Itoa(int(notAckedPacket.header.seqNum)))
				writePacketToConn(c.udpConn, notAckedPacket)

				w.timeoutAck = time.After(timeoutInterval)
			default:

				//limit := uint32(math.Min(float64(uint32(len(w.data))+initialSeqNum-w.nextSeqNum), float64(w.size)))
				limit := uint32(math.Min(float64(uint32(len(w.data))+initialSeqNum-w.base), float64(w.size)))

				if w.nextSeqNum < w.base+uint32(limit) { // if is possible to send

					//debug("WINDOW{base: " + strconv.Itoa(int(w.base)) +
					//	", nextSeqNum: " + strconv.Itoa(int(w.nextSeqNum)) +
					//	", size: " + strconv.Itoa(int(w.size)) +
					//	", dataLen: " + strconv.Itoa(len(w.data)) +
					//	"}\n" +
					//	"LIMIT:" + strconv.Itoa(int(limit)),
					//)

					//debug("SENDING DATA seqNum:" + strconv.Itoa(int(w.nextSeqNum)))
					qtdWrite := uint32(math.Min(512, float64(w.base+limit-w.nextSeqNum)))

					// send packet to network
					begin := w.nextSeqNum - initialSeqNum
					end := begin + qtdWrite
					//debug("POINTERS: " + strconv.Itoa(int(begin)) + "," + strconv.Itoa(int(end)))
					packet := newDataPacket(w.nextSeqNum, c.ID,
						w.data[begin:end], c.localAddr)
					//debug("DATA: " + string(w.data[begin:end]))
					_, err := writePacketToConn(c.udpConn, packet)
					if err != nil {
						c.resultWrite <- err
						break sending
					}

					// save segment sent for possible retransmitting
					segment := &segment{packet, false}
					w.segments = append(w.segments, segment)

					// point to next byte to send
					w.nextSeqNum += qtdWrite

					if w.timeoutAck == nil { // if timeout not set
						w.timeoutAck = time.After(timeoutInterval)
					}

					if end == uint32(len(w.data)) {
						// sent all data
					}
				}

				// FIXME remove this DEBUG
				//time.Sleep(3 * time.Second)
			}
		}
	}()
}

func writePacketToConn(udpConn *net.UDPConn, packet *Packet) (int, error) {
	return udpConn.Write(packet.compact())
}

func writePacketToAddr(sourceConn *net.UDPConn, addr *Addr, packet *Packet) (int, error) {
	return sourceConn.WriteToUDP(packet.compact(), addr.udpAddr)
}

// Write a packet to a connection
func (c ConnServer) Write(p []byte) (n int, err error) {
	//debug("WRITING " + strconv.Itoa(len(p)) + " bytes : " + string(p))
	// TODO Write: coordinate sending of packets, break bytes on packets

	c.outPacket <- p // put data to send in queue

	err = <-c.resultWrite

	n = len(p)

	debug("WROTE ALL")

	return n, err
}

// Differentiate the received packet and forwards it to the right place.
func (c *ConnServer) demultiplexer(packet *Packet) {
	if packet.header.ack {
		c.ackPacket <- packet
	} else {
		if packet.sourceAddr == c.remoteAddr {

			// TODO plan a way of unify the sending and receiving of data in both conn types

			// TODO allow ConnServer to receive data packets
			//c.newPacket <- packet
		} else {
			// PACKET WITHOUT SYN WITHOUT KNOWN ID
			// IGNORE?
		}
	}
}
