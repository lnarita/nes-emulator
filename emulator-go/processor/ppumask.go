package processor

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
