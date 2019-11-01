package opcodes

import "students.ic.unicamp.br/goten/processor"

type ora struct{}

func (o ora) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	console.CPU.A |= value
	console.CPU.SetZN(console.CPU.A)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o ora) getVariations() []Variation {
	return []Variation{
		{opcode: 0x01, addressingMode: IndirectX, cycles: 6},
		{opcode: 0x05, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x09, addressingMode: Immediate, cycles: 2},
		{opcode: 0x0D, addressingMode: Absolute, cycles: 4},
		{opcode: 0x11, addressingMode: IndirectY, cycles: 5},
		{opcode: 0x15, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0x19, addressingMode: AbsoluteY, cycles: 4},
		{opcode: 0x1D, addressingMode: AbsoluteX, cycles: 4},
	}
}

func (o ora) GetName() string {
	return "ORA"
}

type and struct{}

func (o and) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	console.CPU.A &= value
	console.CPU.SetZN(console.CPU.A)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o and) getVariations() []Variation {
	return []Variation{
		{opcode: 0x21, addressingMode: IndirectX, cycles: 6},
		{opcode: 0x25, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x29, addressingMode: Immediate, cycles: 2},
		{opcode: 0x2D, addressingMode: Absolute, cycles: 4},
		{opcode: 0x31, addressingMode: IndirectY, cycles: 5},
		{opcode: 0x35, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0x39, addressingMode: AbsoluteY, cycles: 4},
		{opcode: 0x3D, addressingMode: AbsoluteX, cycles: 4},
	}
}

func (o and) GetName() string {
	return "AND"
}

type eor struct{}

func (o eor) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	console.CPU.A ^= value
	console.CPU.SetZN(console.CPU.A)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o eor) getVariations() []Variation {
	return []Variation{
		{opcode: 0x41, addressingMode: IndirectX, cycles: 6},
		{opcode: 0x45, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x49, addressingMode: Immediate, cycles: 2},
		{opcode: 0x4D, addressingMode: Absolute, cycles: 4},
		{opcode: 0x51, addressingMode: IndirectY, cycles: 5},
		{opcode: 0x55, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0x59, addressingMode: AbsoluteY, cycles: 4},
		{opcode: 0x5D, addressingMode: AbsoluteX, cycles: 4},
	}
}

func (o eor) GetName() string {
	return "EOR"
}

type adc struct{}

func (o adc) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	add := console.CPU.A

	result := add + value + (console.CPU.Flags & processor.CarryBit)
	signal1 := (add >> 7) & 0x00FF
	signal2 := (value >> 7) & 0x00FF
	signalResult := (result >> 7) & 0x00FF
	overflow := (signal1 == signal2) && (signal1 != signalResult)
	carry := add > result

	console.CPU.A = result
	console.CPU.SetCarry(carry)
	console.CPU.SetOverflow(overflow)
	console.CPU.SetZN(console.CPU.A)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o adc) getVariations() []Variation {
	return []Variation{
		{opcode: 0x61, addressingMode: IndirectX, cycles: 6},
		{opcode: 0x65, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x69, addressingMode: Immediate, cycles: 2},
		{opcode: 0x6D, addressingMode: Absolute, cycles: 4},
		{opcode: 0x71, addressingMode: IndirectY, cycles: 5},
		{opcode: 0x75, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0x79, addressingMode: AbsoluteY, cycles: 4},
		{opcode: 0x7D, addressingMode: AbsoluteX, cycles: 4},
	}
}

func (o adc) GetName() string {
	return "ADC"
}

type sbc struct{}

func (o sbc) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	sub := console.CPU.A

	off := 1
	if console.CPU.HasCarry() {
		off = 0
	}
	result := int(sub) - int(value) - off
	a := uint8(result)

	overflow := ((sub^value)&processor.NegativeBit > 0) && ((sub^a)&processor.NegativeBit > 0)
	carry := result >= 0

	console.CPU.A = a
	console.CPU.SetCarry(carry)
	console.CPU.SetOverflow(overflow)
	console.CPU.SetZN(console.CPU.A)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o sbc) getVariations() []Variation {
	return []Variation{
		{opcode: 0xE1, addressingMode: IndirectX, cycles: 6},
		{opcode: 0xE5, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0xE9, addressingMode: Immediate, cycles: 2},
		{opcode: 0xED, addressingMode: Absolute, cycles: 4},
		{opcode: 0xF1, addressingMode: IndirectY, cycles: 5},
		{opcode: 0xF5, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0xF9, addressingMode: AbsoluteY, cycles: 4},
		{opcode: 0xFD, addressingMode: AbsoluteX, cycles: 4},
	}
}

func (o sbc) GetName() string {
	return "SBC"
}

type cmp struct{}

func (o cmp) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	sub := console.CPU.A

	result := int(sub) - int(value)
	a := uint8(result)
	carry := sub >= a

	console.CPU.SetCarry(carry)
	console.CPU.SetZN(a)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o cmp) getVariations() []Variation {
	return []Variation{
		{opcode: 0xC1, addressingMode: IndirectX, cycles: 6},
		{opcode: 0xC5, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0xC9, addressingMode: Immediate, cycles: 2},
		{opcode: 0xCD, addressingMode: Absolute, cycles: 4},
		{opcode: 0xD1, addressingMode: IndirectY, cycles: 5},
		{opcode: 0xD5, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0xD9, addressingMode: AbsoluteY, cycles: 4},
		{opcode: 0xDD, addressingMode: AbsoluteX, cycles: 4},
	}
}

func (o cmp) GetName() string {
	return "CMP"
}

type cpx struct{}

