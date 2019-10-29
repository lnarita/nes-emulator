package processor

import (
	"fmt"
)

const NMI = 0xFFFA
const Reset = 0xFFFC
const IRQ = 0xFFFE

type Memory struct {
	RAM    []byte
	mapper Mapper
}

func (mem Memory) String() string {
	return fmt.Sprintf("MEM { RAM: [% X], ROM: %s }", mem.RAM, mem.mapper)
}

func WrapInt(start, end, value int) int {
	return start + (value-start)%((end+1)-start)
}

func WrapUint16(start, end, value uint16) uint16 {
	return start + (value-start)%((end+1)-start)
}

func Load(cartridge *Cartridge) *Memory {
	return &Memory{RAM: make([]byte, 2*KB), mapper: CreateMapper(cartridge)}
}

func (mem *Memory) Read(address uint16) byte {
	switch {
	case address >= 0x6000:
		return mem.mapper.Read(address)
	default:
		return mem.RAM[address%0x0800]
	}
}

func (mem *Memory) Write(address uint16, data byte) {
	switch {
	case address >= 0x6000:
		mem.mapper.Write(address, data)
	default:
		mem.RAM[address%0x0800] = data
	}
}
