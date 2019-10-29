package opcodes

import (
	"students.ic.unicamp.br/goten/processor"
)

type AddressMode interface {
	WriteTo(console *processor.Console, address uint16, value byte)
	ReadFrom(console *processor.Console, address uint16) int
	FetchAddress(console *processor.Console, state *State) (uint16, bool)
	AddressFormatString() (string, int)
}

type indirect struct{}

func (a indirect) WriteTo(console *processor.Console, address uint16, value byte) {
	// do nothing
}

func (a indirect) ReadFrom(console *processor.Console, address uint16) int {
	return 0x00 // FIXME??: this should do nothing
}

func (a indirect) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	pointer := console.FetchAddress(console.CPU.PC)
	console.CPU.PC += 2
	address := console.FetchAddress(pointer)

	/// log
	state.ParameterCount = 2
	state.Parameter1 = byte(pointer & 0x00FF)
	state.Parameter2 = byte(pointer & 0xFF00 >> 8)
	state.Address = make([]interface{}, 2)
	state.Address[0] = pointer
	state.Address[1] = address

	return address, false
}

func (a indirect) AddressFormatString() (string, int) {
	return "($%04X) = %04X", 2
}

type indirectX struct{}

func (a indirectX) WriteTo(console *processor.Console, address uint16, value byte) {
	console.StoreData(address, value)
}

func (a indirectX) ReadFrom(console *processor.Console, address uint16) int {
	data := console.FetchData(address)
	return int(data)
}

func (a indirectX) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	acc := console.FetchData(console.CPU.PC)
	console.CPU.PC++

	pointer := console.CPU.X + acc

	address := console.FetchAddress(uint16(pointer))

	/// log
	state.ParameterCount = 1
	state.Parameter1 = acc & 0x00FF
	state.Address = make([]interface{}, 3)
	state.Address[0] = acc
	state.Address[1] = pointer
	state.Address[2] = address

	return address, false
}

func (a indirectX) AddressFormatString() (string, int) {
	return "($%02X,X) @ %02X = %04X", 3
}

type indirectY struct{}

func (a indirectY) WriteTo(console *processor.Console, address uint16, value byte) {
	console.StoreData(address, value)
}

func (a indirectY) ReadFrom(console *processor.Console, address uint16) int {
	data := console.FetchData(address)
	return int(data)
}

func (a indirectY) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	pointer := uint16(console.FetchData(console.CPU.PC))
	console.CPU.PC++

	startAddressLow := uint16(console.FetchData(pointer))

	startAddressHigh := uint16(console.FetchData(processor.WrapUint16(0x00, 0xFF, pointer+1))) << 8

	baseAddress := startAddressLow + startAddressHigh
	address := baseAddress + uint16(console.CPU.Y)

	/// log
	state.ParameterCount = 1
	state.Parameter1 = byte(pointer & 0x00FF)
	state.Address = make([]interface{}, 3)
	state.Address[0] = pointer
	state.Address[1] = startAddressHigh + startAddressLow
	state.Address[2] = address

	return address, pageCross(baseAddress, address)
}

func (a indirectY) AddressFormatString() (string, int) {
	return "($%02X),Y = %04X @ %04X", 3
}

type zeroPage struct{}

func (a zeroPage) WriteTo(console *processor.Console, address uint16, value byte) {
	console.StoreData(address, value)
}

func (a zeroPage) ReadFrom(console *processor.Console, address uint16) int {
	data := console.FetchData(address)
	return int(data)
}

func (a zeroPage) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	address := uint16(console.FetchData(console.CPU.PC))
	console.CPU.PC++

	/// log
	state.ParameterCount = 1
	state.Parameter1 = byte(address & 0x00FF)
	state.Address = make([]interface{}, 1)
	state.Address[0] = address

	return address, false
}

func (a zeroPage) AddressFormatString() (string, int) {
	return "$%02X", 1
}

type zeroPageX struct{}

func (a zeroPageX) WriteTo(console *processor.Console, address uint16, value byte) {
	console.StoreData(address, value)
}

func (a zeroPageX) ReadFrom(console *processor.Console, address uint16) int {
	data := console.FetchData(address)
	return int(data)
}

func (a zeroPageX) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	baseAddress := console.FetchData(console.CPU.PC)
	console.CPU.PC++

	address := baseAddress + console.CPU.X

	/// log
	state.ParameterCount = 1
	state.Parameter1 = byte(baseAddress & 0x00FF)
	state.Address = make([]interface{}, 2)
	state.Address[0] = baseAddress
	state.Address[1] = address

	return uint16(address), false
}

func (a zeroPageX) AddressFormatString() (string, int) {
	return "$%02X,X @ %02X", 2
}

type zeroPageY struct{}

func (a zeroPageY) WriteTo(console *processor.Console, address uint16, value byte) {
	console.StoreData(address, value)
}

func (a zeroPageY) ReadFrom(console *processor.Console, address uint16) int {
	data := console.FetchData(address)
	return int(data)
}

func (a zeroPageY) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	baseAddress := console.FetchData(console.CPU.PC)
	console.CPU.PC++

	address := baseAddress + console.CPU.Y

	/// log
	state.ParameterCount = 1
	state.Parameter1 = byte(baseAddress & 0x00FF)
	state.Address = make([]interface{}, 2)
	state.Address[0] = baseAddress
	state.Address[1] = address

	return uint16(address), false
}

func (a zeroPageY) AddressFormatString() (string, int) {
	return "$%02X,Y @ %02X", 2
}

type absolute struct{}

