package opcodes

import "students.ic.unicamp.br/goten/processor"

type bit struct{}

func (o bit) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o bit) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x24, addresingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x2C, addresingMode: processor.Absolute, cycles: 4},
	}
}

func (o bit) getName() string {
	return "BIT"
}

type clc struct{}

func (o clc) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o clc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x18, addresingMode: nil, cycles: 2},
	}
}

func (o clc) getName() string {
	return "CLC"
}

type sec struct{}

func (o sec) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o sec) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x38, addresingMode: nil, cycles: 2},
	}
}

func (o sec) getName() string {
	return "SEC"
}

type cld struct{}

func (o cld) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o cld) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xD8, addresingMode: nil, cycles: 2},
	}
}

func (o cld) getName() string {
	return "CLD"
}

type sed struct{}

func (o sed) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o sed) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xF8, addresingMode: nil, cycles: 2},
	}
}

func (o sed) getName() string {
	return "SED"
}

type cli struct{}

func (o cli) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o cli) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x58, addresingMode: nil, cycles: 2},
	}
}

func (o cli) getName() string {
	return "CLI"
}

type sei struct{}

func (o sei) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o sei) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x78, addresingMode: nil, cycles: 2},
	}
}

func (o sei) getName() string {
	return "SEI"
}

type clv struct{}

func (o clv) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o clv) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xB8, addresingMode: nil, cycles: 2},
	}
}

func (o clv) getName() string {
	return "CLV"
}

type nop struct{}

func (o nop) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o nop) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xEA, addresingMode: nil, cycles: 2},
	}
}

func (o nop) getName() string {
	return "NOP"
}

// FlagOpcodes flag manipulating opcodes
var FlagOpcodes = []OpCode{
	bit{},
	clc{},
	sec{},
	cld{},
	sed{},
	cli{},
	sei{},
	clv{},
	nop{},
}
