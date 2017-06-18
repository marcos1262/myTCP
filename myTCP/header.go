package myTCP

import (
	"encoding/binary"
)

type Header struct {
	seqNum uint32
	ackNum uint32
	connID uint16
	ACK    bool
	SYN    bool
	FIN    bool
}

func (h *Header) compact() []byte {
	compact := make([]byte, 12)

	seqNumByte := compact[0:4]
	binary.LittleEndian.PutUint32(seqNumByte, h.seqNum)

	ackNumByte := compact[4:8]
	binary.LittleEndian.PutUint32(ackNumByte, h.ackNum)

	connIDByte := compact[8:10]
	binary.LittleEndian.PutUint16(connIDByte, h.connID)

	restByte := compact[10:12]
	var rest uint16 = 0
	if h.ACK {
		setBit(int(rest), 2)
	}
	if h.SYN {
		setBit(int(rest), 1)
	}
	if h.FIN {
		setBit(int(rest), 0)
	}
	binary.LittleEndian.PutUint16(restByte, rest)
	return compact
}

func decompact(header [12]byte) *Header {
	rest := binary.LittleEndian.Uint16(header[10:12])

	return &Header{
		seqNum: binary.LittleEndian.Uint32(header[0:4]),
		ackNum:binary.LittleEndian.Uint32(header[4:8]),
		connID:binary.LittleEndian.Uint16(header[8:10]),
		ACK:hasBit(int(rest), 2),
		SYN:hasBit(int(rest), 1),
		FIN:hasBit(int(rest), 0),
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
