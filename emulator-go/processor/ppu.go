package processor

import (
	"fmt"
	"image/color"
)

type PPU struct {
	Console *Console // parent pointer :/

	Cycle    int // 0-340
	ScanLine int // 0-261, 0-239=visible, 240=post, 241-260=vblank, 261=pre
	Frame    uint64

	Ctrl   *PPUCTRL
	Mask   *PPUMASK
	Status *PPUSTATUS

	palettes      [32]byte // 4 colours each; 4 background palettes; 4 foreground palettes
	nameTableData [2048]byte
	oamMem        [256]byte

	oam     [8]Sprite
	copyOam [8]Sprite
	Pixels  [256 * 240]color.RGBA

	// PPUSCROLL registers
	vRamAddress          uint16
	temporaryVRamAddress uint16
	x                    byte // fine x scroll (3 bit)
	latch                bool // write latch
	odd                  bool // even/odd frame flag

	oamAddress   byte
	bufferedData byte

	mirroring byte // mirroring mode. 0: horizontal; 1: vertical; 2: 4-screen VRAM
}

type Sprite struct {
	id         uint8 // index in oam
	x          uint8
	y          uint8
	tile       uint8 // index
	attributes uint8
	low        uint8 // tile data
	high       uint8 // tile data
}

type PPUCTRL struct {
	nameTable          uint16 // Nametable ($2000 / $2400 / $2800 / $2C00)
	addressIncrement   uint16 // Address increment (1 across / 32 down)
	spritePatternTable uint16 // Sprite pattern table ($0000 / $1000)
	bgPatternTable     uint16 // BG pattern table ($0000 / $1000)
	spriteSize         byte   // Sprite size (8x8 / 8x16)
	slave              bool   // PPU master/slave select (read backdrop from EXT pins; output color on EXT pins)
	nmiEnabled         bool   // Enable NMI
}

func createControlFromInt(flag uint8) *PPUCTRL {
	ctrl := &PPUCTRL{}
	ctrl.fromFlag(flag)
	return ctrl
}

func (ctrl *PPUCTRL) fromFlag(flag uint8) {
	var increment uint16
	if flag&0b0010_0000 > 0 {
		increment = 32
	} else {
		increment = 1
	}
	var sprite byte
	if flag&0b0000_0100 > 0 {
		sprite = 16
	} else {
		sprite = 8
	}
	ctrl.nameTable = uint16(flag&0b1100_0000) << 8
	ctrl.addressIncrement = increment
	ctrl.spritePatternTable = uint16(flag&0b0001_0000) >> 4 << 16
	ctrl.bgPatternTable = uint16(flag&0b0000_1000) >> 3 << 16
	ctrl.spriteSize = sprite
	ctrl.slave = flag&0b0000_0010 > 0
	ctrl.nmiEnabled = flag&0x1 == 0x1
}

func (ctrl *PPUCTRL) toFlag() uint8 {
	var flag byte = 0
	flag |= uint8(ctrl.nameTable>>8) & 0b1100_0000
	if ctrl.addressIncrement == 32 {
		flag |= 0b0010_0000
	}
	flag |= uint8(ctrl.spritePatternTable>>16) << 4
	flag |= uint8(ctrl.bgPatternTable>>16) >> 3
	if ctrl.spriteSize == 16 {
		flag |= 0b0000_0100
	}
	if ctrl.slave {
		flag |= 0b0000_0010
	}
	if ctrl.nmiEnabled {
		flag |= 0b0000_0000
	}
	return flag
}

type PPUMASK struct {
	grayscale      bool // Grayscale.
	bgLeft         bool // Show background in leftmost 8 pixels.
	sprLeft        bool // Show sprite in leftmost 8 pixels.
	showBg         bool // Show background.
	showSpr        bool // Show sprites.
	intensifyRed   bool // Intensify reds.
	intensifyGreen bool // Intensify greens.
	intensifyBlue  bool // Intensify blues.
}

func createMaskFromInt(flag uint8) *PPUMASK {
	mask := &PPUMASK{}
	mask.fromFlag(flag)
	return mask
}

