package processor

// CPUFrequency frequency of the nes cpu
const CPUFrequency = 1.79e6

// CyclePeriod period of the nes cpu
const CyclePeriod = 1. / CPUFrequency

// KB bits in a kb
const KB = 1024

// HighBitsMask binary mask for the high 8 bits
const HighBitsMask = 0b1111111100000000

// LowBitsMask binary mask for the low 8 bits
const LowBitsMask = 0b0000000011111111

// NegativeBit binary mask for the negative bit
const NegativeBit = 0b0000000010000000

const NEGATIVE_BIT byte = 0b1000_0000
const OVERFLOW_BIT byte = 0b0100_0000
const B_FLAG byte = 0b0010_0000
const BREAK_BIT byte = 0b0001_0000
const DECIMAL_BIT byte = 0b0000_1000
const INTERRUPTS_BIT byte = 0b0000_0100
const ZERO_BIT byte = 0b0000_0010
const CARRY_BIT byte = 0b0000_0001
const NOT_NEGATIVE_BIT = ^NEGATIVE_BIT
const NOT_OVERFLOW_BIT = ^OVERFLOW_BIT
const NOT_B_FLAG = ^B_FLAG
const NOT_BREAK_BIT = ^BREAK_BIT
const NOT_DECIMAL_BIT = ^DECIMAL_BIT
const NOT_INTERRUPTS_BIT = ^INTERRUPTS_BIT
const NOT_ZERO_BIT = ^ZERO_BIT
const NOT_CARRY_BIT = ^CARRY_BIT
