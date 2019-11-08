package apu

type pulseRegister struct {
	// byte 1
	duty           byte // DDLC VVVV
	envelopeLoop   bool
	constantVolume bool
	volume         byte

	// byte 2
	sweepUnitEnabled bool // EPPP NSSS
	period           byte
	negate           bool
	shift            byte

	// byte 3
	timerLow byte

	// byte 4
	lengthCounterLoad byte // LLLL LTTT
	timerHigh         byte
}

func (p *pulseRegister) readByte1() byte {
	var result byte = 0
	result |= p.duty << 6
	if p.envelopeLoop {
		result |= 0b00100000
	}
	if p.constantVolume {
		result |= 0b00010000
	}
	result |= p.volume

	return result
}

func (p *pulseRegister) writeByte1(data byte) {
	p.duty = data >> 6
	p.envelopeLoop = ((data << 2) >> 7) == 1
	p.constantVolume = ((data << 3) >> 7) == 1
	p.volume = data & 0b0001111
}

func (p *pulseRegister) readByte2() byte {
	var result byte = 0
	if p.sweepUnitEnabled {
		result |= 0b10000000
	}
	result |= p.period << 4
	if p.negate {
		result |= 0b00001000
	}
	result |= p.shift

	return result
}

func (p *pulseRegister) writeByte2(data byte) {
	p.sweepUnitEnabled = data>>7 == 1
	p.period = (data & 0b01110000) >> 4
	p.constantVolume = ((data << 4) >> 7) == 1
	p.volume = data & 0b0000111
}

func (p *pulseRegister) readByte3() byte {
	return p.timerLow
}

func (p *pulseRegister) writeByte3(data byte) {
	p.timerLow = data
}

func (p *pulseRegister) readByte4() byte {
	var result byte = 0
	result |= p.lengthCounterLoad << 3
	result |= p.timerLow

	return result
}

func (p *pulseRegister) writeByte4(data byte) {
	p.lengthCounterLoad = data >> 3
	p.timerLow = data & 0b0000111
}
