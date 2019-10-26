package opcodes

import "students.ic.unicamp.br/goten/processor"

type bit struct{}

func (o bit) exec(console *processor.Console, variation *Variation) int {
	address, stall := variation.addressingMode.FetchAddress(console)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	console.CPU.SetZN(value)
	console.CPU.SetOverflow(value&processor.OverflowBit != 0)

	var cycleAcc int = 0
	if stall {
		cycleAcc++
	}
	return variation.cycles + cycleAcc
}

func (o bit) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x24, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x2C, addressingMode: processor.Absolute, cycles: 4},
	}
}

func (o bit) getName() string {
	return "BIT"
}

type clc struct{}

func (o clc) exec(console *processor.Console, variation *Variation) int {
	console.CPU.SetCarry(false)
	return variation.cycles
}

func (o clc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x18, addressingMode: nil, cycles: 2},
	}
}

func (o clc) getName() string {
	return "CLC"
}

type sec struct{}

func (o sec) exec(console *processor.Console, variation *Variation) int {
	console.CPU.SetCarry(true)
	return variation.cycles
}

func (o sec) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x38, addressingMode: nil, cycles: 2},
	}
}

func (o sec) getName() string {
	return "SEC"
}

type cld struct{}

func (o cld) exec(console *processor.Console, variation *Variation) int {
	console.CPU.SetDecimalMode(false)
	return variation.cycles
}

func (o cld) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xD8, addressingMode: nil, cycles: 2},
	}
}

func (o cld) getName() string {
	return "CLD"
}

type sed struct{}

func (o sed) exec(console *processor.Console, variation *Variation) int {
	console.CPU.SetDecimalMode(true)
	return variation.cycles
}

func (o sed) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xF8, addressingMode: nil, cycles: 2},
	}
}

func (o sed) getName() string {
	return "SED"
}

type cli struct{}

func (o cli) exec(console *processor.Console, variation *Variation) int {
	console.CPU.DisableInterrupts(false)
	return variation.cycles
}

func (o cli) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x58, addressingMode: nil, cycles: 2},
	}
}

func (o cli) getName() string {
	return "CLI"
}

type sei struct{}

func (o sei) exec(console *processor.Console, variation *Variation) int {
	console.CPU.DisableInterrupts(true)
	return variation.cycles
}

func (o sei) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x78, addressingMode: nil, cycles: 2},
	}
}

func (o sei) getName() string {
	return "SEI"
}

type clv struct{}

func (o clv) exec(console *processor.Console, variation *Variation) int {
	console.CPU.SetOverflow(false)
	return variation.cycles
}

func (o clv) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xB8, addressingMode: nil, cycles: 2},
	}
}

func (o clv) getName() string {
	return "CLV"
}

type nop struct{}

func (o nop) exec(console *processor.Console, variation *Variation) int {
	return variation.cycles
}

func (o nop) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xEA, addressingMode: nil, cycles: 2},
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
