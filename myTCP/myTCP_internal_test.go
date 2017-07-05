package myTCP

import (
	"testing"
	"encoding/binary"
	"net"
)

// TODO test internal funcs

func TestHeaderCompact(t *testing.T) {
	h := &Header{1, 2, 3, true, true, false}
	c := h.compact()
	var sN, aN, cID, flags = binary.LittleEndian.Uint32(c[0:4]),
		binary.LittleEndian.Uint32(c[4:8]),
		binary.LittleEndian.Uint16(c[8:10]),
		binary.LittleEndian.Uint16(c[10:12])

	if sN != 1 {
		t.Errorf("Error on compacting header seqNum. Expected: 1, got: %d", sN)
	}
	if aN != 2 {
		t.Errorf("Error on compacting header ackNum. Expected: 2, got: %d", aN)
	}
	if cID != 3 {
		t.Errorf("Error on compacting header connID. Expected: 3, got: %d", cID)
	}
	if flags != 6 {
		t.Errorf("Error on compacting header flags. Expected: 6, got: %d", flags)
	}
}

func TestHeaderDecompact(t *testing.T) {
	h := &Header{1, 2, 3, true, true, false}
	c := h.compact()
	d := decompactHeader(c)
	if d.seqNum != 1 {
		t.Errorf("Error on decompacting header seqNum. Expected: 1, got: %d", d.seqNum)
	}
	if d.ackNum != 2 {
		t.Errorf("Error on decompacting header ackNum. Expected: 2, got: %d", d.ackNum)
	}
	if d.connID != 3 {
		t.Errorf("Error on decompacting header connID. Expected: 3, got: %d", d.connID)
	}
	if d.ack != true {
		t.Errorf("Error on decompacting header ack. Expected: 6, got: %d", d.ack)
	}

	if d.syn != true {
		t.Errorf("Error on decompacting header syn. Expected: 6, got: %d", d.syn)
	}

	if d.fin != false {
		t.Errorf("Error on decompacting header fin. Expected: 6, got: %d", d.fin)
	}
}

//func TestReadPacket(t *testing.T) {
//	var addr, _ = ResolveName("127.0.0.1:12345")
//	udpConn, _ := net.ListenUDP("udp", addr.udpAddr)
//
//	var l = newListener(addr, udpConn)
//	var emptyPayload [512]byte = [512]byte{}
//	header := newHeader(
//		4321,
//		4322,
//		42,
//		true, true, false,
//	)
//	packet := newPacket(header, &emptyPayload, l.addr)
//
//	go func() {
//		fmt.Println("wangga")
//		//p, addr := l.readPacket()
//		//var d = decompactPacket(p, addr)
//		//if addr.String() != "127.0.0.1:12345" {
//		//	t.Errorf("Error on reading packet from network. Wrong Packet ADDR.")
//		//}
//		//if d.header.ackNum != 4322 {
//		//	t.Errorf("Error on reading packet. Wrong Packet Attr.")
//		//}
//		//
//		//l.Close()
//		//return
//	}()
//
//	la, _ := ResolveName("127.0.0.1:0")
//	conn, _ := net.DialUDP("udp", la.udpAddr, addr.udpAddr)
//
//	var p = packet.compact()
//	conn.Write(p[:])
//}

func TestReceivePacket(t *testing.T) {
	var addr, _ = ResolveName("127.0.0.1:12345")
	udpConn, _ := net.ListenUDP("udp", addr.udpAddr)

	var l = newListener(addr, udpConn)

	var emptyPayload [512]byte = [512]byte{}
	header := newHeader(
		4321,
		4322,
		42,
		true, false, false,
	)
	packet := newPacket(header, &emptyPayload, l.addr)

	go func() {
		p := <-l.newConn
		if p.header.syn != true {
			t.Errorf("Error on receiving packet. Packet SYN not forwarded by channel.")
		}
		return
	}()

	l.receivePacket(packet.compact(), addr)

	conn := newConn(nil, addr)
	l.conns[1] = conn
	header = newHeader(
		4321,
		0,
		42,
		false, false, false,
	)
	packet = newPacket(header, &emptyPayload, l.addr)

	go func() {
		p := <-conn.newPacket
		if p.header.syn || p.header.ackNum == 4322 {
			t.Errorf("Error on receiving packet. Data Packet not forwarded by channel.")
		}
		return
	}()

	l.receivePacket(packet.compact(), addr)
	//l.Close()
}

func TestListenPacket(t *testing.T) {
	//func (l *Listener) listenPacket() {
	//	go func() {
	//		for {
	//			packetByte, addr := l.readPacket()
	//
	//			go l.receivePacket(packetByte, addr)
	//		}
	//	}()
	//}

}

func TestPacketCompact(t *testing.T) {
	//func (p *Packet) compact() *[524]byte {
	//	var packet [524]byte
	//	header := p.header.compact()
	//
	//	copy(packet[:], append((*header)[:], (*p.payload)[:]...))
	//	return &packet
	//}

}

func TestPacketDecompact(t *testing.T) {
	//	func decompactPacket(packet *[524]byte, sourceAddr *Addr) *Packet {
	//		var header [12]byte
	//		copy(header[:], (*packet)[0:12])
	//
	//		var payload [512]byte
	//		copy(payload[:], (*packet)[12:524])
	//
	//		return &Packet{
	//			header:     decompactHeader(&header),
	//			payload:    &payload,
	//			sourceAddr: sourceAddr,
	//		}
	//}

}
