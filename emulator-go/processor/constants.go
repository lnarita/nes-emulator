package processor

const (
	// CPUFrequency frequency of the nes cpu
	CPUFrequency = 1.79e6

	// CyclePeriod period of the nes cpu
	CyclePeriod = 1. / CPUFrequency

	// KB bits in a kb
	KB = 1024

	// HighBitsMask binary mask for the high 8 bits
	HighBitsMask = 0xFF00

	// LowBitsMask binary mask for the low 8 bits
	LowBitsMask = 0x00FF

	NegativeBit      byte = 0b1000_0000
	OverflowBit      byte = 0b0100_0000
	BFlag            byte = 0b0010_0000
	BreakBit         byte = 0b0001_0000
	DecimalBit       byte = 0b0000_1000
	InterruptsBit    byte = 0b0000_0100
	ZeroBit          byte = 0b0000_0010
	CarryBit         byte = 0b0000_0001
	NotNegativeBit        = ^NegativeBit
	NotOverflowBit        = ^OverflowBit
	NotBFlag              = ^BFlag
	NotBreakBit           = ^BreakBit
	NotDecimalBit         = ^DecimalBit
	NotInterruptsBit      = ^InterruptsBit
	NotZeroBit            = ^ZeroBit
	NotCarryBit           = ^CarryBit
)
