package apu

type noiseRegister struct {
	timer                       uint16
	lengthCounter               uint16
	envelope                    uint16
	linearFeedbackShiftRegister uint16
}

func (n *noiseRegister) outputValue() byte {
	return 0
}
