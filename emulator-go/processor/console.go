package processor

import "fmt"

type Console struct {
	CPU    *CPU
	PPU    *PPU
	Memory *Memory
}

func (console Console) String() string {
	return fmt.Sprintf("Console { CPU: %s, PPU: %s, MEM: %s }", console.CPU, console.PPU, console.Memory)
}

func (console Console) Tick() {
	console.PPU.Tick()
	console.PPU.Tick()
	console.PPU.Tick()
	console.CPU.Tick()
}
