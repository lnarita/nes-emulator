package opcodes

import "students.ic.unicamp.br/goten/processor"

type ora struct{}

func (o ora) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o ora) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x01, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x05, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x09, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x0D, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x11, addressingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0x15, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x19, addressingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0x1D, addressingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o ora) GetName() string {
	return "ORA"
}

type and struct{}

func (o and) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o and) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x21, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x25, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x29, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x2D, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x31, addressingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0x35, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x39, addressingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0x3D, addressingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o and) GetName() string {
	return "AND"
}

type eor struct{}

func (o eor) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o eor) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x41, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x45, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x49, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x4D, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x51, addressingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0x55, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x59, addressingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0x5D, addressingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o eor) GetName() string {
	return "EOR"
}

type adc struct{}

func (o adc) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o adc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x61, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x65, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x69, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x6D, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x71, addressingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0x75, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x79, addressingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0x7D, addressingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o adc) GetName() string {
	return "ADC"
}

type sbc struct{}

func (o sbc) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o sbc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xE1, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0xE5, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xE9, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xED, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0xF1, addressingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0xF5, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0xF9, addressingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0xFD, addressingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o sbc) GetName() string {
	return "SBC"
}

type cmp struct{}

func (o cmp) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o cmp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xC1, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0xC5, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xC9, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xCD, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0xD1, addressingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0xD5, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0xD9, addressingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0xDD, addressingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o cmp) GetName() string {
	return "CMP"
}

type cpx struct{}

func (o cpx) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o cpx) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xE0, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xE4, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xEC, addressingMode: processor.Absolute, cycles: 4},
	}
}

func (o cpx) GetName() string {
	return "CPX"
}

type cpy struct{}

func (o cpy) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o cpy) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xC0, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xC4, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xCC, addressingMode: processor.Absolute, cycles: 4},
	}
}

func (o cpy) GetName() string {
	return "CPY"
}

type dec struct{}

func (o dec) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o dec) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xC6, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0xCE, addressingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0xD6, addressingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0xDE, addressingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o dec) GetName() string {
	return "DEC"
}

type dex struct{}

func (o dex) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o dex) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xCA, addressingMode: nil, cycles: 2},
	}
}

func (o dex) GetName() string {
	return "DEX"
}

type dey struct{}

func (o dey) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o dey) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x88, addressingMode: nil, cycles: 2},
	}
}

func (o dey) GetName() string {
	return "DEY"
}

type inc struct{}

func (o inc) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o inc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xE6, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0xEE, addressingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0xF6, addressingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0xFE, addressingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o inc) GetName() string {
	return "INC"
}

type inx struct{}

func (o inx) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o inx) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xE8, addressingMode: nil, cycles: 2},
	}
}

func (o inx) GetName() string {
	return "INX"
}

type iny struct{}

func (o iny) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o iny) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xC8, addressingMode: nil, cycles: 2},
	}
}

func (o iny) GetName() string {
	return "INY"
}

type asl struct{}

func (o asl) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o asl) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x06, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x0A, addressingMode: processor.ZeroPage, cycles: 2},
		Variation{opcode: 0x0E, addressingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x16, addressingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x1E, addressingMode: processor.ZeroPage, cycles: 7},
	}
}

func (o asl) GetName() string {
	return "ASL"
}

type rol struct{}

func (o rol) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o rol) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x26, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x2A, addressingMode: processor.ZeroPage, cycles: 2},
		Variation{opcode: 0x2E, addressingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x36, addressingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x3E, addressingMode: processor.ZeroPage, cycles: 7},
	}
}

func (o rol) GetName() string {
	return "ROL"
}

type lsr struct{}

func (o lsr) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o lsr) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x46, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x4A, addressingMode: processor.ZeroPage, cycles: 2},
		Variation{opcode: 0x4E, addressingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x56, addressingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x5E, addressingMode: processor.ZeroPage, cycles: 7},
	}
}

func (o lsr) GetName() string {
	return "LSR"
}

type ror struct{}

func (o ror) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o ror) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x66, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x6A, addressingMode: processor.ZeroPage, cycles: 2},
		Variation{opcode: 0x6E, addressingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x76, addressingMode: processor.ZeroPage, cycles: 6},
		Variation{opcode: 0x7E, addressingMode: processor.ZeroPage, cycles: 7},
	}
}

func (o ror) GetName() string {
	return "ROR"
}

var ArithmeticAndLogicalOpCodes = []OpCode{
	ora{},
	and{},
	eor{},
	adc{},
	sbc{},
	cmp{},
	cpx{},
	cpy{},
	dec{},
	dex{},
	dey{},
	inc{},
	inx{},
	iny{},
	asl{},
	rol{},
	lsr{},
	ror{},
}
