package opcodes

import (
	"students.ic.unicamp.br/goten/processor"
)

type lda struct{}

func (o lda) Exec(console *processor.Console, variation *Variation, state *State) int {
	var stall bool = false
	var cycleAcc int = 0
	if variation.addressingMode != nil {
		var address uint16
		address, stall = variation.addressingMode.FetchAddress(console, state)
		value := byte(variation.addressingMode.ReadFrom(console, address))

		/// log
		state.HasData = true
		state.Data = value

		console.CPU.A = value
		console.CPU.SetZN(value)
	}

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o lda) getVariations() []Variation {
	return []Variation{
		{opcode: 0xA1, addressingMode: IndirectX, cycles: 6},
		{opcode: 0xA5, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0xA9, addressingMode: Immediate, cycles: 2},
		{opcode: 0xAD, addressingMode: Absolute, cycles: 4},
		{opcode: 0xB1, addressingMode: IndirectY, cycles: 5},
		{opcode: 0xB5, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0xB9, addressingMode: AbsoluteY, cycles: 4},
		{opcode: 0xBd, addressingMode: AbsoluteX, cycles: 4},
	}
}

func (o lda) GetName() string {
	return "LDA"
}

type sta struct{}

func (o sta) Exec(console *processor.Console, variation *Variation, state *State) int {
	var stall bool = false
	var cycleAcc int = 0
	if variation.addressingMode != nil {
		var address uint16
		address, stall = variation.addressingMode.FetchAddress(console, state)

		/// log
		state.HasData = true
		state.Data = console.FetchData(address)

		variation.addressingMode.WriteTo(console, address, console.CPU.A)
	}

	if stall {
		cycleAcc++
	}

	// uhh, write instructions don't stall???
	return variation.cycles
}

func (o sta) getVariations() []Variation {
	return []Variation{
		{opcode: 0x81, addressingMode: IndirectX, cycles: 6},
		{opcode: 0x85, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x8D, addressingMode: Absolute, cycles: 4},
		{opcode: 0x91, addressingMode: IndirectY, cycles: 6},
		{opcode: 0x95, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0x99, addressingMode: AbsoluteY, cycles: 5},
		{opcode: 0x9D, addressingMode: AbsoluteX, cycles: 5},
	}
}

func (o sta) GetName() string {
	return "STA"
}

type ldx struct{}

func (o ldx) Exec(console *processor.Console, variation *Variation, state *State) int {
	var stall bool = false
	var cycleAcc int = 0
	if variation.addressingMode != nil {
		var address uint16
		address, stall = variation.addressingMode.FetchAddress(console, state)
		value := byte(variation.addressingMode.ReadFrom(console, address))

		/// log
		state.HasData = true
		state.Data = value

		console.CPU.X = value
		console.CPU.SetZN(value)
	}

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o ldx) getVariations() []Variation {
	return []Variation{
		{opcode: 0xA2, addressingMode: Immediate, cycles: 2},
		{opcode: 0xA6, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0xAE, addressingMode: Absolute, cycles: 4},
		{opcode: 0xB6, addressingMode: ZeroPageY, cycles: 4},
		{opcode: 0xBE, addressingMode: AbsoluteY, cycles: 4},
	}
}

func (o ldx) GetName() string {
	return "LDX"
}

type stx struct{}

func (o stx) Exec(console *processor.Console, variation *Variation, state *State) int {
	var stall bool = false
	var cycleAcc int = 0
	if variation.addressingMode != nil {
		var address uint16
		address, stall = variation.addressingMode.FetchAddress(console, state)

		/// log
		state.HasData = true
		state.Data = console.FetchData(address)

		variation.addressingMode.WriteTo(console, address, console.CPU.X)
	}

	if stall {
		cycleAcc++
	}

	// uhh, write instructions don't stall???
	return variation.cycles
}

func (o stx) getVariations() []Variation {
	return []Variation{
		{opcode: 0x86, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x8E, addressingMode: Absolute, cycles: 4},
		{opcode: 0x96, addressingMode: ZeroPageY, cycles: 4},
	}
}

func (o stx) GetName() string {
	return "STX"
}

type ldy struct{}

func (o ldy) Exec(console *processor.Console, variation *Variation, state *State) int {
	var stall bool = false
	var cycleAcc int = 0
	if variation.addressingMode != nil {
		var address uint16
		address, stall = variation.addressingMode.FetchAddress(console, state)
		value := byte(variation.addressingMode.ReadFrom(console, address))

		/// log
		state.HasData = true
		state.Data = value

		console.CPU.Y = value
		console.CPU.SetZN(value)
	}

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o ldy) getVariations() []Variation {
	return []Variation{
		{opcode: 0xA0, addressingMode: Immediate, cycles: 2},
		{opcode: 0xA4, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0xAC, addressingMode: Absolute, cycles: 4},
		{opcode: 0xB4, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0xBC, addressingMode: AbsoluteX, cycles: 4},
	}
}

