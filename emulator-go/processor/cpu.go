package processor

import (
	"fmt"
)

type CPU struct {
	PC    uint16
	SP    uint16
	A     byte
	X     byte
	Y     byte
	Flags byte

	interrupt byte
	Stall     int

	Cycle int
}

const (
	interruptNone = 0
	interruptNMI  = 1
	interruptIRQ  = 2
)

func (cpu CPU) String() string {
	return fmt.Sprintf("CPU { \"cycle\": %d, \"pc\": 0x%04X, \"sp\": 0x%04X, \"a\": %02X, \"x\": %02X, \"y\": %02X, \"p[NV-BDIZC]\": %08b (%02X) }",
		cpu.Cycle, cpu.PC, cpu.SP, cpu.A, cpu.X, cpu.Y, cpu.Flags, cpu.Flags)
}

func (cpu *CPU) Tick() {
	cpu.Cycle++
}

func Setup(memory *Memory, automation bool) *CPU {
	cpu := &CPU{}
	//cpu := CPU{Flags: 0x34, PC: resetAddress, SP: 0x1FF}
	cpu.Reset(memory)
	if automation {
		cpu.PC = 0xC000
		cpu.Cycle = 7
	}
	return cpu
}

func (cpu *CPU) Reset(memory *Memory) {
	resetAddressLow := uint16(memory.Read(Reset))
	resetAddressHigh := uint16(memory.Read(Reset + 1))
	resetAddress := resetAddressHigh<<8 | resetAddressLow
	cpu.PC = resetAddress
	cpu.Flags = 0x24
	cpu.SP = 0x1FF
}

func (cpu *CPU) updateFlagBit(value bool, bit byte, inv byte) {
	if value {
		cpu.Flags |= bit
	} else {
		cpu.Flags &= inv
	}
}

func (cpu *CPU) SetNegative(value bool) {
	cpu.updateFlagBit(value, NegativeBit, NotNegativeBit)
}
func (cpu *CPU) SetOverflow(value bool) {
	cpu.updateFlagBit(value, OverflowBit, NotOverflowBit)
}
func (cpu *CPU) SetBreak(value bool) {
	cpu.updateFlagBit(value, BreakBit, NotBreakBit)
}
func (cpu *CPU) SetDecimalMode(value bool) {
	cpu.updateFlagBit(value, DecimalBit, NotDecimalBit)
}
func (cpu *CPU) DisableInterrupts(value bool) {
	cpu.updateFlagBit(value, InterruptsBit, NotInterruptsBit)
}
func (cpu *CPU) SetZero(value bool) {
	cpu.updateFlagBit(value, ZeroBit, NotZeroBit)
}
func (cpu *CPU) SetCarry(value bool) {
	cpu.updateFlagBit(value, CarryBit, NotCarryBit)
}
func (cpu *CPU) SetZN(value byte) {
	cpu.SetNegative(value&NegativeBit != 0)
	cpu.SetZero(value == 0)
}
func (cpu *CPU) IsNegative() bool {
	return cpu.Flags&NegativeBit != 0
}
func (cpu *CPU) IsOverflow() bool {
	return cpu.Flags&OverflowBit != 0
}
func (cpu *CPU) IsBreak() bool {
	return cpu.Flags&BreakBit != 0
}
func (cpu *CPU) IsDecimalMode() bool {
	return cpu.Flags&DecimalBit != 0
}
func (cpu *CPU) AreInterruptsDisabled() bool {
	return cpu.Flags&InterruptsBit != 0
}
func (cpu *CPU) IsZero() bool {
	return cpu.Flags&ZeroBit != 0
}
func (cpu *CPU) HasCarry() bool {
	return cpu.Flags&CarryBit != 0
}
