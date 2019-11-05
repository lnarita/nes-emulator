package processor

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
	if (flag>>2)&1 == 0 {
		increment = 1
	} else {
		increment = 32
	}
	var sprite byte
	if (flag>>5)&1 == 0 {
		sprite = 8
	} else {
		sprite = 16
	}
	//ctrl.nameTable = uint16((flag>>0)&3) * 0x100
	ctrl.addressIncrement = increment
	ctrl.spritePatternTable = uint16((flag>>3)&1) * 0x1000
	ctrl.bgPatternTable = uint16((flag>>4)&1) * 0x1000
	ctrl.spriteSize = sprite
	ctrl.slave = (flag>>6)&1 != 0
	ctrl.nmiEnabled = (flag>>7)&1 != 0
}