func (a absolute) WriteTo(console *processor.Console, address uint16, value byte) {
	console.StoreData(address, value)
}

func (a absolute) ReadFrom(console *processor.Console, address uint16) int {
	data := console.FetchData(address)
	return int(data)
}

func (a absolute) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	startAddressLow := uint16(console.FetchData(console.CPU.PC))
	console.CPU.PC++

	startAddressHigh := uint16(console.FetchData(console.CPU.PC)) << 8
	console.CPU.PC++

	address := startAddressLow + startAddressHigh

	/// log
	state.ParameterCount = 2
	state.Parameter1 = byte(startAddressLow)
	state.Parameter2 = byte(startAddressHigh >> 8)
	state.Address = make([]interface{}, 1)
	state.Address[0] = address

	return address, false
}

func (a absolute) AddressFormatString() (string, int) {
	return "$%04X", 1
}
func pageCross(a uint16, b uint16) bool {
	return (a & 0xFF00) != (b & 0xFF00)
}

type absoluteY struct{}

func (a absoluteY) WriteTo(console *processor.Console, address uint16, value byte) {
	console.StoreData(address, value)
}

func (a absoluteY) ReadFrom(console *processor.Console, address uint16) int {
	data := console.FetchData(address)
	return int(data)
}

func (a absoluteY) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	startAddressLow := uint16(console.FetchData(console.CPU.PC))
	console.CPU.PC++

	startAddressHigh := uint16(console.FetchData(console.CPU.PC)) << 8
	console.CPU.PC++

	baseAddress := startAddressLow + startAddressHigh
	address := baseAddress + uint16(console.CPU.Y)

	/// log
	state.ParameterCount = 2
	state.Parameter1 = byte(startAddressLow)
	state.Parameter2 = byte(startAddressHigh >> 8)
	state.Address = make([]interface{}, 2)
	state.Address[0] = baseAddress
	state.Address[1] = address

	return address, pageCross(baseAddress, address)
}

func (a absoluteY) AddressFormatString() (string, int) {
	return "$%04X,Y @ %04X", 2
}

type absoluteX struct{}

func (a absoluteX) WriteTo(console *processor.Console, address uint16, value byte) {
	console.StoreData(address, value)
}

func (a absoluteX) ReadFrom(console *processor.Console, address uint16) int {
	data := console.FetchData(address)
	return int(data)
}

func (a absoluteX) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	startAddressLow := uint16(console.FetchData(console.CPU.PC))
	console.CPU.PC++

	startAddressHigh := uint16(console.FetchData(console.CPU.PC)) << 8
	console.CPU.PC++

	baseAddress := startAddressLow + startAddressHigh
	address := baseAddress + uint16(console.CPU.X)

	/// log
	state.ParameterCount = 2
	state.Parameter1 = byte(startAddressLow)
	state.Parameter2 = byte(startAddressHigh >> 8)
	state.Address = make([]interface{}, 2)
	state.Address[0] = baseAddress
	state.Address[1] = address

	return address, pageCross(baseAddress, address)
}

func (a absoluteX) AddressFormatString() (string, int) {
	return "$%04X,X @ %04X", 2
}

type immediate struct{}

func (a immediate) WriteTo(console *processor.Console, address uint16, value byte) {
	// do nothing
}

func (a immediate) ReadFrom(console *processor.Console, address uint16) int {
	return int(address)
}

func (a immediate) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	address := uint16(console.FetchData(console.CPU.PC))
	console.CPU.PC++

	/// log
	state.ParameterCount = 1
	state.Parameter1 = byte(address)
	state.Address = make([]interface{}, 1)
	state.Address[0] = address

	return address, false
}

func (a immediate) AddressFormatString() (string, int) {
	return "#$%02X", 1
}

type accumulator struct{}

func (a accumulator) WriteTo(console *processor.Console, address uint16, value byte) {
	console.CPU.A = value
}

func (a accumulator) ReadFrom(console *processor.Console, address uint16) int {
	return int(console.CPU.A)
}

func (a accumulator) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	/// log
	state.ParameterCount = 0
	state.Address = make([]interface{}, 0)
	return 0x00, false // FIXME??: this should do nothing
}

func (a accumulator) AddressFormatString() (string, int) {
	return "A", 0
}

type relative struct{}

func (a relative) WriteTo(console *processor.Console, address uint16, value byte) {
	// do nothing
}

func (a relative) ReadFrom(console *processor.Console, address uint16) int {
	pointer := console.FetchData(address)
	console.CPU.PC++
	return int(pointer)
}

func (a relative) FetchAddress(console *processor.Console, state *State) (uint16, bool) {
	address := uint16(console.FetchData(console.CPU.PC))
	console.CPU.PC++

	// ????
	if address >= 0x80 {
		address -= 0x100
	}

	/// log
	state.ParameterCount = 1
	state.Parameter1 = byte(address)
	state.Address = make([]interface{}, 1)
	state.Address[0] = address + console.CPU.PC

	return address, false
}

func (a relative) AddressFormatString() (string, int) {
	return "$%04X", 1
}

var (
	Indirect    = indirect{}
	IndirectX   = indirectX{}
	IndirectY   = indirectY{}
	ZeroPage    = zeroPage{}
	ZeroPageX   = zeroPageX{}
	ZeroPageY   = zeroPageY{}
	Absolute    = absolute{}
	AbsoluteY   = absoluteY{}
	AbsoluteX   = absoluteX{}
	Immediate   = immediate{}
	Accumulator = accumulator{}
	Relative    = relative{}
)
