package apu

import "log"

const cpuFrequency = 1.79e6 // to avoid cyclic imports
const frameCounterRate = cpuFrequency / 240.0

type APU struct {
	// APU Registers per channel
	pulse1       pulseRegister        // $4000-$4003
	pulse2       pulseRegister        // $4004-$4007
	triangle     triangleRegister     // $4008-$400B
	noise        noiseRegister        // $400C-$400F
	dmc          dmcRegister          // $4010-$4013
	status       statusRegister       // $4015
	frameCounter frameCounterRegister // $4017

	cycle uint64

	mixer

	Channel    chan float64
	SampleRate float64
}

func (apu *APU) Init() {
	apu.mixer.init()
	apu.pulse1.init()
	apu.pulse2.init()
	apu.Channel = make(chan float64, 44100)
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
		apu.status.write(data, apu)
	case 0x4017:
		apu.frameCounter.write(data)
	}

}

func (apu *APU) Step() {
	apu.cycle++
	apu.pulse1.stepTimer()
	apu.pulse2.stepTimer()

	s1 := int(float64(apu.cycle-1) / apu.SampleRate)
	s2 := int(float64(apu.cycle) / apu.SampleRate)
	if s1 != s2 {
		log.Println("send stuff")

		select {
		case apu.Channel <- apu.mixer.output(apu):
		default:
		}
	}
}