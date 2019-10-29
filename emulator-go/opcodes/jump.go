package opcodes

import (
	"students.ic.unicamp.br/goten/processor"
)

func branch(console *processor.Console, state *State, variation *Variation, shouldTakeBranch bool) int {
	var cycleAcc int = 0
	acc, _ := variation.addressingMode.FetchAddress(console, state)

	branchAddress := console.CPU.PC + acc
	overflow := (console.CPU.PC & 0xFF00) != (branchAddress & 0xFF00)

	if shouldTakeBranch {
		console.CPU.PC = branchAddress
		cycleAcc++
		if overflow {
			cycleAcc++
		}
	}

	return variation.cycles + cycleAcc
}

type bpl struct{}

func (o bpl) Exec(console *processor.Console, variation *Variation, state *State) int {
	return branch(console, state, variation, !console.CPU.IsNegative())
}

func (o bpl) getVariations() []Variation {
	return []Variation{
		{opcode: 0x10, addressingMode: Relative, cycles: 2},
	}
}

func (o bpl) GetName() string {
	return "BPL"
}

type bmi struct{}

func (o bmi) Exec(console *processor.Console, variation *Variation, state *State) int {
	return branch(console, state, variation, console.CPU.IsNegative())
}

func (o bmi) getVariations() []Variation {
	return []Variation{
		{opcode: 0x30, addressingMode: Relative, cycles: 2},
	}
}

func (o bmi) GetName() string {
	return "BMI"
}

type bvc struct{}

func (o bvc) Exec(console *processor.Console, variation *Variation, state *State) int {
	return branch(console, state, variation, !console.CPU.IsOverflow())
}

func (o bvc) getVariations() []Variation {
	return []Variation{
		{opcode: 0x50, addressingMode: Relative, cycles: 2},
	}
}

func (o bvc) GetName() string {
	return "BVC"
}

type bvs struct{}

func (o bvs) Exec(console *processor.Console, variation *Variation, state *State) int {
	return branch(console, state, variation, console.CPU.IsOverflow())
}

func (o bvs) getVariations() []Variation {
	return []Variation{
		{opcode: 0x70, addressingMode: Relative, cycles: 2},
	}
}

func (o bvs) GetName() string {
	return "BVS"
}

type bcc struct{}

func (o bcc) Exec(console *processor.Console, variation *Variation, state *State) int {
	return branch(console, state, variation, !console.CPU.HasCarry())
}

func (o bcc) getVariations() []Variation {
	return []Variation{
		{opcode: 0x90, addressingMode: Relative, cycles: 2},
	}
}

func (o bcc) GetName() string {
	return "BCC"
}

type bcs struct{}

func (o bcs) Exec(console *processor.Console, variation *Variation, state *State) int {
	return branch(console, state, variation, console.CPU.HasCarry())
}

func (o bcs) getVariations() []Variation {
	return []Variation{
		{opcode: 0xB0, addressingMode: Relative, cycles: 2},
	}
}

func (o bcs) GetName() string {
	return "BCS"
}

type bne struct{}

func (o bne) Exec(console *processor.Console, variation *Variation, state *State) int {
	return branch(console, state, variation, !console.CPU.IsZero())
}

func (o bne) getVariations() []Variation {
	return []Variation{
		{opcode: 0xD0, addressingMode: Relative, cycles: 2},
	}
}

func (o bne) GetName() string {
	return "BNE"
}

type beq struct{}

func (o beq) Exec(console *processor.Console, variation *Variation, state *State) int {
	return branch(console, state, variation, console.CPU.IsZero())
}

func (o beq) getVariations() []Variation {
	return []Variation{
		{opcode: 0xF0, addressingMode: Relative, cycles: 2},
	}
}

func (o beq) GetName() string {
	return "BEQ"
}

type brk struct{}

func (o brk) Exec(console *processor.Console, variation *Variation, state *State) int {
	console.Memory.StackPushAddress(console.CPU, console.CPU.PC)

	flags := console.CPU.Flags
	console.Memory.StackPushData(console.CPU, flags|processor.BreakBit)

	irq := uint16(console.Memory.FetchAddress(processor.IRQ))

	console.CPU.PC = irq
	console.CPU.SetBreak(true)

	return variation.cycles
}

func (o brk) getVariations() []Variation {
	return []Variation{
		{opcode: 0x00, addressingMode: nil, cycles: 7},
	}
}

func (o brk) GetName() string {
	return "BRK"
}

type rti struct{}

func (o rti) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.Memory.StackPopData(console.CPU)

	flags := value&processor.NotBreakBit | processor.BFlag
	console.CPU.Flags = flags
	pc := uint16(console.Memory.StackPopAddress(console.CPU))
	console.CPU.PC = pc
	return variation.cycles
}

func (o rti) getVariations() []Variation {
	return []Variation{
		{opcode: 0x40, addressingMode: nil, cycles: 6},
	}
}

func (o rti) GetName() string {
	return "RTI"
}

type jsr struct{}

func (o jsr) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)

	console.Memory.StackPushAddress(console.CPU, console.CPU.PC-1)

	console.CPU.PC = address

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o jsr) getVariations() []Variation {
	return []Variation{
		{opcode: 0x20, addressingMode: Absolute, cycles: 6},
	}
}

func (o jsr) GetName() string {
	return "JSR"
}

type rts struct{}

func (o rts) Exec(console *processor.Console, variation *Variation, state *State) int {
	address := console.Memory.StackPopAddress(console.CPU)
	console.CPU.PC = uint16(address) + 1
	return variation.cycles
}

func (o rts) getVariations() []Variation {
	return []Variation{
		{opcode: 0x60, addressingMode: nil, cycles: 6},
	}
}

func (o rts) GetName() string {
	return "RTS"
}

type jmp struct{}


func (o jmp) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)

	console.CPU.PC = address

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o jmp) getVariations() []Variation {
	return []Variation{
		{opcode: 0x4C, addressingMode: Absolute, cycles: 3},
		{opcode: 0x6C, addressingMode: Indirect, cycles: 5},
	}
}

func (o jmp) GetName() string {
	return "JMP"
}

var JumpOpCodes = []OpCode{
	bpl{},
	bmi{},
	bvc{},
	bvs{},
	bcc{},
	bcs{},
	bne{},
	beq{},
	brk{},
	rti{},
	jsr{},
	rts{},
	jmp{},
}
