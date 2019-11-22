package apu

import "testing"

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