func (msk *PPUMASK) fromFlag(flag uint8) {
	msk.grayscale = flag&0b1000_0000 > 0
	msk.bgLeft = flag&0b0100_0000 > 0
	msk.sprLeft = flag&0b0010_0000 > 0
	msk.showBg = flag&0b0001_0000 > 0
	msk.showSpr = flag&0b0000_1000 > 0
	msk.intensifyRed = flag&0b0000_0100 > 0
	msk.intensifyGreen = flag&0b0000_0010 > 0
	msk.intensifyBlue = flag&0b0000_0001 > 0
}

func (msk *PPUMASK) toFlag() uint8 {
	var flag byte = 0
	if msk.grayscale {
		flag |= 0b1000_0000
	}
	if msk.bgLeft {
		flag |= 0b0100_0000
	}
	if msk.sprLeft {
		flag |= 0b0010_0000
	}
	if msk.showBg {
		flag |= 0b0001_0000
	}
	if msk.showSpr {
		flag |= 0b0000_1000
	}
	if msk.intensifyRed {
		flag |= 0b0000_0100
	}
	if msk.intensifyGreen {
		flag |= 0b0000_0010
	}
	if msk.intensifyBlue {
		flag |= 0b0000_0001
	}
	return flag
}

type PPUSTATUS struct {
	bus            byte // Not significant.
	spriteOverflow bool // Sprite overflow.
	spriteHit      bool // Sprite 0 Hit.
	vBlank         bool // In VBlank?
}

func createStatusFromInt(flag uint8) *PPUSTATUS {
	status := &PPUSTATUS{}
	status.fromFlag(flag)
	return status
}

func (status *PPUSTATUS) fromFlag(flag uint8) {
	status.bus = flag & 0b0001_1111
	status.spriteOverflow = flag&0b0010_0000 > 0
	status.spriteHit = flag&0b0100_0000 > 0
	status.vBlank = flag&0b1000_0000 > 0
}

func (status *PPUSTATUS) toFlag() uint8 {
	var flag byte = 0
	flag |= status.bus & 0b0001_1111
	if status.spriteOverflow {
		flag |= 0b0010_0000
	}
	if status.spriteHit {
		flag |= 0b0100_0000
	}
	if status.vBlank {
		flag |= 0b1000_0000
	}
	return flag
}

func (ppu *PPU) ExecCycle() {
	ppu.Tick()
	ppu.Tick()
	ppu.Tick()
}

func (ppu *PPU) Tick() {
	renderingEnabled := ppu.Mask.showBg || ppu.Mask.showSpr
	// There are a total of 262 scanlines per frame
	//   Scanlines 0 to 239 are for display (NES is 256 x 240)
	//   Scanline  240 is a post-render scanline (idle)
	//   Scanlines 241 to 260 are the vblank interval
	//   Scanline  261 is a pre-render scanline
	visibleLine := ppu.ScanLine < 240
	postLine := ppu.ScanLine == 240
	nmiLine := ppu.ScanLine == 241
	preLine := ppu.ScanLine == 261

	ppu.scanlineCycle(renderingEnabled, visibleLine, postLine, nmiLine, preLine)

	ppu.Cycle++
}

func (ppu *PPU) scanlineCycle(render, visible, post, nmi, pre bool) {
	console := ppu.Console
	if nmi && ppu.Cycle == 1 {
		ppu.Status.vBlank = true
		if ppu.Ctrl.nmiEnabled {
			console.TriggerNMI()
		}
	} else if post && ppu.Cycle == 0 {
		// signal render
	} else if visible || pre {
		// Sprites
		switch ppu.Cycle {
		case 1:
			ppu.clearOAM()
			if pre {
				ppu.Status.spriteHit = false
				ppu.Status.spriteOverflow = false
			}
		case 257:
			ppu.evalSprites()
		case 321:
			ppu.loadSprites()
		}
		// BG
		switch {
		case ppu.Cycle < 2:
			if pre {

			}
		}
	}
}

func (ppu PPU) String() string {
	return fmt.Sprintf("PPU { \"cycle\": %d }", ppu.Cycle)
}

