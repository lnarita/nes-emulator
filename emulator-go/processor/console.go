package processor

import (
	"fmt"
	"log"
)

type Console struct {
	CPU         *CPU
	PPU         *PPU
	Memory      *Memory
	Controller1 *Controller
	Controller2 *Controller
}

func (console Console) String() string {
	return fmt.Sprintf("Console { CPU: %s, PPU: %s, MEM: %s }", console.CPU, console.PPU, console.Memory)
}

func (console *Console) Tick() {
	console.PPU.ExecCycle()
	console.CPU.Tick()
}

func (console *Console) Reset() {
	console.CPU.Reset(console.Memory)
	console.PPU.Reset()
}

func (console *Console) FetchData(address uint16) byte {
	switch {
	case address >= 0x6000 || address < 0x2000:
		return console.Memory.Read(address)
	case address < 0x4000:
		return console.PPU.Read(0x2000 + (address % 8))
	case address == 0x4014:
		return console.PPU.Read(address)
	case address == 0x4015:
		// APU registers
	case address == 0x4016:
		return console.Controller1.Read()
	case address == 0x4017:
		return console.Controller2.Read()
	case address < 0x6000:
		// I/O registers
	default:
		log.Fatalf("unhandled cpu memory read at address: 0x%04X\n", address)
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
		console.PPU.Write(0x2000+(address%8), data)
	case address == 0x4014:
		console.PPU.Write(address, data)
	case address == 0x4015:
		// APU registers
	case address == 0x4016:
		console.Controller1.Write(data)
		console.Controller2.Write(data)
	case address == 0x4017:
		// APU ???
	case address < 0x6000:
		// I/O registers
	default:
		log.Fatalf("unhandled cpu memory write at address: 0x%04X\n", address)
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

func (console *Console) TriggerNMI() {
	console.CPU.interrupt = interruptNMI
}

func (console *Console) TriggerIRQ() {
	if !console.CPU.AreInterruptsDisabled() {
		console.CPU.interrupt = interruptIRQ
	}
}

func (console *Console) ClearInterrupt() {
	console.CPU.interrupt = interruptNone
}

func (console *Console) CheckInterrupts() {
	switch console.CPU.interrupt {
	case interruptNMI:
		console.nmi()
		log.Println("NMI trigged!!!")
	}
	console.ClearInterrupt()
}

func (console *Console) nmi() {
	console.StackPushAddress(console.CPU.PC)
	value := console.CPU.Flags | BreakBit | BFlag
	console.StackPushData(value)
	console.CPU.PC = console.FetchAddress(0xFFFA)
	console.CPU.DisableInterrupts(true)
	for i := 0; i < 7; i++ {
		console.CPU.Cycle += 7
		//console.Tick()
	}
}
