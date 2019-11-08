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

func TestWriteByte1(t *testing.T) {
	p := pulseRegister{}
	p.writeByte1(0b11111111)

	if p.duty != 0b11 {
		t.Errorf("Expected duty of 0b11, got %d", p.duty)
	}

	if !p.envelopeLoop {
		t.Errorf("Expected envelopeLoop true, got %t", p.envelopeLoop)
	}

	if !p.constantVolume {
		t.Errorf("Expected constantVolume true, got %t", p.constantVolume)
	}

	if p.volume != 0b1111 {
		t.Errorf("Expected volume of 0b1111, got %d", p.volume)
	}
}
