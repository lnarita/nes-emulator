package opcodes

import "students.ic.unicamp.br/goten/processor"

type ign struct{}

func (o ign) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	if stall {
		cycleAcc++
	}
	return variation.cycles + cycleAcc
}

func (o ign) getVariations() []Variation {
	return []Variation{
		{opcode: 0x0C, addressingMode: Absolute, cycles: 4},
		{opcode: 0x04, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x44, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x64, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x14, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0x34, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0x54, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0x74, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0xD4, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0xF4, addressingMode: ZeroPageX, cycles: 4},
		{opcode: 0x1C, addressingMode: AbsoluteX, cycles: 4},
		{opcode: 0x3C, addressingMode: AbsoluteX, cycles: 4},
		{opcode: 0x5C, addressingMode: AbsoluteX, cycles: 4},
		{opcode: 0x7C, addressingMode: AbsoluteX, cycles: 4},
		{opcode: 0xDC, addressingMode: AbsoluteX, cycles: 4},
		{opcode: 0xFC, addressingMode: AbsoluteX, cycles: 4},
	}
}

func (o ign) GetName() string {
	return "*NOP"
}

type skb struct{}

func (o skb) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)

	value := byte(variation.addressingMode.ReadFrom(console, address))

	/// log
	state.hasData = true
	state.data = value

	if stall {
		cycleAcc++
	}
	return variation.cycles + cycleAcc
}

func (o skb) getVariations() []Variation {
	return []Variation{
		{opcode: 0x80, addressingMode: Immediate, cycles: 2},
		{opcode: 0x82, addressingMode: Immediate, cycles: 2},
		{opcode: 0x89, addressingMode: Immediate, cycles: 2},
		{opcode: 0xC2, addressingMode: Immediate, cycles: 2},
		{opcode: 0xE2, addressingMode: Immediate, cycles: 2},
	}
}

func (o skb) GetName() string {
	return "*NOP"
}

type unofficialNop struct{}

func (o unofficialNop) Exec(console *processor.Console, variation *Variation, state *State) int {
	return variation.cycles
}

func (o unofficialNop) getVariations() []Variation {
	return []Variation{
		{opcode: 0x1A, addressingMode: nil, cycles: 2},
		{opcode: 0x3A, addressingMode: nil, cycles: 2},
		{opcode: 0x5A, addressingMode: nil, cycles: 2},
		{opcode: 0x7A, addressingMode: nil, cycles: 2},
		{opcode: 0xDA, addressingMode: nil, cycles: 2},
		{opcode: 0xFA, addressingMode: nil, cycles: 2},
	}
}

func (o unofficialNop) GetName() string {
	return "*NOP"
}

type lax struct{}

func (o lax) Exec(console *processor.Console, variation *Variation, state *State) int {
	// LDA
	var stall bool = false
	var cycleAcc int = 0
	if variation.addressingMode != nil {
		var address uint16
		address, stall = variation.addressingMode.FetchAddress(console, state)
		value := byte(variation.addressingMode.ReadFrom(console, address))

		/// log
		state.hasData = true
		state.data = value

		console.CPU.A = value
		console.CPU.SetZN(value)
	}

	if stall {
		cycleAcc++
	}

	// TAX
	value := console.CPU.A
	console.CPU.X = value
	console.CPU.SetZN(value)

	return variation.cycles + cycleAcc
}

func (o lax) getVariations() []Variation {
	return []Variation{
		{opcode: 0xA3, addressingMode: IndirectX, cycles: 6},
		{opcode: 0xA7, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0xAF, addressingMode: Absolute, cycles: 4},
		{opcode: 0xB3, addressingMode: IndirectY, cycles: 5},
		{opcode: 0xB7, addressingMode: ZeroPageY, cycles: 4},
		{opcode: 0xBF, addressingMode: AbsoluteY, cycles: 4},
	}
}

func (o lax) GetName() string {
	return "*LAX"
}

type sax struct{}

