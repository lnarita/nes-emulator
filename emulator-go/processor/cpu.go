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

	Cycle int
}

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

func (cpu *CPU) Tick() {
	cpu.Cycle += 1
}

func Setup(memory *Memory) *CPU {
	resetAddress := memory.FetchAddress(RESET)
	cpu := CPU{Flags: 0x34, PC: resetAddress}
	return &cpu
}

func (cpu *CPU) updateFlagBit(value bool, bit byte, inv byte) {
	if value {
		cpu.Flags = cpu.Flags | bit
	} else {
		cpu.Flags = cpu.Flags & inv
	}
}
func (cpu *CPU) updateNegative(value bool) {
	cpu.updateFlagBit(value, NEGATIVE_BIT, NOT_NEGATIVE_BIT)
}
func (cpu *CPU) updateOverflow(value bool) {
	cpu.updateFlagBit(value, OVERFLOW_BIT, NOT_OVERFLOW_BIT)
}
func (cpu *CPU) updateBreakCommand(value bool) {
	cpu.updateFlagBit(value, BREAK_BIT, NOT_BREAK_BIT)
}

func (cpu *CPU) updateDecimal(value bool) {
	cpu.updateFlagBit(value, DECIMAL_BIT, NOT_DECIMAL_BIT)
}
func (cpu *CPU) updateInterruptsDisabled(value bool) {
	cpu.updateFlagBit(value, OVERFLOW_BIT, NOT_OVERFLOW_BIT)
}
func (cpu *CPU) updateZero(value bool) {
	cpu.updateFlagBit(value, ZERO_BIT, NOT_ZERO_BIT)
}
func (cpu *CPU) updateCarry(value bool) {
	cpu.updateFlagBit(value, CARRY_BIT, NOT_CARRY_BIT)
}
