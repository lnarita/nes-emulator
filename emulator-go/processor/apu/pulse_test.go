package apu

import "testing"

func TestReadByte1(t *testing.T) {
	p := pulseRegister{}
	byte1 := p.readByte1()
	if byte1 != 0 {
		t.Errorf("Expected result of 0, got %d", byte1)
	}

	// duty           byte // DDLC VVVV
	// envelopeLoop   bool
	// constantVolume bool
	// volume         byte
	p.duty = 0b11
	p.envelopeLoop = true
	p.constantVolume = true
	p.volume = 0b1111
	byte1 = p.readByte1()

	if byte1 != 0b11111111 {
		t.Errorf("Expected result of 0b11111111, got %d", byte1)
	}

	p.duty = 0b10
	p.envelopeLoop = true
	p.constantVolume = false
	p.volume = 0b1010
	byte1 = p.readByte1()

	if byte1 != 0b10101010 {
		t.Errorf("Expected result of 0b10101010, got %d", byte1)
	}
}
