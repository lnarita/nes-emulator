package processor

import (
	"fmt"
	"log"
)

type Console struct {
	CPU    *CPU
	PPU    *PPU
	Memory *Memory
}

func (console Console) String() string {
	return fmt.Sprintf("Console { CPU: %s, PPU: %s, MEM: %s }", console.CPU, console.PPU, console.Memory)
}

func (console *Console) Tick() {
	console.PPU.Tick()
	console.PPU.Tick()
	console.PPU.Tick()
	console.CPU.Tick()
}

func (console *Console) FetchData(address uint16) byte {
	switch {
	case address >= 0x6000 || address < 0x2000:
		return console.Memory.Read(address)
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

func (console *Console) FetchAddress(address uint16) uint16 {
	var ptrLow, ptrHigh uint16
	ptrLow = address
	if (ptrLow & 0x00FF) == 0x00FF {
		ptrHigh = address & 0xFF00 // NES hardware bug: wrap address within page
	} else {
		ptrHigh = address + 1
	}
	var low, high uint16
	low = uint16(console.FetchData(ptrLow))
	high = uint16(console.FetchData(ptrHigh))
	return (high << 8) | low
}

func (console *Console) StoreData(address uint16, data byte) {
	switch {
	case address >= 0x6000 || address < 0x2000:
		console.Memory.Write(address, data)
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

func (console *Console) StackPushData(data byte) {
	console.StoreData(console.CPU.SP, data)
	console.CPU.SP = WrapUint16(0x0100, 0x01FF, console.CPU.SP-1)
}

func (console *Console) StackPopData() byte {
	console.CPU.SP = WrapUint16(0x0100, 0x01FF, console.CPU.SP+1)
	return console.FetchData(console.CPU.SP)
}

func (console *Console) StackPushAddress(data uint16) {
	console.StackPushData(byte(data >> 8))
	console.StackPushData(byte(data & 0x00FF))
}

func (console *Console) StackPopAddress() uint16 {
	low := console.StackPopData()
	high := console.StackPopData()
	return uint16(high)<<8 | uint16(low)
}