func (ppu *PPU) Reset() {
	ppu.Cycle = 340
	ppu.ScanLine = 240
	ppu.Frame = 0
	if ppu.Ctrl == nil {
		ppu.Ctrl = createControlFromInt(0)
	} else {
		ppu.Ctrl.fromFlag(0)
	}
	if ppu.Mask == nil {
		ppu.Mask = createMaskFromInt(0)
	} else {
		ppu.Mask.fromFlag(0)
	}
	if ppu.Status == nil {
		ppu.Status = createStatusFromInt(0)
	} else {
		ppu.Status.fromFlag(0)
	}
}

func (ppu *PPU) Read(address uint16) byte {
	switch address {
	case 0x2002:
		return ppu.readStatus()
	case 0x2004:
		return ppu.readOAM()
	case 0x2007:
		return ppu.readData()
	}
	return 0
}

func (ppu *PPU) readOAM() byte {
	res := ppu.oamMem[ppu.oamAddress]
	return res
}

func (ppu *PPU) readStatus() byte {
	res := ppu.Status.toFlag()
	ppu.latch = false
	ppu.Status.vBlank = false
	return res
}

func (ppu *PPU) Write(address uint16, data byte) {
	ppu.updateStatusBus(data)
	switch address {
	case 0x2000:
		ppu.writeCtrl(data)
	case 0x2001:
		ppu.writeMask(data)
	case 0x2003:
		ppu.writeOAMAddress(data)
	case 0x2004:
		ppu.writeOAMData(data)
	case 0x2005:
		ppu.writeScroll(data)
	case 0x2006:
		ppu.writeAddress(data)
	case 0x2007:
		ppu.writeData(data)
	case 0x4014:
		ppu.writeDMA(data)
	}
}

func (ppu *PPU) getCiramAddress(address uint16) uint16 {
	switch ppu.mirroring {
	//horizontal
	case 0:
		return address % 0x800
	//vertical
	case 1:
		return ((address / 2) & 0x400) + (address % 0400)
	//4-screen VRAM
	default:
		return address - 0x2000
	}
}

func (ppu *PPU) access(address uint16) byte {
	switch {
	case address < 0x2000:
		return ppu.Console.Memory.mapper.Read(address)
	case address < 0x3F00:
		return ppu.nameTableData[ppu.getCiramAddress(address)]
	case address < 0x4000:
		if (address & 0x13) == 0x10 {
			index := address & ^uint16(0x10)
			var mask byte
			if ppu.Mask.grayscale {
				mask = 0x30
			} else {
				mask = 0xFF
			}
			return ppu.palettes[index&0x1F] & mask
		}
	}
	return 0
}

func (ppu *PPU) readData() byte {
	var res byte
	if ppu.vRamAddress <= 0x3EFF {
		res = ppu.bufferedData
		ppu.bufferedData = ppu.access(ppu.vRamAddress)
	} else {
		ppu.bufferedData = ppu.access(ppu.vRamAddress)
		res = ppu.bufferedData
	}
	ppu.vRamAddress += ppu.Ctrl.addressIncrement
	return res
}

func (ppu *PPU) updateStatusBus(data byte) {
	ppu.Status.bus = data & 0x1F
}

func (ppu *PPU) writeCtrl(data byte) {
	ppu.Ctrl.fromFlag(data)
	ppu.temporaryVRamAddress = (ppu.temporaryVRamAddress & 0xF3FF) | (uint16(data) & 0x03 << 10)
}

func (ppu *PPU) writeMask(data byte) {
	ppu.Mask.fromFlag(data)
}

func (ppu *PPU) writeOAMAddress(data byte) {
	ppu.oamAddress = data
}

func (ppu *PPU) writeOAMData(data byte) {
	ppu.oamMem[ppu.oamAddress] = data
	ppu.oamAddress++
}

func (ppu *PPU) writeScroll(data byte) {
	if !ppu.latch {
		ppu.temporaryVRamAddress = (ppu.temporaryVRamAddress & 0xFFE0) | (uint16(data) >> 3)
		ppu.x = data & 0x07
	} else {
		ppu.temporaryVRamAddress = (ppu.temporaryVRamAddress & 0x8FFF) | ((uint16(data) & 0x07) << 12)
		ppu.temporaryVRamAddress = (ppu.temporaryVRamAddress & 0xFC1F) | ((uint16(data) & 0xF8) << 2)
	}
	ppu.latch = !ppu.latch
}

