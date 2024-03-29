package opcodes

import "students.ic.unicamp.br/goten/processor"

type bit struct{}

func (o bit) Exec(console *processor.Console, variation *Variation, state *State) int {
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	console.CPU.SetZN(value)
	console.CPU.SetZero(value & console.CPU.A == 0)
	console.CPU.SetOverflow(value&processor.OverflowBit != 0)

	var cycleAcc int = 0
	if stall {
		cycleAcc++
	}
	return variation.cycles + cycleAcc
}

func (o bit) getVariations() []Variation {
	return []Variation{
		{opcode: 0x24, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x2C, addressingMode: Absolute, cycles: 4},
	}
}

func (o bit) GetName() string {
	return "BIT"
}

type clc struct{}

func (o clc) Exec(console *processor.Console, variation *Variation, state *State) int {
	console.CPU.SetCarry(false)
	return variation.cycles
}

func (o clc) getVariations() []Variation {
	return []Variation{
		{opcode: 0x18, addressingMode: nil, cycles: 2},
	}
}

func (o clc) GetName() string {
	return "CLC"
}

type sec struct{}

func (o sec) Exec(console *processor.Console, variation *Variation, state *State) int {
	console.CPU.SetCarry(true)
	return variation.cycles
}

func (o sec) getVariations() []Variation {
	return []Variation{
		{opcode: 0x38, addressingMode: nil, cycles: 2},
	}
}

func (o sec) GetName() string {
	return "SEC"
}

type cld struct{}

func (o cld) Exec(console *processor.Console, variation *Variation, state *State) int {
	console.CPU.SetDecimalMode(false)
	return variation.cycles
}

func (o cld) getVariations() []Variation {
	return []Variation{
		{opcode: 0xD8, addressingMode: nil, cycles: 2},
	}
}

func (o cld) GetName() string {
	return "CLD"
}

type sed struct{}

func (o sed) Exec(console *processor.Console, variation *Variation, state *State) int {
	console.CPU.SetDecimalMode(true)
	return variation.cycles
}

func (o sed) getVariations() []Variation {
	return []Variation{
		{opcode: 0xF8, addressingMode: nil, cycles: 2},
	}
}

func (o sed) GetName() string {
	return "SED"
}

type cli struct{}

func (o cli) Exec(console *processor.Console, variation *Variation, state *State) int {
	console.CPU.DisableInterrupts(false)
	return variation.cycles
}

func (o cli) getVariations() []Variation {
	return []Variation{
		{opcode: 0x58, addressingMode: nil, cycles: 2},
	}
}

func (o cli) GetName() string {
	return "CLI"
}

type sei struct{}

func (o sei) Exec(console *processor.Console, variation *Variation, state *State) int {
	console.CPU.DisableInterrupts(true)
	return variation.cycles
}

func (o sei) getVariations() []Variation {
	return []Variation{
		{opcode: 0x78, addressingMode: nil, cycles: 2},
	}
}

func (o sei) GetName() string {
	return "SEI"
}

type clv struct{}

func (o clv) Exec(console *processor.Console, variation *Variation, state *State) int {
	console.CPU.SetOverflow(false)
	return variation.cycles
}

func (o clv) getVariations() []Variation {
	return []Variation{
		{opcode: 0xB8, addressingMode: nil, cycles: 2},
	}
}

func (o clv) GetName() string {
	return "CLV"
}

type nop struct{}

func (o nop) Exec(console *processor.Console, variation *Variation, state *State) int {
	return variation.cycles
}

func (o nop) getVariations() []Variation {
	return []Variation{
		{opcode: 0xEA, addressingMode: nil, cycles: 2},
	}
}

func (o nop) GetName() string {
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
