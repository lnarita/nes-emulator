package processor

type AddressMode interface {
	WriteTo(console *Console, address int, value byte)
	ReadFrom(console *Console, address int) int
	FetchAddress(console *Console) (int, bool)
}

type indirect struct{}

func (a indirect) WriteTo(console *Console, address int, value byte) {
	// do nothing
}

func (a indirect) ReadFrom(console *Console, address int) int {
	return 0x00 // FIXME??: this should do nothing
}

func (a indirect) FetchAddress(console *Console) (int, bool) {
	pointer := console.Memory.FetchAddress(console.CPU.PC)
	console.CPU.PC += 2
	address := console.Memory.FetchAddress(pointer)
	return address, false
}

type indirectX struct{}

func (a indirectX) WriteTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
}

func (a indirectX) ReadFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	return int(data)
}

func (a indirectX) FetchAddress(console *Console) (int, bool) {
	acc := console.Memory.FetchData(console.CPU.PC)
	console.CPU.PC++

	pointer := Wrap(0x00, 0xFF, int(console.CPU.X+acc))

	address := console.Memory.FetchAddress(pointer)

	return address, false
}

type indirectY struct{}

func (a indirectY) WriteTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
}

func (a indirectY) ReadFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	return int(data)
}

func (a indirectY) FetchAddress(console *Console) (int, bool) {
	pointer := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++

	startAddressLow := int(console.Memory.FetchData(pointer))

	startAddressHigh := int(console.Memory.FetchData(Wrap(0x00, 0xFF, pointer+1))) << 8

	address := startAddressLow + int(console.CPU.Y) + startAddressHigh

	return address, pageCross(address, startAddressHigh)
}

type zeroPage struct{}

func (a zeroPage) WriteTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
}

func (a zeroPage) ReadFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	return int(data)
}

func (a zeroPage) FetchAddress(console *Console) (int, bool) {
	address := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	return address, false
}

type zeroPageX struct{}

func (a zeroPageX) WriteTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
}

func (a zeroPageX) ReadFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	return int(data)
}

func (a zeroPageX) FetchAddress(console *Console) (int, bool) {
	baseAddress := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++

	address := Wrap(0x00, 0xFF, baseAddress+int(console.CPU.X))

	return address, false
}

type zeroPageY struct{}

func (a zeroPageY) WriteTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
}

func (a zeroPageY) ReadFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	return int(data)
}

func (a zeroPageY) FetchAddress(console *Console) (int, bool) {
	baseAddress := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++

	address := Wrap(0x00, 0xFF, baseAddress+int(console.CPU.Y))

	return address, false
}

type absolute struct{}

func (a absolute) WriteTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
}

func (a absolute) ReadFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	return int(data)
}

func (a absolute) FetchAddress(console *Console) (int, bool) {
	startAddressLow := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++

	startAddressHigh := int(console.Memory.FetchData(console.CPU.PC)) << 8
	console.CPU.PC++

	address := startAddressLow + startAddressHigh

	return address, pageCross(address, startAddressHigh)
}

func pageCross(address int, high int) bool {
	return (address & HighBitsMask) != (high & HighBitsMask)
}

type absoluteY struct{}

func (a absoluteY) WriteTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
}

func (a absoluteY) ReadFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	return int(data)
}

func (a absoluteY) FetchAddress(console *Console) (int, bool) {
	startAddressLow := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++

	startAddressHigh := int(console.Memory.FetchData(console.CPU.PC)) << 8
	console.CPU.PC++
	address := startAddressLow + startAddressHigh + int(console.CPU.Y)
	return address, pageCross(address, startAddressHigh)
}

type absoluteX struct{}

func (a absoluteX) WriteTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
}

func (a absoluteX) ReadFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	return int(data)
}

func (a absoluteX) FetchAddress(console *Console) (int, bool) {
	startAddressLow := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++

	startAddressHigh := int(console.Memory.FetchData(console.CPU.PC)) << 8
	console.CPU.PC++

	address := startAddressLow + startAddressHigh + int(console.CPU.X)
	return address, pageCross(address, startAddressHigh)
}

type immediate struct{}

func (a immediate) WriteTo(console *Console, address int, value byte) {
	// do nothing
}

func (a immediate) ReadFrom(console *Console, address int) int {
	return address
}

func (a immediate) FetchAddress(console *Console) (int, bool) {
	address := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	return address, false
}

type accumulator struct{}

func (a accumulator) WriteTo(console *Console, address int, value byte) {
	console.CPU.A = value
}

func (a accumulator) ReadFrom(console *Console, address int) int {
	return int(console.CPU.A)
}

func (a accumulator) FetchAddress(console *Console) (int, bool) {
	return 0x00, false // FIXME??: this should do nothing
}

type relative struct{}

func (a relative) WriteTo(console *Console, address int, value byte) {
	// do nothing
}

func (a relative) ReadFrom(console *Console, address int) int {
	pointer := console.Memory.FetchData(address)
	console.CPU.PC++
	return int(pointer)
}

func (a relative) FetchAddress(console *Console) (int, bool) {
	address := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	return address, false
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
