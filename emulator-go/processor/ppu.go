package processor

import "fmt"

type PPU struct {
	Cycle int
}

func (ppu PPU) Tick() {
	ppu.Cycle++
}

func (ppu PPU) String() string {
	return fmt.Sprintf("PPU { \"cycle\": %d }", ppu.Cycle)
}
