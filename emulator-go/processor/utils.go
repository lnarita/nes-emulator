package processor

func Tick(cpu CPU, ppu PPU) {
	ppu.Tick()
	ppu.Tick()
	ppu.Tick()
	cpu.Tick()
}
