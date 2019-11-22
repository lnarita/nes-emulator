package apu

var lengthTable = []byte{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

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

	initEnvelope   bool
	envelopeVolume byte
	envelopeValue  byte

	initSweep  bool
	sweepValue byte

	timerPeriod uint16
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
	p.initEnvelope = true
}

func (p *pulseRegister) writeByte2(data byte) {
	p.sweepUnitEnabled = data>>7 == 1
	p.sweepPeriod = (data & 0b01110000) >> 4
	p.sweepNegate = ((data << 4) >> 7) == 1
	p.sweepShift = data & 0b0000111
	p.initSweep = true
}

func (p *pulseRegister) writeByte3(data byte) {
	p.timerLow = data
	p.timerPeriod = (p.timerPeriod & 0xFF00) | uint16(data)
}

func (p *pulseRegister) writeByte4(data byte) {
	p.lengthCounter = lengthTable[data>>3]
	p.timerHigh = data & 0b0000111
	p.timerPeriod = (p.timerPeriod & 0x00FF) | (uint16(data&7) << 8)
	p.dutyCounter = 0
}

func (p *pulseRegister) stepTimer() {
	if p.timerValue == 0 {
		p.timerValue = p.timerPeriod
		p.dutyCounter = (p.dutyCounter + 1) % 8
	} else {
		p.timerValue--
	}
}

func (p *pulseRegister) stepEnvelope() {
	if p.initEnvelope {
		p.envelopeVolume = 15
		p.envelopeValue = p.volume
		p.initEnvelope = false
	} else if p.envelopeValue > 0 {
		p.envelopeValue--
	} else {
		if p.envelopeVolume > 0 {
			p.envelopeVolume--
		} else if p.envelopeLoop {
			p.envelopeVolume = 15
		}
		p.envelopeValue = p.volume
	}
}

func (p *pulseRegister) stepSweep(pulse1 bool) {
	if p.initSweep {
		if p.sweepUnitEnabled && p.sweepValue == 0 {
			p.doSweep(pulse1)
		}
		p.sweepValue = p.sweepPeriod
		p.initSweep = false
	} else if p.sweepValue > 0 {
		p.sweepValue--
	} else {
		if p.sweepUnitEnabled {
			p.doSweep(pulse1)
		}
		p.sweepValue = p.sweepPeriod
	}
}

func (p *pulseRegister) doSweep(pulse1 bool) {
	barrel := p.timerPeriod >> p.sweepPeriod
	if p.sweepNegate {
		p.timerPeriod -= barrel
		if pulse1 {
			p.timerPeriod--
		}
	} else {
		p.timerPeriod += barrel
	}
}

func (p *pulseRegister) outputValue() byte {
	if p.lengthCounter == 0 { // disabled
		return 0
	}

	if p.dutyCycleSequences[p.duty][p.dutyCounter] == 0 {
		return 0
	}

	if p.timerPeriod < 8 || p.timerPeriod > 0x7FF { // if timer size is invalid
		return 0
	}

	if p.constantVolume {
		return p.volume
	} else {
		return p.envelopeVolume
	}
}