func (o sax) Exec(console *processor.Console, variation *Variation, state *State) int {
	var stall bool = false
	var cycleAcc int = 0

	value := console.CPU.A & console.CPU.X
	if variation.addressingMode != nil {
		var address uint16
		address, stall = variation.addressingMode.FetchAddress(console, state)

		/// log
		state.hasData = true
		state.data = console.FetchData(address)

		variation.addressingMode.WriteTo(console, address, value)
	}

	if stall {
		cycleAcc++
	}

	return variation.cycles + cycleAcc
}

func (o sax) getVariations() []Variation {
	return []Variation{
		{opcode: 0x83, addressingMode: IndirectX, cycles: 6},
		{opcode: 0x87, addressingMode: ZeroPage, cycles: 3},
		{opcode: 0x8F, addressingMode: Absolute, cycles: 4},
		{opcode: 0x97, addressingMode: ZeroPageY, cycles: 4},
	}
}

func (o sax) GetName() string {
	return "*SAX"
}

type unofficialSbc struct{}

func (o unofficialSbc) Exec(console *processor.Console, variation *Variation, state *State) int {
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

func (o unofficialSbc) getVariations() []Variation {
	return []Variation{
		{opcode: 0xEB, addressingMode: Immediate, cycles: 2},
	}
}

func (o unofficialSbc) GetName() string {
	return "*SBC"
}

type dcp struct{}

func (o dcp) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)
	value := processor.WrapInt(0x00, 0xFF, oldValue-1)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	variation.addressingMode.WriteTo(console, address, byte(value))

	sub := int(console.CPU.A)
	result := sub - value
	a := uint8(result)
	carry := sub >= int(a)

	console.CPU.SetCarry(carry)
	console.CPU.SetZN(a)

	if stall {
		cycleAcc++
	}

	// uhh, write instructions don't stall???
	return variation.cycles
}

