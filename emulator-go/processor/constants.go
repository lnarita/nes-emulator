package processor

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

	NEGATIVE_BIT       byte = 0b1000_0000
	OVERFLOW_BIT       byte = 0b0100_0000
	B_FLAG             byte = 0b0010_0000
	BREAK_BIT          byte = 0b0001_0000
	DECIMAL_BIT        byte = 0b0000_1000
	INTERRUPTS_BIT     byte = 0b0000_0100
	ZERO_BIT           byte = 0b0000_0010
	CARRY_BIT          byte = 0b0000_0001
	NOT_NEGATIVE_BIT        = ^NEGATIVE_BIT
	NOT_OVERFLOW_BIT        = ^OVERFLOW_BIT
	NOT_B_FLAG              = ^B_FLAG
	NOT_BREAK_BIT           = ^BREAK_BIT
	NOT_DECIMAL_BIT         = ^DECIMAL_BIT
	NOT_INTERRUPTS_BIT      = ^INTERRUPTS_BIT
	NOT_ZERO_BIT            = ^ZERO_BIT
	NOT_CARRY_BIT           = ^CARRY_BIT
)
