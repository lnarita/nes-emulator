package processor

type AddressMode interface {
	writeTo(console *Console, address int, value byte)
	readFrom(console *Console, address int) int
	fetchAddress(console *Console) int
}

type indirect struct{}

func (a indirect) writeTo(console *Console, address int, value byte) {
	// do nothing
}

func (a indirect) readFrom(console *Console, address int) int {
	return 0x00 // FIXME??: this should do nothing
}

func (a indirect) fetchAddress(console *Console) int {
	pointer := console.Memory.FetchAddress(console.CPU.PC)
	console.CPU.PC += 2
	console.Tick()
	console.Tick()
	address := console.Memory.FetchAddress(pointer)
	console.Tick()
	console.Tick()
	return address
}

type indirectX struct{}

func (a indirectX) writeTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
	console.Tick()
}

func (a indirectX) readFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	console.Tick()
	return int(data)
}

func (a indirectX) fetchAddress(console *Console) int {
	acc := console.Memory.FetchData(console.CPU.PC)
	console.CPU.PC++
	console.Tick()
	pointer := Wrap(0x00, 0xFF, int(console.CPU.X+acc))
	console.Tick()
	address := console.Memory.FetchAddress(pointer)
	console.Tick()
	console.Tick()
	return address
}

type indirectY struct{}

func (a indirectY) writeTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
	console.Tick()
}

func (a indirectY) readFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	console.Tick()
	return int(data)
}

func (a indirectY) fetchAddress(console *Console) int {
	pointer := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	console.Tick()

	startAddressLow := int(console.Memory.FetchData(pointer))
	console.Tick()
	startAddressHigh := int(console.Memory.FetchData(Wrap(0x00, 0xFF, pointer+1))) << 8
	console.Tick()
	address := startAddressLow + int(console.CPU.Y) + startAddressHigh
	if (address & 0xFF00) != (startAddressHigh & 0xFF00) {
		console.Tick() // Overflow: oops cycle
	}
	return address
}

type zeroPage struct{}

func (a zeroPage) writeTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
	console.Tick()
}

func (a zeroPage) readFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	console.Tick()
	return int(data)
}

func (a zeroPage) fetchAddress(console *Console) int {
	address := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	console.Tick()
	return address
}

type zeroPageX struct{}

func (a zeroPageX) writeTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
	console.Tick()
}

func (a zeroPageX) readFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	console.Tick()
	return int(data)
}

func (a zeroPageX) fetchAddress(console *Console) int {
	baseAddress := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	console.Tick()
	address := Wrap(0x00, 0xFF, baseAddress+int(console.CPU.X))
	console.Tick()
	return address
}

type zeroPageY struct{}

func (a zeroPageY) writeTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
	console.Tick()
}

func (a zeroPageY) readFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	console.Tick()
	return int(data)
}

func (a zeroPageY) fetchAddress(console *Console) int {
	baseAddress := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	console.Tick()
	address := Wrap(0x00, 0xFF, baseAddress+int(console.CPU.Y))
	console.Tick()
	return address
}

type absolute struct{}

func (a absolute) writeTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
	console.Tick()
}

func (a absolute) readFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	console.Tick()
	return int(data)
}

func (a absolute) fetchAddress(console *Console) int {
	startAddressLow := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	console.Tick()
	startAddressHigh := int(console.Memory.FetchData(console.CPU.PC)) << 8
	console.CPU.PC++
	console.Tick()
	address := startAddressLow + startAddressHigh
	if (address & 0xFF00) != (startAddressHigh & 0xFF00) {
		console.Tick() // Overflow: oops cycle
	}
	return address
}

type absoluteY struct{}

func (a absoluteY) writeTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
	console.Tick()
}

func (a absoluteY) readFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	console.Tick()
	return int(data)
}

func (a absoluteY) fetchAddress(console *Console) int {
	startAddressLow := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	console.Tick()
	startAddressHigh := int(console.Memory.FetchData(console.CPU.PC)) << 8
	console.CPU.PC++
	console.Tick()
	address := startAddressLow + startAddressHigh + int(console.CPU.Y)
	if (address & 0xFF00) != (startAddressHigh & 0xFF00) {
		console.Tick() // Overflow: oops cycle
	}
	return address
}

type absoluteX struct{}

func (a absoluteX) writeTo(console *Console, address int, value byte) {
	console.Memory.StoreData(address, value)
	console.Tick()
}

func (a absoluteX) readFrom(console *Console, address int) int {
	data := console.Memory.FetchData(address)
	console.Tick()
	return int(data)
}

func (a absoluteX) fetchAddress(console *Console) int {
	startAddressLow := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	console.Tick()
	startAddressHigh := int(console.Memory.FetchData(console.CPU.PC)) << 8
	console.CPU.PC++
	console.Tick()
	address := startAddressLow + startAddressHigh + int(console.CPU.X)
	if (address & 0xFF00) != (startAddressHigh & 0xFF00) {
		console.Tick() // Overflow: oops cycle
	}
	return address
}

type immediate struct{}

func (a immediate) writeTo(console *Console, address int, value byte) {
	// do nothing
}

func (a immediate) readFrom(console *Console, address int) int {
	return address
}

func (a immediate) fetchAddress(console *Console) int {
	address := int(console.Memory.FetchData(console.CPU.PC))
	console.CPU.PC++
	return address
}

type accumulator struct{}

func (a accumulator) writeTo(console *Console, address int, value byte) {
	console.CPU.A = value
}

func (a accumulator) readFrom(console *Console, address int) int {
	return int(console.CPU.A)
}

func (a accumulator) fetchAddress(console *Console) int {
	console.Tick()
	return 0x00 // FIXME??: this should do nothing
}

type relative struct{}

func (a relative) writeTo(console *Console, address int, value byte) {
	// do nothing
}

func (a relative) readFrom(console *Console, address int) int {
	pointer := console.Memory.FetchData(address)
	console.CPU.PC++
	return int(pointer)
}

func (a relative) fetchAddress(console *Console) int {
	return 0x00
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