func (o dcp) getVariations() []Variation {
	return []Variation{
		{opcode: 0xC3, addressingMode: IndirectX, cycles: 8},
		{opcode: 0xC7, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0xCF, addressingMode: Absolute, cycles: 6},
		{opcode: 0xD3, addressingMode: IndirectY, cycles: 8},
		{opcode: 0xD7, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0xDB, addressingMode: AbsoluteY, cycles: 7},
		{opcode: 0xDF, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o dcp) GetName() string {
	return "*DCP"
}

type isc struct{}

func (o isc) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	value := byte(processor.WrapInt(0x00, 0xFF, oldValue+1))
	variation.addressingMode.WriteTo(console, address, value)

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

	// uhh, write instructions don't stall???
	return variation.cycles
}

func (o isc) getVariations() []Variation {
	return []Variation{
		{opcode: 0xE3, addressingMode: IndirectX, cycles: 8},
		{opcode: 0xE7, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0xEF, addressingMode: Absolute, cycles: 6},
		{opcode: 0xF3, addressingMode: IndirectY, cycles: 8},
		{opcode: 0xF7, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0xFB, addressingMode: AbsoluteY, cycles: 7},
		{opcode: 0xFF, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o isc) GetName() string {
	return "*ISB"
}

type slo struct{}

func (o slo) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	variation.addressingMode.WriteTo(console, address, byte(oldValue))
	value := oldValue << 1
	variation.addressingMode.WriteTo(console, address, byte(value))
	carry := (value & 0xFF00) > 0

	console.CPU.A |= byte(value)
	console.CPU.SetCarry(carry)
	console.CPU.SetZN(console.CPU.A)

	if stall {
		cycleAcc++
	}

	// uhh, write instructions don't stall???
	return variation.cycles
}

func (o slo) getVariations() []Variation {
	return []Variation{
		{opcode: 0x03, addressingMode: IndirectX, cycles: 8},
		{opcode: 0x07, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0x0F, addressingMode: Absolute, cycles: 6},
		{opcode: 0x13, addressingMode: IndirectY, cycles: 8},
		{opcode: 0x17, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0x1B, addressingMode: AbsoluteY, cycles: 7},
		{opcode: 0x1F, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o slo) GetName() string {
	return "*SLO"
}

type rla struct{}

func (o rla) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	variation.addressingMode.WriteTo(console, address, byte(oldValue))
	value := oldValue << 1 & 0x00FF
	if console.CPU.HasCarry() {
		value |= 0b0000_0001
	} else {
		value &= 0b1111_1110
	}
	variation.addressingMode.WriteTo(console, address, byte(value))
	carry := (oldValue & 0b1000_0000) > 0

	console.CPU.A &= byte(value)
	console.CPU.SetCarry(carry)
	console.CPU.SetZN(console.CPU.A)

	if stall {
		cycleAcc++
	}

	// uhh, write instructions don't stall???
	return variation.cycles
}

func (o rla) getVariations() []Variation {
	return []Variation{
		{opcode: 0x23, addressingMode: IndirectX, cycles: 8},
		{opcode: 0x27, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0x2F, addressingMode: Absolute, cycles: 6},
		{opcode: 0x33, addressingMode: IndirectY, cycles: 8},
		{opcode: 0x37, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0x3B, addressingMode: AbsoluteY, cycles: 7},
		{opcode: 0x3F, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o rla) GetName() string {
	return "*RLA"
}

type sre struct{}

func (o sre) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	variation.addressingMode.WriteTo(console, address, byte(oldValue))
	value := oldValue >> 1 & 0x00FF
	variation.addressingMode.WriteTo(console, address, byte(value))
	carry := (oldValue & 0b0000_0001) > 0

	console.CPU.A ^= byte(value)
	console.CPU.SetCarry(carry)
	console.CPU.SetZN(console.CPU.A)

	if stall {
		cycleAcc++
	}

	// uhh, write instructions don't stall???
	return variation.cycles
}

func (o sre) getVariations() []Variation {
	return []Variation{
		{opcode: 0x43, addressingMode: IndirectX, cycles: 8},
		{opcode: 0x47, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0x4F, addressingMode: Absolute, cycles: 6},
		{opcode: 0x53, addressingMode: IndirectY, cycles: 8},
		{opcode: 0x57, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0x5B, addressingMode: AbsoluteY, cycles: 7},
		{opcode: 0x5F, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o sre) GetName() string {
	return "*SRE"
}

type rra struct{}

func (o rra) Exec(console *processor.Console, variation *Variation, state *State) int {
	var cycleAcc int = 0
	address, stall := variation.addressingMode.FetchAddress(console, state)
	oldValue := variation.addressingMode.ReadFrom(console, address)

	/// log
	state.hasData = true
	state.data = byte(oldValue)

	variation.addressingMode.WriteTo(console, address, byte(oldValue))
	value := oldValue >> 1 & 0x00FF
	if console.CPU.HasCarry() {
		value |= 0b1000_0000
	} else {
		value &= 0b0111_1111
	}
	variation.addressingMode.WriteTo(console, address, byte(value))
	carry := (oldValue & 0b0000_0001) > 0

	add := int(console.CPU.A)

	result := add + value

	if carry {
		result++
	}
	carry = result>>8 != 0

	console.CPU.A = byte(result)

	overflow := (add>>7 == value>>7) && (add>>7) != int(console.CPU.A)>>7
	console.CPU.SetOverflow(overflow)
	console.CPU.SetCarry(carry)
	console.CPU.SetZN(console.CPU.A)

	if stall {
		cycleAcc++
	}

	// uhh, write instructions don't stall???
	return variation.cycles
}

func (o rra) getVariations() []Variation {
	return []Variation{
		{opcode: 0x63, addressingMode: IndirectX, cycles: 8},
		{opcode: 0x67, addressingMode: ZeroPage, cycles: 5},
		{opcode: 0x6F, addressingMode: Absolute, cycles: 6},
		{opcode: 0x73, addressingMode: IndirectY, cycles: 8},
		{opcode: 0x77, addressingMode: ZeroPageX, cycles: 6},
		{opcode: 0x7B, addressingMode: AbsoluteY, cycles: 7},
		{opcode: 0x7F, addressingMode: AbsoluteX, cycles: 7},
	}
}

func (o rra) GetName() string {
	return "*RRA"
}

var UnofficialOpcodes = []OpCode{
	ign{},
	unofficialNop{},
	skb{},
	lax{},
	sax{},
	unofficialSbc{},
	dcp{},
	isc{},
	slo{},
	rla{},
	sre{},
	rra{},
}
