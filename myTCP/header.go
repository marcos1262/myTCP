package myTCP

import (
	"encoding/binary"
)

type Header struct {
	seqNum uint32
	ackNum uint32
	connID uint16
	ack    bool
	syn    bool
	fin    bool
}

func newHeader(seqNum uint32, ackNum uint32, connID uint16, ack bool, syn bool, fin bool) *Header {
	return &Header{
		seqNum: seqNum,
		ackNum: ackNum,
		connID: connID,
		ack:    ack,
		syn:    syn,
		fin:    fin,
	}
}

// Compacts a header into a byte array
func (h *Header) compact() *[12]byte {
	var compact [12]byte

	seqNumByte := compact[0:4]
	binary.LittleEndian.PutUint32(seqNumByte, h.seqNum)

	ackNumByte := compact[4:8]
	binary.LittleEndian.PutUint32(ackNumByte, h.ackNum)

	connIDByte := compact[8:10]
	binary.LittleEndian.PutUint16(connIDByte, h.connID)

	restByte := compact[10:12]
	var rest uint16 = 0
	if h.ack {
		rest = uint16(setBit(int(rest), 2))
	}
	if h.syn {
		rest = uint16(setBit(int(rest), 1))
	}
	if h.fin {
		rest = uint16(setBit(int(rest), 0))
	}

	binary.LittleEndian.PutUint16(restByte, rest)
	return &compact
}

// Decompacts a byte array into a Header
func decompactHeader(header *[12]byte) *Header {
	rest := binary.LittleEndian.Uint16((*header)[10:12])

	return &Header{
		seqNum: binary.LittleEndian.Uint32((*header)[0:4]),
		ackNum: binary.LittleEndian.Uint32((*header)[4:8]),
		connID: binary.LittleEndian.Uint16((*header)[8:10]),
		ack:    hasBit(int(rest), 2),
		syn:    hasBit(int(rest), 1),
		fin:    hasBit(int(rest), 0),
	}
}

// Sets the bit at pos in the integer n.
func setBit(n int, pos uint) int {
	return n | (1 << pos)
}

// Clears the bit at pos in n.
func clearBit(n int, pos uint) int {
	return n &^ (1 << pos)
}

// Verifies if there's a 1 bit at pos in n.
func hasBit(n int, pos uint) bool {
	return (n & (1 << pos)) > 0
}
