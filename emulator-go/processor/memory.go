package processor

import (
	"fmt"
	"log"
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

func (mem *Memory) FetchData(address uint16) byte {
	switch {
	case address >= 0x6000:
		return mem.mapper.Read(address)
	case address < 0x2000:
		return mem.RAM[address%0x0800]
	case address < 0x4000:
		// PPU registers
	case address == 0x4014:
		// PPU registers
	case address == 0x4015:
		// APU registers
	case address == 0x4016:
		// Controller 1
	case address == 0x4017:
		// Controller 2
	case address < 0x6000:
		// I/O registers
	default:
		log.Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return 0xFF
}

func (mem *Memory) FetchAddress(address uint16) uint16 {
	var ptrLow, ptrHigh uint16
	ptrLow = address
	if (ptrLow & 0x00FF) == 0x00FF {
		ptrHigh = address & 0xFF00 // NES hardware bug: wrap address within page
	} else {
		ptrHigh = address + 1
	}
	var low, high uint16
	low = uint16(mem.FetchData(ptrLow))
	high = uint16(mem.FetchData(ptrHigh))
	return (high << 8) | low
}

func (mem *Memory) StoreData(address uint16, data byte) {
	switch {
	case address >= 0x6000:
		mem.mapper.Write(address, data)
	case address < 0x2000:
		mem.RAM[address%0x0800] = data
	case address < 0x4000:
		// PPU registers
	case address == 0x4014:
		// PPU registers
	case address == 0x4015:
		// APU registers
	case address == 0x4016:
		// Controller 1
	case address == 0x4017:
		// Controller 2
	case address < 0x6000:
		// I/O registers
	default:
		log.Fatalf("unhandled cpu memory write at address: 0x%04X", address)
	}
}

func (mem *Memory) StackPushData(cpu *CPU, data byte) {
	mem.StoreData(cpu.SP, data)
	cpu.SP = WrapUint16(0x0100, 0x01FF, cpu.SP-1)
}

func (mem *Memory) StackPopData(cpu *CPU) byte {
	cpu.SP = WrapUint16(0x0100, 0x01FF, cpu.SP+1)
	return mem.FetchData(cpu.SP)
}

func (mem *Memory) StackPushAddress(cpu *CPU, data uint16) {
	mem.StackPushData(cpu, byte(data>>8))
	mem.StackPushData(cpu, byte(data&0x00FF))
}

func (mem *Memory) StackPopAddress(cpu *CPU) uint16 {
	low := mem.StackPopData(cpu)
	high := mem.StackPopData(cpu)
	return uint16(high)<<8 | uint16(low)
}
