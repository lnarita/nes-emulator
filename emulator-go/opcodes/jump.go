package opcodes

import "students.ic.unicamp.br/goten/processor"

func branch(console *processor.Console, variation *Variation, shouldTakeBranch bool) int {
	var cycleAcc int = 0
	acc, stall := variation.addressingMode.FetchAddress(console)

	branchAddress := int(int16(console.CPU.PC) + int16(int8(acc)))
	overflow := (console.CPU.PC & processor.HighBitsMask) != (branchAddress & processor.HighBitsMask)

	if shouldTakeBranch {
		console.CPU.PC = branchAddress
		cycleAcc++
		if overflow {
			cycleAcc++
		}
	}

	if stall {
		cycleAcc++
	}
	return variation.cycles + cycleAcc
}

type blp struct{}

func (o blp) exec(console *processor.Console, variation *Variation) int {
	return branch(console, variation, !console.CPU.IsNegative())
}

func (o blp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x10, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o blp) getName() string {
	return "BLP"
}

type bmi struct{}

func (o bmi) exec(console *processor.Console, variation *Variation) int {
	return branch(console, variation, console.CPU.IsNegative())
}

func (o bmi) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x30, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bmi) getName() string {
	return "BMI"
}

type bvc struct{}

func (o bvc) exec(console *processor.Console, variation *Variation) int {
	return branch(console, variation, !console.CPU.IsOverflow())
}

func (o bvc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x50, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bvc) getName() string {
	return "BVC"
}

type bvs struct{}

func (o bvs) exec(console *processor.Console, variation *Variation) int {
	return branch(console, variation, console.CPU.IsOverflow())
}

func (o bvs) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x70, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bvs) getName() string {
	return "BVS"
}

type bcc struct{}

func (o bcc) exec(console *processor.Console, variation *Variation) int {
	return branch(console, variation, !console.CPU.HasCarry())

}

func (o bcc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x90, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bcc) getName() string {
	return "bcc"
}

type bcs struct{}

func (o bcs) exec(console *processor.Console, variation *Variation) int {
	return branch(console, variation, console.CPU.HasCarry())
}

func (o bcs) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xB0, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bcs) getName() string {
	return "bcs"
}

type bne struct{}

func (o bne) exec(console *processor.Console, variation *Variation) int {
	return branch(console, variation, !console.CPU.IsZero())
}

func (o bne) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xD0, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bne) getName() string {
	return "BNE"
}

type beq struct{}

func (o beq) exec(console *processor.Console, variation *Variation) int {
	return branch(console, variation, console.CPU.IsZero())
}

func (o beq) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xF0, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o beq) getName() string {
	return "BEQ"
}

type brk struct{}

func (o brk) exec(console *processor.Console, variation *Variation) int {
	console.Memory.StackPushAddress(console.CPU, console.CPU.PC)

	flags := console.CPU.Flags
	console.Memory.StackPushData(console.CPU, flags|processor.BreakBit)

	irq := console.Memory.FetchAddress(processor.IRQ)

	console.CPU.PC = irq
	console.CPU.SetBreak(true)

	return variation.cycles
}

func (o brk) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x00, addressingMode: nil, cycles: 7},
	}
}

func (o brk) getName() string {
	return "BRK"
}

type rti struct{}

func (o rti) exec(console *processor.Console, variation *Variation) int {
	value := console.Memory.StackPopData(console.CPU)

	flags := value&processor.NotBreakBit | processor.BFlag
	console.CPU.Flags = flags
	pc := console.Memory.StackPopAddress(console.CPU)
	console.CPU.PC = pc
	return variation.cycles
}

func (o rti) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x00, addressingMode: nil, cycles: 6},
	}
}

func (o rti) getName() string {
	return "RTI"
}

type jsr struct{}

func (o jsr) exec(console *processor.Console, variation *Variation) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console)

	console.Memory.StackPushAddress(console.CPU, console.CPU.PC-1)

	console.CPU.PC = address

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o jsr) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x20, addressingMode: processor.Absolute, cycles: 6},
	}
}

func (o jsr) getName() string {
	return "JSR"
}

type rts struct{}

func (o rts) exec(console *processor.Console, variation *Variation) int {
	address := console.Memory.StackPopAddress(console.CPU)
	console.CPU.PC = address + 1
	return variation.cycles
}

func (o rts) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x60, addressingMode: nil, cycles: 6},
	}
}

func (o rts) getName() string {
	return "JSR"
}

type jmp struct{}

func (o jmp) exec(console *processor.Console, variation *Variation) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console)

	console.CPU.PC = address

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o jmp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x4C, addressingMode: processor.Absolute, cycles: 3},
		Variation{opcode: 0x6C, addressingMode: processor.Indirect, cycles: 5},
	}
}

func (o jmp) getName() string {
	return "JMP"
}

var JumpOpCodes = []OpCode{
	blp{},
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
