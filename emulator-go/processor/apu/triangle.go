package apu

type triangleRegister struct {
	timer         uint16
	lengthCounter uint16
	linearCounter uint16
}

func (t *triangleRegister) outputValue() byte {
	return 0
}