func (o cpx) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	sub := console.CPU.X

	result := int(sub) - int(value)
	a := uint8(result)
	carry := sub >= a

	console.CPU.SetCarry(carry)
	console.CPU.SetZN(a)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o cpx) getVariations() []Variation {
	return []Variation{
		{opcode: 0xE0, addressingMode: Immediate, cycles: 2},
		{opcode: 0xE4, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0xEC, addressingMode: Absolute, cycles: 4},
	}
}

func (o cpx) GetName() string {
	return "CPX"
}

type cpy struct{}

func (o cpy) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	sub := console.CPU.Y

	result := int(sub) - int(value)
	a := uint8(result)
	carry := sub >= a

	console.CPU.SetCarry(carry)
	console.CPU.SetZN(a)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o cpy) getVariations() []Variation {
	return []Variation{
		{opcode: 0xC0, addressingMode: Immediate, cycles: 2},
		{opcode: 0xC4, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0xCC, addressingMode: Absolute, cycles: 4},
	}
}

func (o cpy) GetName() string {
	return "CPY"
}

type dec struct{}

func (o dec) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = oldValue

	value := oldValue - 1

	variation.addressingMode.WriteTo(console, address, value)
	console.CPU.SetZN(value)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o dec) getVariations() []Variation {
	return []Variation{
		{opcode: 0xC6, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0xCE, addressingMode: Absolute, cycles: 6},
		{opcode: 0xD6, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0xDE, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o dec) GetName() string {
	return "DEC"
}

type dex struct{}

func (o dex) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.CPU.X - 1

	console.CPU.X = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o dex) getVariations() []Variation {
	return []Variation{
		{opcode: 0xCA, addressingMode: nil, cycles: 2},
	}
}

func (o dex) GetName() string {
	return "DEX"
}

type dey struct{}

func (o dey) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.CPU.Y - 1

	console.CPU.Y = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o dey) getVariations() []Variation {
	return []Variation{
		{opcode: 0x88, addressingMode: nil, cycles: 2},
	}
}

func (o dey) GetName() string {
	return "DEY"
}

type inc struct{}

func (o inc) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = oldValue

	value := oldValue + 1

	variation.addressingMode.WriteTo(console, address, value)
	console.CPU.SetZN(value)

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o inc) getVariations() []Variation {
	return []Variation{
		{opcode: 0xE6, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0xEE, addressingMode: Absolute, cycles: 6},
		{opcode: 0xF6, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0xFE, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o inc) GetName() string {
	return "INC"
}

type inx struct{}

func (o inx) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.CPU.X + 1

	console.CPU.X = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o inx) getVariations() []Variation {
	return []Variation{
		{opcode: 0xE8, addressingMode: nil, cycles: 2},
	}
}

func (o inx) GetName() string {
	return "INX"
}

type iny struct{}

func (o iny) Exec(console *processor.Console, variation *Variation, state *State) int {
	value := console.CPU.Y + 1

	console.CPU.Y = value
	console.CPU.SetZN(value)
	return variation.cycles
}

func (o iny) getVariations() []Variation {
	return []Variation{
		{opcode: 0xC8, addressingMode: nil, cycles: 2},
	}
}

func (o iny) GetName() string {
	return "INY"
}

type asl struct{}

func (o asl) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	value := oldValue << 1
	carry := (value & 0xFF00) > 0

	variation.addressingMode.WriteTo(console, address, byte(value))
	console.CPU.SetCarry(carry)
	console.CPU.SetZN(byte(value))

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o asl) getVariations() []Variation {
	return []Variation{
		{opcode: 0x06, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0x0A, addressingMode: Accumulator, cycles: 2},
		{opcode: 0x0E, addressingMode: Absolute, cycles: 6},
		{opcode: 0x16, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0x1E, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o asl) GetName() string {
	return "ASL"
}

type rol struct{}

func (o rol) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	value := oldValue << 1 & 0x00FF
	if console.CPU.HasCarry() {
		value |= 0b0000_0001
	} else {
		value &= 0b1111_1110
	}
	carry := (oldValue & 0b1000_0000) > 0

	variation.addressingMode.WriteTo(console, address, byte(value))
	console.CPU.SetCarry(carry)
	console.CPU.SetZN(byte(value))

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o rol) getVariations() []Variation {
	return []Variation{
		{opcode: 0x26, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0x2A, addressingMode: Accumulator, cycles: 2},
		{opcode: 0x2E, addressingMode: Absolute, cycles: 6},
		{opcode: 0x36, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0x3E, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o rol) GetName() string {
	return "ROL"
}

type lsr struct{}

func (o lsr) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	value := oldValue >> 1 & 0x00FF
	carry := (oldValue & 0b0000_0001) > 0

	variation.addressingMode.WriteTo(console, address, byte(value))

	console.CPU.SetCarry(carry)
	console.CPU.SetZN(byte(value))

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o lsr) getVariations() []Variation {
	return []Variation{
		{opcode: 0x46, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0x4A, addressingMode: Accumulator, cycles: 2},
		{opcode: 0x4E, addressingMode: Absolute, cycles: 6},
		{opcode: 0x56, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0x5E, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o lsr) GetName() string {
	return "LSR"
}

type ror struct{}

func (o ror) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	value := oldValue >> 1 & 0x00FF
	if console.CPU.HasCarry() {
		value |= 0b1000_0000
	} else {
		value &= 0b0111_1111
	}
	carry := (oldValue & 0b0000_0001) > 0

	variation.addressingMode.WriteTo(console, address,byte(value))
	console.CPU.SetCarry(carry)
	console.CPU.SetZN(byte(value))

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o ror) getVariations() []Variation {
	return []Variation{
		{opcode: 0x66, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0x6A, addressingMode: Accumulator, cycles: 2},
		{opcode: 0x6E, addressingMode: Absolute, cycles: 6},
		{opcode: 0x76, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0x7E, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o ror) GetName() string {
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
