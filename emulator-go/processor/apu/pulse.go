package apu

type pulseRegister struct {
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
	lengthCounterLoad byte // LLLL LTTT
	timerHigh         byte

	dutyCycleSequences map[int]byte
}

func (p *pulseRegister) init() {
	p.dutyCycleSequences = make(map[int]byte)
	p.dutyCycleSequences[0] = 0b01000000 // 12.5%
	p.dutyCycleSequences[1] = 0b01100000 // 25%
	p.dutyCycleSequences[2] = 0b01111000 // 50%
	p.dutyCycleSequences[3] = 0b10011111 // 25% negated
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
	p.lengthCounterLoad = data >> 3
	p.timerLow = data & 0b0000111
}

const cpuFrequency = 1.79e6 // to avoid cyclic imports

func (p *pulseRegister) getFrequency() int {
	period := int((uint16(p.timerLow)) | (uint16(p.timerHigh) << 5))
	return cpuFrequency / (16 * period)
}

func (p *pulseRegister) outputVolume() byte {
	//check if envelope is enabled
	return p.volume
}
