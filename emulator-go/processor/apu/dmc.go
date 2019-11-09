package apu

type dmcRegister struct {
	timer        uint16
	memoryReader uint16
	sampleBuffer uint16
	outputUnit   uint16
}

func (dmc *dmcRegister) outputValue() byte {
	return 0
}
