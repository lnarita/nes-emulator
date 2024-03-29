package processor

import "image/color"

func toColor(c uint, alpha byte) color.RGBA {
	return color.RGBA{R: byte(c >> 16), G: byte(c >> 8), B: byte(c), A: alpha}
}

var Colors = []color.RGBA{
	toColor(0x666666, 0xFF), toColor(0x002A88, 0xFF), toColor(0x1412A7, 0xFF), toColor(0x3B00A4, 0xFF), toColor(0x5C007E, 0xFF), toColor(0x6E0040, 0xFF), toColor(0x6C0600, 0xFF), toColor(0x561D00, 0xFF),
	toColor(0x333500, 0xFF), toColor(0x0B4800, 0xFF), toColor(0x005200, 0xFF), toColor(0x004F08, 0xFF), toColor(0x00404D, 0xFF), toColor(0x000000, 0xFF), toColor(0x000000, 0xFF), toColor(0x000000, 0xFF),
	toColor(0xADADAD, 0xFF), toColor(0x155FD9, 0xFF), toColor(0x4240FF, 0xFF), toColor(0x7527FE, 0xFF), toColor(0xA01ACC, 0xFF), toColor(0xB71E7B, 0xFF), toColor(0xB53120, 0xFF), toColor(0x994E00, 0xFF),
	toColor(0x6B6D00, 0xFF), toColor(0x388700, 0xFF), toColor(0x0C9300, 0xFF), toColor(0x008F32, 0xFF), toColor(0x007C8D, 0xFF), toColor(0x000000, 0xFF), toColor(0x000000, 0xFF), toColor(0x000000, 0xFF),
	toColor(0xFFFEFF, 0xFF), toColor(0x64B0FF, 0xFF), toColor(0x9290FF, 0xFF), toColor(0xC676FF, 0xFF), toColor(0xF36AFF, 0xFF), toColor(0xFE6ECC, 0xFF), toColor(0xFE8170, 0xFF), toColor(0xEA9E22, 0xFF),
	toColor(0xBCBE00, 0xFF), toColor(0x88D800, 0xFF), toColor(0x5CE430, 0xFF), toColor(0x45E082, 0xFF), toColor(0x48CDDE, 0xFF), toColor(0x4F4F4F, 0xFF), toColor(0x000000, 0xFF), toColor(0x000000, 0xFF),
	toColor(0xFFFEFF, 0xFF), toColor(0xC0DFFF, 0xFF), toColor(0xD3D2FF, 0xFF), toColor(0xE8C8FF, 0xFF), toColor(0xFBC2FF, 0xFF), toColor(0xFEC4EA, 0xFF), toColor(0xFECCC5, 0xFF), toColor(0xF7D8A5, 0xFF),
	toColor(0xE4E594, 0xFF), toColor(0xCFEF96, 0xFF), toColor(0xBDF4AB, 0xFF), toColor(0xB3F3CC, 0xFF), toColor(0xB5EBF2, 0xFF), toColor(0xB8B8B8, 0xFF), toColor(0x000000, 0xFF), toColor(0x000000, 0xFF),
}
