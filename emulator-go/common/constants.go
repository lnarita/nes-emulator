package common

const (
	// CPUFrequency frequency of the nes cpu
	CPUFrequency = 1.79e6

	// CyclePeriod period of the nes cpu
	CyclePeriod = 1. / CPUFrequency

	// KB bits in a kb
	KB = 1024

	// HighBitsMask binary mask for the high 8 bits
	HighBitsMask = 0b1111111100000000

	// LowBitsMask binary mask for the low 8 bits
	LowBitsMask = 0b0000000011111111

	// NegativeBit binary mask for the negative bit
	NegativeBit = 0b0000000010000000
)
