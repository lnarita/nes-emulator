package apu

type pulseRegister struct {
	dutyCounter byte

	// byte 1
	duty           byte // DDLC VVVV
	envelopeLoop   bool
	constantVolume bool
	volume         byte

	// byte 2 sweep stuff https://wiki.nesdev.com/w/index.php/APU_Sweep
	sweepUnitEnabled bool // EPPP NSSS
	sweepPeriod      byte
	sweepNegate      bool
	sweepShift       byte

	// byte 3
	timerLow byte

	// byte 4
	lengthCounter byte // LLLL LTTT
	timerHigh     byte

	timerValue uint16

	dutyCycleSequences map[byte][]byte
}

func (p *pulseRegister) init() {
	p.dutyCycleSequences = make(map[byte][]byte)
	p.dutyCycleSequences[0] = []byte{0, 1, 0, 0, 0, 0, 0, 0} // 12.5%
	p.dutyCycleSequences[1] = []byte{0, 1, 1, 0, 0, 0, 0, 0} // 25%
	p.dutyCycleSequences[2] = []byte{0, 1, 1, 1, 1, 0, 0, 0} // 50%
	p.dutyCycleSequences[3] = []byte{1, 0, 0, 1, 1, 1, 1, 1} // 25% negated
}

func (p *pulseRegister) writeByte1(data byte) {
	p.duty = data >> 6
	p.envelopeLoop = ((data << 2) >> 7) == 1
	p.constantVolume = ((data << 3) >> 7) == 1
	p.volume = data & 0b0001111
}

func (p *pulseRegister) writeByte2(data byte) {
	p.sweepUnitEnabled = data>>7 == 1
	p.sweepPeriod = (data & 0b01110000) >> 4
	p.sweepNegate = ((data << 4) >> 7) == 1
	p.sweepShift = data & 0b0000111
}

func (p *pulseRegister) writeByte3(data byte) {
	p.timerLow = data
}

func (p *pulseRegister) writeByte4(data byte) {
	p.lengthCounter = data >> 3
	p.timerLow = data & 0b0000111
}

func (p *pulseRegister) getTimerPeriod() uint16 {
	return uint16((uint16(p.timerLow)) | (uint16(p.timerHigh) << 5))
}

func (p *pulseRegister) getFrequency() int {
	period := int(p.getTimerPeriod())
	return cpuFrequency / (16 * period)
}

func (p *pulseRegister) stepTimer() {
	if p.timerValue == 0 {
		p.timerValue = p.getTimerPeriod()
		p.dutyCounter = (p.dutyCounter + 1) % 8
	} else {
		p.timerValue--
	}
}

func (p *pulseRegister) outputValue() byte {
	if p.lengthCounter == 0 { // disabled
		return 0
	}

	if p.dutyCycleSequences[p.duty][p.dutyCounter] == 0 {
		return 0
	}

	if p.getTimerPeriod() < 8 || p.getTimerPeriod() > 0x7FF { // if timer size is invalid
		return 0
	}

	//check if envelope is enabled
	return p.volume
}
