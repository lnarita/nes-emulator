package processor

import (
	"log"
)

const NMI = 0xFFFA
const RESET = 0xFFFC
const IRQ = 0xFFFE

type Memory struct {
	RAM []byte
	ROM *Cartridge
}

func Load(cartridge *Cartridge) *Memory {
	return &Memory{RAM: make([]byte, 2*KB), ROM: cartridge}
}

func (mem *Memory) FetchData(address int) byte {
	switch {
	case address >= 0x6000:
		// Mappers
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

func (mem *Memory) FetchAddress(address int) int {
	var ptrLow, ptrHigh int
	ptrLow = address
	if ptrLow == 0xFF {
		ptrHigh = address & 0xFF00 // NES hardware bug: wrap address within page
	} else {
		ptrHigh = address + 1
	}
	var low, high int
	low = int(mem.FetchData(ptrLow))
	high = int(mem.FetchData(ptrHigh))
	return (high << 8) | low
}
