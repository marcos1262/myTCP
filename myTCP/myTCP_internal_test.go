package myTCP

import (
	"testing"
	"encoding/binary"
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
