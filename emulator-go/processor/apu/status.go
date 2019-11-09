package apu

type statusRegister struct {
	dmcEnabled      bool
	noiseEnabled    bool
	triangleEnabled bool
	pulse2Enabled   bool
	pulse1Enabled   bool
}

func (s *statusRegister) write(data byte) {
	s.dmcEnabled = (data & 0b00010000) == 0b00010000
	s.noiseEnabled = (data & 0b00001000) == 0b00001000
	s.triangleEnabled = (data & 0b00000100) == 0b00000100
	s.pulse2Enabled = (data & 0b00000010) == 0b00000010
	s.pulse1Enabled = (data & 0b00000001) == 0b00000001
}

func (s *statusRegister) read(apu *APU) byte {
	// FIXME
	var result byte = 0
	if s.dmcEnabled {
		result |= 0b00010000
	}
	if s.noiseEnabled {
		result |= 0b00001000
	}
	if s.triangleEnabled {
		result |= 0b00000100
	}
	if apu.pulse2.lengthCounterLoad > 0 {
		result |= 0b00000010
	}
	if apu.pulse1.lengthCounterLoad > 0 {
		result |= 0b00000001
	}

	return result
}