func (ppu *PPU) writeAddress(data byte) {
	if !ppu.latch {
		ppu.temporaryVRamAddress = (ppu.temporaryVRamAddress & 0x80FF) | ((uint16(data) & 0x3F) << 8)
	} else {
		ppu.temporaryVRamAddress = (ppu.temporaryVRamAddress & 0xFF00) | uint16(data)
		ppu.vRamAddress = ppu.temporaryVRamAddress
	}
	ppu.latch = !ppu.latch
}

func (ppu *PPU) writeData(data byte) {
	ppu.accessWrite(ppu.vRamAddress, data)
	ppu.vRamAddress += ppu.Ctrl.addressIncrement
}

func (ppu *PPU) accessWrite(address uint16, data byte) {
	switch {
	case address < 0x2000:
		ppu.Console.Memory.mapper.Write(address, data)
	case address < 0x3F00:
		ppu.nameTableData[ppu.getCiramAddress(address)] = data
	case address < 0x4000:
		if (address & 0x13) == 0x10 {
			index := address & ^uint16(0x10)
			ppu.palettes[index&0x1F] = data
		}
	}
}

func (ppu *PPU) writeDMA(data byte) {
	console := ppu.Console
	address := uint16(data) << 8
	for i := 0; i < 256; i++ {
		ppu.writeOAMData(console.FetchData(address))
		address++
	}
	console.CPU.Cycle += 513
	if console.CPU.Cycle%2 != 0 {
		console.CPU.Cycle++
	}
}

func (ppu *PPU) clearOAM() {
	for i := 0; i < 8; i++ {
		ppu.copyOam[i].id = 64
		ppu.copyOam[i].x = 0xFF
		ppu.copyOam[i].y = 0xFF
		ppu.copyOam[i].tile = 0xFF
		ppu.copyOam[i].attributes = 0xFF
		ppu.copyOam[i].low = 0
		ppu.copyOam[i].high = 0
	}
}

func (ppu *PPU) evalSprites() {
	n := 0
	var line int
	spriteHeight := int(ppu.Ctrl.spriteSize)
	var i byte
	for i = 0; i < 64; i++ {
		if ppu.ScanLine == 261 {
			line = -1
		} else {
			line = ppu.ScanLine
		}
		line -= int(ppu.oamMem[i*4])
		// If there's a sprite in the scanline, copy its properties into OAM copy:
		if line >= 0 && line < spriteHeight {
			ppu.copyOam[n].id = i
			ppu.copyOam[n].y = ppu.oamMem[i*4]
			ppu.copyOam[n].tile = ppu.oamMem[i*4+1]
			ppu.copyOam[n].attributes = ppu.oamMem[i*4+2]
			ppu.copyOam[n].x = ppu.oamMem[i*4+3]
			n++
			if n > 8 {
				ppu.Status.spriteOverflow = true
				break
			}
		}
	}
}

func (ppu *PPU) loadSprites() {
	var address uint16
	var tile uint16
	var spriteHeight uint16
	var spriteY uint16
	for i := 0; i < 8; i++ {
		ppu.oam[i] = ppu.copyOam[i]
		tile = uint16(ppu.oam[i].tile)
		spriteHeight = uint16(ppu.Ctrl.spriteSize)
		if ppu.Ctrl.spriteSize == 16 {
			address = ((tile & 1) * 0x1000) + ((tile & (^uint16(1))) * 16)
		} else {
			address = (ppu.Ctrl.spritePatternTable * 0x1000) + (tile * 16)
		}
		spriteY = (uint16(ppu.ScanLine) - uint16(ppu.oam[i].y)) % spriteHeight
		if ppu.oam[i].attributes&0x80 != 0 {
			spriteY ^= spriteHeight - 1
		}
		address += spriteY + (spriteY & 8)

		ppu.oam[i].low = ppu.access(address)
		ppu.oam[i].high = ppu.access(address + 8)
	}
}
