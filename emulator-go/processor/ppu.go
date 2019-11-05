package processor

import (
	"fmt"
	"image/color"
)

type PPU struct {
	Console *Console // parent pointer :/

	Cycle    int // 0-340
	ScanLine int
	Frame    uint64

	Ctrl   *PPUCTRL
	Mask   *PPUMASK
	Status *PPUSTATUS

	palettes   [32]byte // 4 colours each; 4 background palettes; 4 foreground palettes
	nametables [2048]byte
	oam        *OAM

	spriteCount byte
	sprites     [8]Sprite
	Pixels      [256 * 240]color.RGBA

	// PPUSCROLL registers
	v     uint16
	t     uint16
	x     byte // fine x scroll (3 bit)
	latch bool // write latch
	odd   bool // even/odd frame flag

	// NMI flags
	nmiOccurred bool
	nmiOutput   bool
	nmiPrevious bool
	nmiDelay    byte

	nt       uint8
	at       uint8
	tileLow  uint8
	tileHigh uint8
	tileData uint64

	readAddress    uint16
	attributeShift uint16
	oamAddress     byte
	bufferedData   byte
}

type Sprite struct {
	id         uint8 // index in sprites
	x          uint8
	y          uint8
	tile       uint32
	attributes uint8
}

func (ppu PPU) String() string {
	return fmt.Sprintf("PPU { Cycle: %d, V: %04X, T: %04X, NT: %04X, AT: %04X, tileData: %04X, tileLow: %04X, tileHigh: %04X, Palettes: % 04X }", ppu.Cycle, ppu.v, ppu.t, ppu.nt, ppu.at, ppu.tileData, ppu.tileLow, ppu.tileHigh, ppu.palettes)
}

func (ppu *PPU) Step() {
	ppu.Tick()
	ppu.Tick()
	ppu.Tick()
}

func (ppu *PPU) Tick() {
	if ppu.nmiDelay > 0 {
		ppu.nmiDelay--
		if ppu.nmiDelay == 0 && ppu.nmiOutput && ppu.nmiOccurred {
			ppu.Console.TriggerNMI()
		}
	}

	if ppu.Mask.showBg || ppu.Mask.showSpr {
		if ppu.odd && ppu.ScanLine == 261 && ppu.Cycle == 339 {
			ppu.Cycle = 0
			ppu.ScanLine = 0
			ppu.Frame++
			ppu.odd = !ppu.odd
			return
		}
	}

	ppu.Cycle++
	if ppu.Cycle > 340 {
		ppu.Cycle = 0
		ppu.ScanLine++
		if ppu.ScanLine > 261 {
			ppu.ScanLine = 0
			ppu.Frame++
			ppu.odd = !ppu.odd
		}
	}

	renderingEnabled := ppu.Mask.showBg || ppu.Mask.showSpr

	// There are a total of 262 scanlines per frame
	//   Scanlines 0 to 239 are for display (NES is 256 x 240)
	//   Scanline  240 is a post-render scanline (idle)
	//   Scanlines 241 to 260 are the vblank interval
	//   Scanline  261 is a pre-render scanline
	visibleLine := ppu.ScanLine < 240
	preLine := ppu.ScanLine == 261

	ppu.scanlineCycle(renderingEnabled, visibleLine, preLine)
}

func (ppu *PPU) scanlineCycle(render, visible, pre bool) {
	renderLine := visible || pre
	preFetchCycle := ppu.Cycle >= 321 && ppu.Cycle <= 336
	visibleCycle := ppu.Cycle >= 1 && ppu.Cycle <= 256
	fetchCycle := preFetchCycle || visibleCycle

	// BG
	if visible && visibleCycle {
		ppu.renderPixel(render)
	}

	if render {
		if renderLine && fetchCycle {
			ppu.tileData <<= 4

			switch ppu.Cycle % 8 {
			case 1:
				ppu.readAddress = ppu.ntAddress()
			case 2:
				ppu.nt = ppu.rd(ppu.readAddress)
			case 3:
				ppu.readAddress, ppu.attributeShift = ppu.atAddress()
			case 4:
				ppu.at = ((ppu.rd(ppu.readAddress) >> ppu.attributeShift) & 3) << 2
			case 5:
				ppu.readAddress = ppu.bgAddress()
			case 6:
				ppu.tileLow = ppu.rd(ppu.readAddress)
			case 7:
				ppu.readAddress += 8
			case 0:
				ppu.tileHigh = ppu.rd(ppu.readAddress)
				ppu.storeTileData()
			}
		}
	}

	if pre && ppu.Cycle >= 280 && ppu.Cycle <= 304 {
		ppu.vUpdate(render)
	}

	if renderLine {
		if fetchCycle && ppu.Cycle%8 == 0 {
			ppu.hScroll(render)
		}
		if ppu.Cycle == 256 {
			ppu.vScroll(render)
		}
		if ppu.Cycle == 257 {
			ppu.hUpdate(render)
		}
	}

	// Sprite
	if render && ppu.Cycle == 257 {
		if visible {
			ppu.evalSprites()
		} else {
			ppu.spriteCount = 0
		}
	}

	// vblank logic
	if ppu.ScanLine == 241 && ppu.Cycle == 1 {
		ppu.setVerticalBlank()
	}
	if pre && ppu.Cycle == 1 {
		ppu.clearVerticalBlank()
		ppu.Status.spriteHit = false
		ppu.Status.spriteOverflow = false
	}
}

