package processor

type PPU struct {
	Cycle int
}

func (ppu PPU) Tick() {
	ppu.Cycle++
}
