package processor

import "fmt"

type PPU struct {
	Console *Console // parent pointer :/

	Cycle    int
	ScanLine int
	Frame    uint64

	palettes  [32]byte // 4 colours each; 4 background palettes; 4 foreground palettes
	nameTable [2048]byte
	oam       [256]byte

	// PPU registers
	V uint16
	T uint16
	X byte
	W byte
	F byte
}

func (ppu *PPU) ExecCycle() {
	ppu.Tick()
	ppu.Tick()
	ppu.Tick()
}

func (ppu *PPU) Tick() {
	ppu.Cycle++
}

func (ppu PPU) String() string {
	return fmt.Sprintf("PPU { \"cycle\": %d }", ppu.Cycle)
}
