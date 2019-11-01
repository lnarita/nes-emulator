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
	odd                  byte // even/odd frame flag

	scanlineAddr uint16
	// Background latches:
	nt  uint16
	at  uint16
	bgL uint16
	bgH uint16
	// Background shift registers:
	atShiftL uint16
	atShiftH uint16
	bgShiftL uint16
	bgShiftH uint16
	atLatchL uint16
	atLatchH uint16

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
	if (flag>>2)&1 > 0 {
		increment = 32
	} else {
		increment = 1
	}
	var sprite byte
	if (flag>>5)&1 > 0 {
		sprite = 16
	} else {
		sprite = 8
	}
	ctrl.nameTable = uint16((flag>>0)&3) * 0x100
	ctrl.addressIncrement = increment
	ctrl.spritePatternTable = uint16((flag>>3)&1) * 0x1000
	ctrl.bgPatternTable = uint16((flag>>4)&1) * 0x1000
	ctrl.spriteSize = sprite
	ctrl.slave = (flag>>6)&1 != 0
	ctrl.nmiEnabled = (flag>>7)&1 != 0
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
	msk.grayscale = (flag>>0)&1 != 0
	msk.bgLeft = (flag>>1)&1 != 0
	msk.sprLeft = (flag>>2)&1 != 0
	msk.showBg = (flag>>3)&1 != 0
	msk.showSpr = (flag>>4)&1 != 0
	msk.intensifyRed = (flag>>5)&1 != 0
	msk.intensifyGreen = (flag>>6)&1 != 0
	msk.intensifyBlue = (flag>>7)&1 != 0
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
	if ppu.Cycle > 340 {
		ppu.Cycle %= 341
		ppu.ScanLine++
		if ppu.ScanLine > 261 {
			ppu.ScanLine = 0
			ppu.odd ^= 1
		}
	}
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
		case 2 < ppu.Cycle && ppu.Cycle < 255 || 322 < ppu.Cycle && ppu.Cycle < 337:
			ppu.renderPixel(render)
			switch ppu.Cycle % 8 {
			// Nametable:
			case 1:
				ppu.scanlineAddr = ppu.ntAddress()
				ppu.reloadShift()
			case 2:
				ppu.nt = uint16(ppu.access(ppu.scanlineAddr))
			// Attribute:
			case 3:
				ppu.scanlineAddr = ppu.atAddress()
			case 4:
				ppu.at = uint16(ppu.access(ppu.scanlineAddr))
				shift := ((ppu.vRamAddress >> 4) & 4) | (ppu.vRamAddress & 2)
				ppu.at = ((ppu.at >> shift) & 3) << 2
			// Background (low bits):
			case 5:
				ppu.scanlineAddr = ppu.bgAddress()
			case 6:
				ppu.bgL = uint16(ppu.access(ppu.scanlineAddr))
			// Background (high bits):
			case 7:
				ppu.scanlineAddr += 8
			case 0:
				ppu.bgH = uint16(ppu.access(ppu.scanlineAddr))
				ppu.hScroll(render)
			}
		case ppu.Cycle == 256:
			ppu.renderPixel(render)
			ppu.bgH = uint16(ppu.access(ppu.scanlineAddr))
			ppu.vScroll(render)
			// Vertical bump.
		case ppu.Cycle == 257:
			ppu.renderPixel(render)
			ppu.reloadShift()
			ppu.hUpdate(render)
			// Update horizontal position.
		case 280 < ppu.Cycle && ppu.Cycle < 304:
			if pre {
				ppu.vUpdate(render)
			}
			// Update vertical position.

		// No shift reloading:
		case ppu.Cycle == 1:
			ppu.scanlineAddr = ppu.ntAddress()
			if pre {
				ppu.Status.vBlank = false
			}
		case ppu.Cycle == 321 || ppu.Cycle == 339:
			ppu.scanlineAddr = ppu.ntAddress()
		// Nametable fetch instead of attribute:
		case ppu.Cycle == 338:
			ppu.scanlineAddr = uint16(ppu.access(ppu.scanlineAddr))
		case ppu.Cycle == 340:
			ppu.nt = uint16(ppu.access(ppu.scanlineAddr))
			if pre && render && ppu.odd == 1 {
				ppu.Cycle++
			}
		}
	}
}