func (o ldy) GetName() string {
	return "LDY"
}

type sty struct{}

func (o sty) Exec(console *processor.Console, variation *Variation, state *State) int {
	var stall bool = false
	var cycleAcc int = 0
	if variation.addressingMode != nil {
		var address uint16
		address, stall = variation.addressingMode.FetchAddress(console, state)

		/// log
		state.HasData = true
		state.Data = console.FetchData(address)

		variation.addressingMode.WriteTo(console, address, console.CPU.Y)
	}

	if stall {
		cycleAcc++
	}

	// uhh, write instructions don't stall???
	return variation.cycles
}

func (o sty) getVariations() []Variation {
	return []Variation{
		{opcode: 0x84, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x8C, addressingMode: Absolute, cycles: 4},
		{opcode: 0x94, addressingMode: ZeroPageX, cycles: 4},
	}
}

func (o sty) GetName() string {
	return "STY"
}

type tax struct{}

func (o tax) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.CPU.A
	console.CPU.X = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o tax) getVariations() []Variation {
	return []Variation{
		{opcode: 0xAA, addressingMode: nil, cycles: 2},
	}
}

func (o tax) GetName() string {
	return "TAX"
}

type txa struct{}

func (o txa) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.CPU.X
	console.CPU.A = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o txa) getVariations() []Variation {
	return []Variation{
		{opcode: 0x8A, addressingMode: nil, cycles: 2},
	}
}

func (o txa) GetName() string {
	return "TXA"
}

type tay struct{}

func (o tay) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.CPU.A
	console.CPU.Y = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o tay) getVariations() []Variation {
	return []Variation{
		{opcode: 0xA8, addressingMode: nil, cycles: 2},
	}
}

func (o tay) GetName() string {
	return "TAY"
}

type tya struct{}

func (o tya) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.CPU.Y
	console.CPU.A = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o tya) getVariations() []Variation {
	return []Variation{
		{opcode: 0x98, addressingMode: nil, cycles: 2},
	}
}

func (o tya) GetName() string {
	return "TYA"
}

type tsx struct{}

func (o tsx) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := byte(console.CPU.SP & 0x00FF)
	console.CPU.X = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o tsx) getVariations() []Variation {
	return []Variation{
		{opcode: 0xBA, addressingMode: nil, cycles: 2},
	}
}

func (o tsx) GetName() string {
	return "TSX"
}

type txs struct{}

func (o txs) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := int(console.CPU.X) | 0x0100
	console.CPU.SP = uint16(value)
	return variation.cycles
}

func (o txs) getVariations() []Variation {
	return []Variation{
		{opcode: 0x9A, addressingMode: nil, cycles: 2},
	}
}

func (o txs) GetName() string {
	return "TXS"
}

type pla struct{}

func (o pla) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.StackPopData()
	console.CPU.A = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o pla) getVariations() []Variation {
	return []Variation{
		{opcode: 0x68, addressingMode: nil, cycles: 4},
	}
}

func (o pla) GetName() string {
	return "PLA"
}

type pha struct{}

func (o pha) Exec(console *processor.Console, variation *Variation, state *State) int {
	console.StackPushData(console.CPU.A)
	return variation.cycles
}

func (o pha) getVariations() []Variation {
	return []Variation{
		{opcode: 0x48, addressingMode: nil, cycles: 3},
	}
}

func (o pha) GetName() string {
	return "PHA"
}

type plp struct{}

func (o plp) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.StackPopData()
	flags := value&processor.NotBreakBit | processor.BFlag
	console.CPU.Flags = flags
	return variation.cycles
}

func (o plp) getVariations() []Variation {
	return []Variation{
		{opcode: 0x28, addressingMode: nil, cycles: 4},
	}
}

func (o plp) GetName() string {
	return "PLP"
}

type php struct{}

func (o php) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.CPU.Flags | processor.BreakBit | processor.BFlag
	console.StackPushData(value)
	return variation.cycles
}

func (o php) getVariations() []Variation {
	return []Variation{
		{opcode: 0x08, addressingMode: nil, cycles: 3},
	}
}

func (o php) GetName() string {
	return "PHP"
}

var MoveOpCodes = []OpCode{
	lda{},
	sta{},
	ldx{},
	stx{},
	ldy{},
	sty{},
	tax{},
	txa{},
	tay{},
	tya{},
	tsx{},
	txs{},
	pla{},
	pha{},
	plp{},
	php{},
}
