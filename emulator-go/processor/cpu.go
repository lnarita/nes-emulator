package processor

import (
	"fmt"
)

type CPU struct {
	PC    int
	SP    int
	A     byte
	X     byte
	Y     byte
	Flags byte

	hasData bool
	Address int
	Data    byte

	interrupt byte

	Cycle int
}

const (
	interruptNone = 0
	interruptNMI  = 1
	interruptIRQ  = 2
)

func (cpu CPU) String() string {
	var data string
	if cpu.hasData {
		data = fmt.Sprintf(", \"MEM[%04X]\": %02X", cpu.Address, cpu.Data)
	} else {
		data = ""
	}
	return fmt.Sprintf("CPU { \"cycle\": %d, \"pc\": 0x%04X, \"sp\": 0x%04X, \"a\": %02X, \"x\": %02X, \"y\": %02X, \"p[NV-BDIZC]\": %08b (%02X)%s }",
		cpu.Cycle, cpu.PC, cpu.SP, cpu.A, cpu.X, cpu.Y, cpu.Flags, cpu.Flags, data)
}

func (cpu CPU) Tick() {
	cpu.Cycle++
}

func Setup(memory *Memory) *CPU {
	resetAddress := memory.FetchAddress(Reset)
	cpu := CPU{Flags: 0x34, PC: resetAddress}
	cpu.ClearMemAccess()
	return &cpu
}

func (cpu CPU) LogMemAccess(address int, value byte) {
	cpu.hasData = true
	cpu.Data = value
	cpu.Address = address
}

func (cpu CPU) ClearMemAccess() {
	cpu.hasData = false
}

func (cpu CPU) updateFlagBit(value bool, bit byte, inv byte) {
	if value {
		cpu.Flags = cpu.Flags | bit
	} else {
		cpu.Flags = cpu.Flags & inv
	}
}

func (cpu CPU) SetNegative(value bool) {
	cpu.updateFlagBit(value, NegativeBit, NotNegativeBit)
}
func (cpu CPU) SetOverflow(value bool) {
	cpu.updateFlagBit(value, OverflowBit, NotOverflowBit)
}
func (cpu CPU) SetBreak(value bool) {
	cpu.updateFlagBit(value, BreakBit, NotBreakBit)
}
func (cpu CPU) SetDecimalMode(value bool) {
	cpu.updateFlagBit(value, DecimalBit, NotDecimalBit)
}
func (cpu CPU) DisableInterrupts(value bool) {
	cpu.updateFlagBit(value, InterruptsBit, NotInterruptsBit)
}
func (cpu CPU) SetZero(value bool) {
	cpu.updateFlagBit(value, ZeroBit, NotZeroBit)
}
func (cpu CPU) SetCarry(value bool) {
	cpu.updateFlagBit(value, CarryBit, NotCarryBit)
}
func (cpu CPU) SetZN(value byte) {
	cpu.SetNegative(value&NegativeBit != 0)
	cpu.SetZero(value == 0)
}
func (cpu CPU) IsNegative() bool {
	return cpu.Flags&NegativeBit != 0
}
func (cpu CPU) IsOverflow() bool {
	return cpu.Flags&OverflowBit != 0
}
func (cpu CPU) IsBreak() bool {
	return cpu.Flags&BreakBit != 0
}
func (cpu CPU) IsDecimalMode() bool {
	return cpu.Flags&DecimalBit != 0
}
func (cpu CPU) AreInterruptsDisabled() bool {
	return cpu.Flags&InterruptsBit != 0
}
func (cpu CPU) IsZero() bool {
	return cpu.Flags&ZeroBit != 0
}
func (cpu CPU) HasCarry() bool {
	return cpu.Flags&CarryBit != 0
}
func (cpu CPU) TriggerNMI() {
	cpu.interrupt = interruptNMI
}
func (cpu CPU) TriggerIRQ() {
	if !cpu.AreInterruptsDisabled() {
		cpu.interrupt = interruptIRQ
	}
}