func (ppu *PPU) ntAddress() uint16 {
	return 0x2000 | ppu.v&0xFFF
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
	}
	ppu.writeMask(0)
	if ppu.Status == nil {
		ppu.Status = createStatusFromInt(0)
	}
	ppu.writeOAMAddress(0)
	if ppu.oam == nil {
		ppu.oam = &OAM{}
	}
}

// OK?
func (ppu *PPU) Read(address uint16) byte {
	var result byte
	switch address {
	case 0x2002:
		result = ppu.readStatus()
	case 0x2004:
		result = ppu.readOAM()
	case 0x2007:
		result = ppu.readData()
	}
	return result
}

// OK?
func (ppu *PPU) ReadDataForLog(address uint16) byte {
	var result byte
	switch address {
	case 0x2007:
		result = ppu.readDataForLog()
	}
	return result
}

// OK?
func (ppu *PPU) readOAM() byte {
	res := ppu.oam.Read(uint16(ppu.oamAddress))
	return res
}

// CHECK FLAGS
func (ppu *PPU) readStatus() byte {
	ppu.Status.vBlank = ppu.nmiOccurred
	res := ppu.Status.toFlag()
	ppu.latch = false
	ppu.nmiOccurred = false
	ppu.nmiChange()
	return res
}

// OK?
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

func (ppu *PPU) rd(address uint16) byte {
	var result byte
	console := ppu.Console
	index := address % 0x4000
	switch {
	case index < 0x2000:
		result = console.Memory.mapper.Read(index)
	case index < 0x3F00:
		mirrorAddress := console.Cartridge.Mirroring.calcAddress(address) % 2048
		result = ppu.nametables[mirrorAddress]
	case index < 0x4000:
		result = ppu.readPalette(index % 0x20)
	}
	return result
}

// OK?
func (ppu *PPU) readData() byte {
	var res byte
	if ppu.v%0x4000 < 0x3F00 {
		res = ppu.bufferedData
		ppu.bufferedData = ppu.rd(ppu.v)
	} else {
		ppu.bufferedData = ppu.rd(ppu.v - 0x1000)
		res = ppu.rd(ppu.v)
	}
	ppu.v += ppu.Ctrl.addressIncrement
	return res
}

// OK?
func (ppu *PPU) readDataForLog() byte {
	var res byte
	if ppu.v%0x4000 < 0x3F00 {
		res = ppu.bufferedData
	} else {
		res = ppu.rd(ppu.v)
	}
	return res
}

func (ppu *PPU) updateStatusBus(data byte) {
	ppu.Status.bus = data & 0x1F
}

// OK?
func (ppu *PPU) writeCtrl(data byte) {
	ppu.Ctrl.fromFlag(data)
	ppu.nmiOutput = (data>>7)&1 == 1
	ppu.nmiChange()
	ppu.t = (ppu.t & 0xF3FF) | (uint16(data) & 0x03 << 10)
}

// OK?
func (ppu *PPU) writeMask(data byte) {
	ppu.Mask.fromFlag(data)
}

// OK?
func (ppu *PPU) writeOAMAddress(data byte) {
	ppu.oamAddress = data
}

// OK?
func (ppu *PPU) writeOAMData(data byte) {
	ppu.oam.Write(uint16(ppu.oamAddress), data)
	ppu.oamAddress++
}

// OK?
func (ppu *PPU) writeScroll(data byte) {
	if !ppu.latch {
		ppu.t = (ppu.t & 0xFFE0) | (uint16(data) >> 3)
		ppu.x = data & 0x07
	} else {
		ppu.t = (ppu.t & 0x8FFF) | ((uint16(data) & 0x07) << 12)
		ppu.t = (ppu.t & 0xFC1F) | ((uint16(data) & 0xF8) << 2)
	}
	ppu.latch = !ppu.latch
}

