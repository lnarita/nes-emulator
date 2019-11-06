package processor

type APU struct {
	// APU Registers per channel
	pulse1        pulseRegisters    // $4000-$4003
	pulse2        pulseRegisters    // $4004-$4007
	triangle      triangleRegisters // $4008-$400B
	noise         noiseRegisters    // $400C-$400F
	dmc           dmcRegisters      // $4010-$4013
	channelEnable uint16            // $4015
	frameCounter  uint16            // $4017

}

type pulseRegisters struct {
	// byte 1
	duty           uint8 // DDLC VVVV
	envelopeLoop   bool
	constantVolume bool
	volume         uint8

	// byte 2
	sweepUnitEnabled bool // EPPP NSSS
	period           uint8
	negate           bool
	shift            uint8

	// byte 3
	timerLow uint8

	// byte 4
	lengthCounterLoad uint8 // LLLL LTTT
	timerHigh         uint8
}
type triangleRegisters struct {
	timer         uint16
	lengthCounter uint16
	linearCounter uint16
}

type noiseRegisters struct {
	timer                       uint16
	lengthCounter               uint16
	envelope                    uint16
	linearFeedbackShiftRegister uint16
}

type dmcRegisters struct {
	timer        uint16
	memoryReader uint16
	sampleBuffer uint16
	outputUnit   uint16
}

func (apu *APU) Read(address uint16) byte {
	switch address {
	case 0x4000:
	case 0x4001:
	case 0x4002:
	case 0x4003:
	case 0x4004:
	case 0x4005:
	case 0x4006:
	case 0x4007:
	case 0x4008:
	case 0x4009:
	case 0x400A:
	case 0x400B:
	case 0x400C:
	case 0x400D:
	case 0x400E:
	case 0x400F:
	case 0x4010:
	case 0x4012:
	case 0x4013:
	case 0x4015:
	case 0x4017:
	}

	return 0
}

func (apu *APU) Write(address uint16, data byte) {

}