func (ppu *PPU) ntAddress() uint16 {
	return 0x2000 | ppu.vRamAddress&0xFFF
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
	}
	ppu.writeCtrl(0)
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
			if n > 7 {
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
			address = ppu.Ctrl.spritePatternTable + (tile * 16)
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

func (ppu *PPU) renderPixel(render bool) {
	var palette byte = 0
	var objPalette byte = 0
	objPriority := false
	x := ppu.Cycle - 2

	if ppu.ScanLine < 240 && x >= 0 && x < 256 {
		if ppu.Mask.showBg && !(!ppu.Mask.bgLeft && x < 8) {
			// Background:
			palette = byte((nthBit(ppu.bgShiftH, 15-uint16(ppu.x)) << 1) | nthBit(ppu.bgShiftL, 15-uint16(ppu.x)))
			if palette != 0 {
				palette |= byte((nthBit(ppu.atShiftH, 7-uint16(ppu.x))<<1)|nthBit(ppu.atShiftL, 7-uint16(ppu.x))) << 2
			}
		}
		// Sprites:
		if ppu.Mask.showSpr && !(!ppu.Mask.sprLeft && x < 8) {
			for i := 7; i >= 0; i-- {
				if ppu.oam[i].id == 64 {
					continue // Void entry.
				}
				sprX := byte(x) - ppu.oam[i].x
				if sprX >= 8 {
					continue // Not in range.
				}
				if (ppu.oam[i].attributes & 0x40) != 0 {
					sprX ^= 7 // Horizontal flip.
				}

				sprPalette := (nthBitByte(ppu.oam[i].high, 7-sprX) << 1) | nthBitByte(ppu.oam[i].low, 7-sprX)
				if sprPalette == 0 {
					continue // Transparent pixel.
				}
				if ppu.oam[i].id == 0 && palette != 0 && x != 255 {
					ppu.Status.spriteHit = true
				}
				sprPalette |= (ppu.oam[i].attributes & 3) << 2
				objPalette = sprPalette + 16
				objPriority = (ppu.oam[i].attributes & 0x20) != 0
			}
		}
		// Evaluate priority:
		if objPalette != 0 && (palette == 0 || !objPriority) {
			palette = objPalette
		}
		var offset uint16
		if render {
			offset = uint16(palette)
		} else {
			offset = 0
		}
		ppu.Pixels[ppu.ScanLine*256+x] = Colors[ppu.access(0x3F00+offset)]
	}
	// Perform background shifts:
	ppu.bgShiftL <<= 1
	ppu.bgShiftH <<= 1
	ppu.atShiftL = (ppu.atShiftL << 1) | ppu.atLatchL
	ppu.atShiftH = (ppu.atShiftH << 1) | ppu.atLatchH
}

func nthBit(x uint16, n uint16) uint16 {
	return ((x) >> (n)) & 1
}

func nthBitByte(x byte, n byte) byte {
	return ((x) >> (n)) & 1
}

func (ppu *PPU) reloadShift() {
	ppu.bgShiftL = (ppu.bgShiftL & 0xFF00) | ppu.bgL
	ppu.bgShiftH = (ppu.bgShiftH & 0xFF00) | ppu.bgH
	ppu.atLatchL = ppu.at & 1
	ppu.atLatchH = ppu.at & 2
}

func (ppu *PPU) atAddress() uint16 {
	return 0x23C0 | (ppu.vRamAddress & 0x0C00) | ((ppu.vRamAddress >> 4) & 0x38) | ((ppu.vRamAddress >> 2) & 0x07)
}

func (ppu *PPU) bgAddress() uint16 {
	fineY := (ppu.vRamAddress >> 12) & 7
	return ppu.Ctrl.bgPatternTable + (ppu.nt * 16) + fineY
}

func (ppu *PPU) hScroll(render bool) {
	if !render {
		return
	}
	if ppu.vRamAddress&0x001F == 31 {
		// coarse X = 0
		ppu.vRamAddress &= 0xFFE0
		// switch horizontal nametable
		ppu.vRamAddress ^= 0x0400
	} else {
		// increment coarse X
		ppu.vRamAddress++
	}
}

func (ppu *PPU) vScroll(render bool) {
	if !render {
		return
	}
	// increment vert(v)
	// if fine Y < 7
	if ppu.vRamAddress&0x7000 != 0x7000 {
		// increment fine Y
		ppu.vRamAddress += 0x1000
	} else {
		// fine Y = 0
		ppu.vRamAddress &= 0x8FFF
		// let y = coarse Y
		y := (ppu.vRamAddress & 0x03E0) >> 5
		if y == 29 {
			// coarse Y = 0
			y = 0
			// switch vertical nametable
			ppu.vRamAddress ^= 0x0800
		} else if y == 31 {
			// coarse Y = 0, nametable not switched
			y = 0
		} else {
			// increment coarse Y
			y++
		}
		// put coarse Y back into v
		ppu.vRamAddress = (ppu.vRamAddress & 0xFC1F) | (y << 5)
	}
}

func (ppu *PPU) hUpdate(render bool) {
	if !render {
		return
	}
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	ppu.vRamAddress = (ppu.vRamAddress & 0xFBE0) | (ppu.temporaryVRamAddress & 0x041F)
}
func (ppu *PPU) vUpdate(render bool) {
	if !render {
		return
	}
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	ppu.vRamAddress = (ppu.vRamAddress & 0x841F) | (ppu.temporaryVRamAddress & 0x7BE0)
}