// OK?
func (ppu *PPU) writeAddress(data byte) {
	if !ppu.latch {
		ppu.t = (ppu.t & 0x80FF) | ((uint16(data) & 0x3F) << 8)
	} else {
		ppu.t = (ppu.t & 0xFF00) | uint16(data)
		ppu.v = ppu.t
	}
	ppu.latch = !ppu.latch
}

// OK?
func (ppu *PPU) writeData(data byte) {
	ppu.wr(ppu.v, data)
	ppu.v += ppu.Ctrl.addressIncrement
}

func (ppu *PPU) wr(address uint16, data byte) {
	console := ppu.Console
	index := address % 0x4000
	switch {
	case index < 0x2000:
		console.Memory.mapper.Write(index, data)
	case index < 0x3F00:
		mirrorAddress := console.Cartridge.Mirroring.calcAddress(index) % 2048
		ppu.nametables[mirrorAddress] = data
	case index < 0x4000:
		ppu.writePalette(index%32, data)
	}
}

func (ppu *PPU) writeDMA(data byte) {
	console := ppu.Console
	address := uint16(data) << 8
	for i := 0; i < 256; i++ {
		ppu.writeOAMData(console.FetchData(address))
		address++
	}
	console.CPU.Stall += 513
	if console.CPU.Cycle%2 != 0 {
		console.CPU.Stall++
	}
}

func (ppu *PPU) evalSprites() {
	count := 0
	var line int
	var y uint8
	size := int(ppu.Ctrl.spriteSize)
	var baseAddress uint16
	var i uint16
	for i = 0; i < 64; i++ {
		baseAddress = i * 4
		y = ppu.oam.Read(baseAddress)
		line = ppu.ScanLine - int(y)
		if line < 0 || line >= size {
			continue
		}
		ppu.sprites[count].id = byte(i)
		ppu.sprites[count].attributes = ppu.oam.Read(baseAddress + 2)
		ppu.sprites[count].tile = ppu.fetchSpritePattern(count, i, line, ppu.sprites[count].attributes)
		ppu.sprites[count].x = ppu.oam.Read(baseAddress + 3)
		ppu.sprites[count].y = y
		count++
		if count > 7 {
			ppu.Status.spriteOverflow = true
			break
		}
	}
	ppu.spriteCount = byte(count)
}

func (ppu *PPU) fetchSpritePattern(id int, index uint16, line int, attributes byte) uint32 {
	tile := uint16(ppu.oam.Read(index*4 + 1))
	var address uint16
	if ppu.Ctrl.spriteSize == 8 {
		if attributes&0x80 == 0x80 {
			line = 7 - line
		}
		address = ppu.Ctrl.spritePatternTable + (tile * 16) + uint16(line)
	} else {
		if attributes&0x80 == 0x80 {
			line = 15 - line
		}
		table := tile & 1
		tile &= 0xFE
		if line > 7 {
			tile++
			line -= 8
		}

		address = 0x1000*table + (tile * 16) + uint16(line)
	}

	a := (attributes & 3) << 2
	lowTileByte := ppu.rd(address)
	highTileByte := ppu.rd(address + 8)
	var data uint32
	for i := 0; i < 8; i++ {
		var p1, p2 byte
		if attributes&0x40 == 0x40 {
			p1 = (lowTileByte & 1) << 0
			p2 = (highTileByte & 1) << 1
			lowTileByte >>= 1
			highTileByte >>= 1
		} else {
			p1 = (lowTileByte & 0x80) >> 7
			p2 = (highTileByte & 0x80) >> 6
			lowTileByte <<= 1
			highTileByte <<= 1
		}
		data <<= 4
		data |= uint32(a | p1 | p2)
	}
	return data
}

func (ppu *PPU) atAddress() (uint16, uint16) {
	v := ppu.v
	address := 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07)
	shift := ((v >> 4) & 4) | (v & 2)
	return address, shift
}
func (ppu *PPU) bgAddress() uint16 {
	fineY := (ppu.v >> 12) & 7
	tile := uint16(ppu.nt)
	return ppu.Ctrl.bgPatternTable + (tile * 16) + fineY
}
func (ppu *PPU) hScroll(render bool) {
	if !render {
		return
	}
	if ppu.v&0x001F == 31 {
		// coarse X = 0
		ppu.v &= 0xFFE0
		// switch horizontal nametable
		ppu.v ^= 0x0400
	} else {
		// increment coarse X
		ppu.v++
	}
}

