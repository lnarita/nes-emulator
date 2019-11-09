package apu

type APU struct {
	// APU Registers per channel
	pulse1       pulseRegister        // $4000-$4003
	pulse2       pulseRegister        // $4004-$4007
	triangle     triangleRegister     // $4008-$400B
	noise        noiseRegister        // $400C-$400F
	dmc          dmcRegister          // $4010-$4013
	status       statusRegister       // $4015
	frameCounter frameCounterRegister // $4017

}

type triangleRegister struct {
	timer         uint16
	lengthCounter uint16
	linearCounter uint16
}

type noiseRegister struct {
	timer                       uint16
	lengthCounter               uint16
	envelope                    uint16
	linearFeedbackShiftRegister uint16
}

type dmcRegister struct {
	timer        uint16
	memoryReader uint16
	sampleBuffer uint16
	outputUnit   uint16
}

func (apu *APU) Read(address uint16) byte {
	switch address {
	case 0x4015:
		return apu.status.read(apu)
	}

	return 0
}

func (apu *APU) Write(address uint16, data byte) {
	switch address {
	case 0x4000:
		apu.pulse1.writeByte1(data)
	case 0x4001:
		apu.pulse1.writeByte2(data)
	case 0x4002:
		apu.pulse1.writeByte3(data)
	case 0x4003:
		apu.pulse1.writeByte4(data)
	case 0x4004:
		apu.pulse2.writeByte1(data)
	case 0x4005:
		apu.pulse2.writeByte2(data)
	case 0x4006:
		apu.pulse2.writeByte3(data)
	case 0x4007:
		apu.pulse2.writeByte4(data)
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
		apu.status.write(data)
	case 0x4017:
		apu.frameCounter.write(data)
	}

}
