package processor

import (
	"log"
	"students.ic.unicamp.br/goten/common"
)

type Memory struct {
	RAM []byte
	ROM *Cartridge
}

func Load(cartridge *Cartridge) *Memory {
	return &Memory{RAM: make([]byte, 2*common.KB), ROM: cartridge}
}

func (mem *Memory) Fetch(address int) byte {
	switch {
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
	case address >= 0x6000:
		// Mappers
	default:
		log.Fatalf("unhandled cpu memory read at address: 0x%04X", address)
	}
	return 0xFF
}
