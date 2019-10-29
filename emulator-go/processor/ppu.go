package processor

import "fmt"

type PPU struct {
	Cycle    int
	ScanLine int
	Frame    uint64

	palettes  [32]byte // 4 colours each; 4 background palettes; 4 foreground palettes
	nameTable [2048]byte
	oam       [256]byte

	V uint16
	T uint16
	X byte
	W byte
	F byte

	register byte
}

func (ppu PPU) Tick() {
	ppu.Cycle++
}

func (ppu PPU) String() string {
	return fmt.Sprintf("PPU { \"cycle\": %d }", ppu.Cycle)
}
