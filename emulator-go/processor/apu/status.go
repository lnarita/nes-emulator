package apu

type statusRegister struct {
	dmcEnabled      bool
	noiseEnabled    bool
	triangleEnabled bool
	pulse2Enabled   bool
	pulse1Enabled   bool
}

func (s *statusRegister) write(data byte, apu *APU) {
	s.dmcEnabled = (data & 0b00010000) == 0b00010000
	s.noiseEnabled = (data & 0b00001000) == 0b00001000
	s.triangleEnabled = (data & 0b00000100) == 0b00000100
	s.pulse2Enabled = (data & 0b00000010) == 0b00000010
	s.pulse1Enabled = (data & 0b00000001) == 0b00000001

	// TODO zero counters

	if !s.pulse1Enabled {
		apu.pulse1.lengthCounter = 0
	}
	if !s.pulse2Enabled {
		apu.pulse2.lengthCounter = 0
	}
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
	if apu.pulse2.lengthCounter > 0 {
		result |= 0b00000010
	}
	if apu.pulse1.lengthCounter > 0 {
		result |= 0b00000001
	}

	return result
}
