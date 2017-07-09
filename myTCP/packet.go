package myTCP

const seqNumInitialClient uint32 = 12345
const seqNumInitialServer uint32 = 4321

// Packet represents a structured packet
type Packet struct {
	header     *Header
	payload    []byte
	sourceAddr *Addr
}

func newPacket(header *Header, payload []byte, sourceAddr *Addr) *Packet {
	return &Packet{
		header:     header,
		payload:    payload,
		sourceAddr: sourceAddr,
	}
}

func newSYNPacket(sourceAddr *Addr) *Packet {
	header := newHeader(seqNumInitialClient, 0, 0,
		false, true, false)
	return newPacket(header, nil, sourceAddr)
}

func newSYNACKPacket(connID uint16, sourceAddr *Addr) *Packet {
	header := newHeader(
		seqNumInitialServer,
		seqNumInitialClient+1,
		connID,
		true, true, false,
	)
	return newPacket(header, nil, sourceAddr)
}

func newACKPacket(seqNum uint32, ackNum uint32, connID uint16, sourceAddr *Addr) *Packet {
	header := newHeader(seqNum, ackNum, connID,
		true, false, false)
	return newPacket(header, nil, sourceAddr)
}

func newDataPacket(seqNum uint32, connID uint16, payload []byte, sourceAddr *Addr) *Packet {
	header := newHeader(seqNum, 0, connID,
		false, false, false)
	return newPacket(header, payload, sourceAddr)
}

// decompactPacket decompresses a byte array into a Packet struct
func decompactPacket(packet []byte, sourceAddr *Addr) *Packet {
	var header [12]byte
	copy(header[:], packet[0:12])

	var payload = make([]byte, 524)
	if len(packet) > 12 {
		copy(payload, packet[12:])
	}

	return &Packet{
		header:     decompactHeader(&header),
		payload:    payload[:len(packet)-12],
		sourceAddr: sourceAddr,
	}
}

// compact compresses a packet into a byte array
func (p *Packet) compact() []byte {
	var packet [524]byte
	header := p.header.compact()

	copy(packet[:], append(header[:], p.payload[:]...))
	return packet[:12+len(p.payload)]
}