func (ppu *PPU) vScroll(render bool) {
	if !render {
		return
	}
	// increment vert(v)
	// if fine Y < 7
	if ppu.v&0x7000 != 0x7000 {
		// increment fine Y
		ppu.v += 0x1000
	} else {
		// fine Y = 0
		ppu.v &= 0x8FFF
		// let y = coarse Y
		y := (ppu.v & 0x03E0) >> 5
		if y == 29 {
			// coarse Y = 0
			y = 0
			// switch vertical nametable
			ppu.v ^= 0x0800
		} else if y == 31 {
			// coarse Y = 0, nametable not switched
			y = 0
		} else {
			// increment coarse Y
			y++
		}
		// put coarse Y back into v
		ppu.v = (ppu.v & 0xFC1F) | (y << 5)
	}
}

func (ppu *PPU) hUpdate(render bool) {
	if !render {
		return
	}
	// v: .....F.. ...EDCBA = t: .....F.. ...EDCBA
	ppu.v = (ppu.v & 0xFBE0) | (ppu.t & 0x041F)
}
func (ppu *PPU) vUpdate(render bool) {
	if !render {
		return
	}
	// v: .IHGF.ED CBA..... = t: .IHGF.ED CBA.....
	ppu.v = (ppu.v & 0x841F) | (ppu.t & 0x7BE0)
}

func (ppu *PPU) writePalette(address uint16, data byte) {
	index := address
	if index >= 16 && index%4 == 0 {
		index -= 16
	}
	ppu.palettes[index] = data
}

func (ppu *PPU) readPalette(address uint16) byte {
	index := address
	if index >= 16 && index%4 == 0 {
		index -= 16
	}
	return ppu.palettes[index]
}

func (ppu *PPU) storeTileData() {
	var p1 byte
	var p2 byte
	var data uint32
	for i := 0; i < 8; i++ {
		p1 = (ppu.tileLow & 0x80) >> 7
		p2 = (ppu.tileHigh & 0x80) >> 6

		ppu.tileLow <<= 1
		ppu.tileHigh <<= 1

		data <<= 4
		data |= uint32(ppu.at | p1 | p2)
	}
	ppu.tileData |= uint64(data)
}

func (ppu *PPU) backgroundPixel() byte {
	if !ppu.Mask.showBg {
		return 0x0
	}

	tileData := ppu.fetchTileData() >> ((7 - ppu.x) * 4)
	colour := byte(tileData & 0x0F)
	return colour
}

func (ppu *PPU) spritePixel() (byte, byte) {
	if !ppu.Mask.showSpr {
		return 0, 0
	}

	var offset int16
	var colour byte
	var i byte
	for i = 0; i < ppu.spriteCount; i++ {
		offset = (int16(ppu.Cycle) - 1) - int16(ppu.sprites[i].x)

		if offset < 0 || offset > 7 {
			continue
		}

		offset = 7 - offset
		colour = byte((ppu.sprites[i].tile >> (offset * 4)) & 0x0F)
		// transparent colour
		if colour%4 == 0 {
			continue
		}

		return i, colour
	}
	return 0, 0
}

func (ppu *PPU) renderPixel(render bool) {
	if !render {
		return
	}
	x := ppu.Cycle - 1
	y := ppu.ScanLine

	background := ppu.backgroundPixel()
	i, sprite := ppu.spritePixel()

	if x < 8 && !ppu.Mask.bgLeft {
		background = 0
	}

	if x < 8 && !ppu.Mask.sprLeft {
		sprite = 0
	}

	b := background%4 != 0
	s := sprite%4 != 0

	var addressLowNyb uint16
	if !b && !s {
		addressLowNyb = 0
	} else if !b && s {
		addressLowNyb = uint16(sprite) | 0x10
	} else if b && !s {
		addressLowNyb = uint16(background)
	} else if b && s {
		if ppu.sprites[i].id == 0 && x < 255 {
			ppu.Status.spriteHit = true
		}

		if (ppu.sprites[i].attributes>>5)&1 == 0 {
			addressLowNyb = uint16(sprite) | 0x10
		} else {
			addressLowNyb = uint16(background)
		}
	}

	address := 0x3F00 | addressLowNyb
	paletteIndex := ppu.rd(address) % 64
	colour := Colors[paletteIndex]
	ppu.Pixels[(y*256+x)%len(ppu.Pixels)] = colour
}

func (ppu *PPU) fetchTileData() uint32 {
	return uint32(ppu.tileData >> 32)
}

func (ppu *PPU) nmiChange() {
	nmi := ppu.nmiOutput && ppu.nmiOccurred
	if nmi && !ppu.nmiPrevious {
		ppu.nmiDelay = 15
	}
	ppu.nmiPrevious = nmi
}

func (ppu *PPU) setVerticalBlank() {
	ppu.nmiOccurred = true
	ppu.nmiChange()
}

func (ppu *PPU) clearVerticalBlank() {
	ppu.nmiOccurred = false
	ppu.nmiChange()
}
