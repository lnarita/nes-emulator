package processor

type AddressMode interface {
	writeTo(*CPU, *Memory, int, int)
	readFrom(*CPU, *Memory, int)
	fetchAddress(*CPU, *Memory)
	read16BitsLow(*Memory, int)
	read16BitsHigh(*Memory, int)
	get16BitsAddrFromHighLow(byte, byte)
}

type indirect struct{}

func (a indirect) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a indirect) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a indirect) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a indirect) read16BitsLow(memory *Memory, addr int) {

}

func (a indirect) read16BitsHigh(memory *Memory, addr int) {

}

func (a indirect) get16BitsAddrFromHighLow(high byte, low byte) {

}

type indirectX struct{}

func (a indirectX) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a indirectX) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a indirectX) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a indirectX) read16BitsLow(memory *Memory, addr int) {

}

func (a indirectX) read16BitsHigh(memory *Memory, addr int) {

}

func (a indirectX) get16BitsAddrFromHighLow(high byte, low byte) {

}

type indirectY struct{}

func (a indirectY) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a indirectY) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a indirectY) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a indirectY) read16BitsLow(memory *Memory, addr int) {

}

func (a indirectY) read16BitsHigh(memory *Memory, addr int) {

}

func (a indirectY) get16BitsAddrFromHighLow(high byte, low byte) {

}

type zeroPage struct{}

func (a zeroPage) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a zeroPage) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a zeroPage) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a zeroPage) read16BitsLow(memory *Memory, addr int) {

}

func (a zeroPage) read16BitsHigh(memory *Memory, addr int) {

}

func (a zeroPage) get16BitsAddrFromHighLow(high byte, low byte) {

}

type zeroPageX struct{}

func (a zeroPageX) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a zeroPageX) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a zeroPageX) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a zeroPageX) read16BitsLow(memory *Memory, addr int) {

}

func (a zeroPageX) read16BitsHigh(memory *Memory, addr int) {

}

func (a zeroPageX) get16BitsAddrFromHighLow(high byte, low byte) {

}

type zeroPageY struct{}

func (a zeroPageY) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a zeroPageY) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a zeroPageY) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a zeroPageY) read16BitsLow(memory *Memory, addr int) {

}

func (a zeroPageY) read16BitsHigh(memory *Memory, addr int) {

}

func (a zeroPageY) get16BitsAddrFromHighLow(high byte, low byte) {

}

type absolute struct{}

func (a absolute) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a absolute) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a absolute) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a absolute) read16BitsLow(memory *Memory, addr int) {

}

func (a absolute) read16BitsHigh(memory *Memory, addr int) {

}

func (a absolute) get16BitsAddrFromHighLow(high byte, low byte) {

}

type absoluteY struct{}

func (a absoluteY) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a absoluteY) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a absoluteY) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a absoluteY) read16BitsLow(memory *Memory, addr int) {

}

func (a absoluteY) read16BitsHigh(memory *Memory, addr int) {

}

func (a absoluteY) get16BitsAddrFromHighLow(high byte, low byte) {

}

type absoluteX struct{}

func (a absoluteX) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a absoluteX) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a absoluteX) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a absoluteX) read16BitsLow(memory *Memory, addr int) {

}

func (a absoluteX) read16BitsHigh(memory *Memory, addr int) {

}

func (a absoluteX) get16BitsAddrFromHighLow(high byte, low byte) {

}

type immediate struct{}

func (a immediate) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a immediate) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a immediate) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a immediate) read16BitsLow(memory *Memory, addr int) {

}

func (a immediate) read16BitsHigh(memory *Memory, addr int) {

}

func (a immediate) get16BitsAddrFromHighLow(high byte, low byte) {

}

type accumulator struct{}

func (a accumulator) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a accumulator) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a accumulator) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a accumulator) read16BitsLow(memory *Memory, addr int) {

}

func (a accumulator) read16BitsHigh(memory *Memory, addr int) {

}

func (a accumulator) get16BitsAddrFromHighLow(high byte, low byte) {

}

type relative struct{}

func (a relative) writeTo(cpu *CPU, memory *Memory, addr int, value int) {

}

func (a relative) readFrom(cpu *CPU, memory *Memory, addr int) {

}

func (a relative) fetchAddress(cpu *CPU, memory *Memory) {

}

func (a relative) read16BitsLow(memory *Memory, addr int) {

}

func (a relative) read16BitsHigh(memory *Memory, addr int) {

}

func (a relative) get16BitsAddrFromHighLow(high byte, low byte) {

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
