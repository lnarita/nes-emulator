package processor

import "fmt"

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

func (cpu CPU) Tick() {
	cpu.Cycle += 1
}

func Setup(memory *Memory) (*CPU, error) {
	cpu := CPU{Flags: 0x34}
	return &cpu, nil
}
