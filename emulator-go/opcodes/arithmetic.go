package opcodes

import "students.ic.unicamp.br/goten/processor"

type ora struct{}

func (o ora) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o ora) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x01, addresingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x05, addresingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x09, addresingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x0D, addresingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x11, addresingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0x15, addresingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x19, addresingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0x1D, addresingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o ora) getName() string {
	return "ORA"
}

type and struct{}

func (o and) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o and) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x21, addresingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x25, addresingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x29, addresingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x2D, addresingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x31, addresingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0x35, addresingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x39, addresingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0x3D, addresingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o and) getName() string {
	return "AND"
}

type eor struct{}

func (o eor) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o eor) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x41, addresingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x45, addresingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x49, addresingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x4D, addresingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x51, addresingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0x55, addresingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x59, addresingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0x5D, addresingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o eor) getName() string {
	return "EOR"
}

type adc struct{}

func (o adc) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o adc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x61, addresingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x65, addresingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x69, addresingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x6D, addresingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x71, addresingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0x75, addresingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x79, addresingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0x7D, addresingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o adc) getName() string {
	return "ADC"
}

type sbc struct{}

func (o sbc) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o sbc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xE1, addresingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0xE5, addresingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xE9, addresingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xED, addresingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0xF1, addresingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0xF5, addresingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0xF9, addresingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0xFD, addresingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o sbc) getName() string {
	return "SBC"
}

type cmp struct{}

func (o cmp) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o cmp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xC1, addresingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0xC5, addresingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xC9, addresingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xCD, addresingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0xD1, addresingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0xD5, addresingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0xD9, addresingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0xDD, addresingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o cmp) getName() string {
	return "CMP"
}

type cpx struct{}

func (o cpx) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o cpx) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xE0, addresingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xE4, addresingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xEC, addresingMode: processor.Absolute, cycles: 4},
	}
}

func (o cpx) getName() string {
	return "CPX"
}

type cpy struct{}

func (o cpy) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o cpy) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xC0, addresingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xC4, addresingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xCC, addresingMode: processor.Absolute, cycles: 4},
	}
}

func (o cpy) getName() string {
	return "CPY"
}

type dec struct{}

func (o dec) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o dec) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xC6, addresingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0xCE, addresingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0xD6, addresingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0xDE, addresingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o dec) getName() string {
	return "DEC"
}

type dex struct{}

func (o dex) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o dex) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xCA, addresingMode: nil, cycles: 2},
	}
}

func (o dex) getName() string {
	return "DEX"
}

type dey struct{}

func (o dey) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o dey) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x88, addresingMode: nil, cycles: 2},
	}
}

func (o dey) getName() string {
	return "DEY"
}

type inc struct{}

func (o inc) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o inc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xE6, addresingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0xEE, addresingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0xF6, addresingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0xFE, addresingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o inc) getName() string {
	return "INC"
}

type inx struct{}

func (o inx) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o inx) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xE8, addresingMode: nil, cycles: 2},
	}
}

func (o inx) getName() string {
	return "INX"
}

type iny struct{}

func (o iny) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o iny) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xC8, addresingMode: nil, cycles: 2},
	}
}

func (o iny) getName() string {
	return "INY"
}

type asl struct{}

func (o asl) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o asl) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x06, addresingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x0A, addresingMode: processor.ZeroPage, cycles: 2},
		Variation{opcode: 0x0E, addresingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x16, addresingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x1E, addresingMode: processor.ZeroPage, cycles: 7},
	}
}

func (o asl) getName() string {
	return "ASL"
}

type rol struct{}

func (o rol) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o rol) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x26, addresingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x2A, addresingMode: processor.ZeroPage, cycles: 2},
		Variation{opcode: 0x2E, addresingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x36, addresingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x3E, addresingMode: processor.ZeroPage, cycles: 7},
	}
}

func (o rol) getName() string {
	return "ROL"
}

type lsr struct{}

func (o lsr) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o lsr) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x46, addresingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x4A, addresingMode: processor.ZeroPage, cycles: 2},
		Variation{opcode: 0x4E, addresingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x56, addresingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x5E, addresingMode: processor.ZeroPage, cycles: 7},
	}
}

func (o lsr) getName() string {
	return "LSR"
}

type ror struct{}

func (o ror) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o ror) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x66, addresingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x6A, addresingMode: processor.ZeroPage, cycles: 2},
		Variation{opcode: 0x6E, addresingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x76, addresingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x7E, addresingMode: processor.ZeroPage, cycles: 7},
	}
}

func (o ror) getName() string {
	return "ROR"
}

var ArithmeticAndLogicalOpCodes = []OpCode{
	ora{},
	and{},
	eor{},
	adc{},
	sbc{},
	cmp{},
	cpx{},
	cpy{},
	dec{},
	dex{},
	dey{},
	inc{},
	inx{},
	iny{},
	asl{},
	rol{},
	lsr{},
	ror{},
}
