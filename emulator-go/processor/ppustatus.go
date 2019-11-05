package processor

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
	status.spriteOverflow = flag&0b0010_0000 != 0
	status.spriteHit = flag&0b0100_0000 != 0
	status.vBlank = flag&0b1000_0000 != 0
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
