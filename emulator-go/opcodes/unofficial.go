package opcodes

import "students.ic.unicamp.br/goten/processor"

type ign struct{}

func (o ign) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o ign) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x0C, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x04, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x44, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x64, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x14, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x34, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x54, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x74, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0xD4, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0xF4, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x1C, addressingMode: processor.AbsoluteX, cycles: 4},
		Variation{opcode: 0x3C, addressingMode: processor.AbsoluteX, cycles: 4},
		Variation{opcode: 0x5C, addressingMode: processor.AbsoluteX, cycles: 4},
		Variation{opcode: 0x7C, addressingMode: processor.AbsoluteX, cycles: 4},
		Variation{opcode: 0xDC, addressingMode: processor.AbsoluteX, cycles: 4},
		Variation{opcode: 0xFC, addressingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o ign) GetName() string {
	return "*NOP"
}

type skb struct{}

func (o skb) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o skb) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x80, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x82, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0x89, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xC2, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xE2, addressingMode: processor.Immediate, cycles: 2},
	}
}

func (o skb) GetName() string {
	return "*NOP"
}

type unofficialNop struct{}

func (o unofficialNop) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o unofficialNop) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x1A, addressingMode: nil, cycles: 2},
		Variation{opcode: 0x3A, addressingMode: nil, cycles: 2},
		Variation{opcode: 0x5A, addressingMode: nil, cycles: 2},
		Variation{opcode: 0x7A, addressingMode: nil, cycles: 2},
		Variation{opcode: 0xDA, addressingMode: nil, cycles: 2},
		Variation{opcode: 0xFA, addressingMode: nil, cycles: 2},
	}
}

func (o unofficialNop) GetName() string {
	return "*NOP"
}

type lax struct{}

func (o lax) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o lax) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xA3, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0xA7, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xAF, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0xB3, addressingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0xB7, addressingMode: processor.ZeroPageY, cycles: 4},
		Variation{opcode: 0xBF, addressingMode: processor.AbsoluteY, cycles: 4},
	}
}

func (o lax) GetName() string {
	return "*LAX"
}

type sax struct{}

func (o sax) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o sax) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x83, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x87, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x8F, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x97, addressingMode: processor.ZeroPageY, cycles: 4},
	}
}

func (o sax) GetName() string {
	return "*SAX"
}

type unofficialSbc struct{}

func (o unofficialSbc) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o unofficialSbc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xEB, addressingMode: processor.Immediate, cycles: 2},
	}
}

func (o unofficialSbc) GetName() string {
	return "*SBC"
}

type dcp struct{}

func (o dcp) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o dcp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xC3, addressingMode: processor.IndirectX, cycles: 8},
		Variation{opcode: 0xC7, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0xCF, addressingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0xD3, addressingMode: processor.IndirectY, cycles: 8},
		Variation{opcode: 0xD7, addressingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0xDB, addressingMode: processor.AbsoluteY, cycles: 7},
		Variation{opcode: 0xDF, addressingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o dcp) GetName() string {
	return "*DCP"
}

type isc struct{}

func (o isc) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o isc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xE3, addressingMode: processor.IndirectX, cycles: 8},
		Variation{opcode: 0xE7, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0xEF, addressingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0xF3, addressingMode: processor.IndirectY, cycles: 8},
		Variation{opcode: 0xF7, addressingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0xFB, addressingMode: processor.AbsoluteY, cycles: 7},
		Variation{opcode: 0xFF, addressingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o isc) GetName() string {
	return "*ISB"
}

type slo struct{}

func (o slo) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o slo) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x03, addressingMode: processor.IndirectX, cycles: 8},
		Variation{opcode: 0x07, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x0F, addressingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0x13, addressingMode: processor.IndirectY, cycles: 8},
		Variation{opcode: 0x17, addressingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0x1B, addressingMode: processor.AbsoluteY, cycles: 7},
		Variation{opcode: 0x1F, addressingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o slo) GetName() string {
	return "*SLO"
}

type rla struct{}

func (o rla) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o rla) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x23, addressingMode: processor.IndirectX, cycles: 8},
		Variation{opcode: 0x27, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x2F, addressingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0x33, addressingMode: processor.IndirectY, cycles: 8},
		Variation{opcode: 0x37, addressingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0x3B, addressingMode: processor.AbsoluteY, cycles: 7},
		Variation{opcode: 0x3F, addressingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o rla) GetName() string {
	return "*RLA"
}

type sre struct{}

func (o sre) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o sre) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x43, addressingMode: processor.IndirectX, cycles: 8},
		Variation{opcode: 0x47, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x4F, addressingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0x53, addressingMode: processor.IndirectY, cycles: 8},
		Variation{opcode: 0x57, addressingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0x5B, addressingMode: processor.AbsoluteY, cycles: 7},
		Variation{opcode: 0x5F, addressingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o sre) GetName() string {
	return "*SRE"
}

type rra struct{}

func (o rra) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o rra) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x63, addressingMode: processor.IndirectX, cycles: 8},
		Variation{opcode: 0x67, addressingMode: processor.ZeroPage, cycles: 5},
		Variation{opcode: 0x6F, addressingMode: processor.Absolute, cycles: 6},
		Variation{opcode: 0x73, addressingMode: processor.IndirectY, cycles: 8},
		Variation{opcode: 0x77, addressingMode: processor.ZeroPageX, cycles: 6},
		Variation{opcode: 0x7B, addressingMode: processor.AbsoluteY, cycles: 7},
		Variation{opcode: 0x7F, addressingMode: processor.AbsoluteX, cycles: 7},
	}
}

func (o rra) GetName() string {
	return "*RRA"
}

var UnofficialOpcodes = []OpCode{
	ign{},
	unofficialNop{},
	skb{},
	lax{},
	sax{},
	unofficialSbc{},
	dcp{},
	isc{},
	slo{},
	rla{},
	sre{},
	rra{},
}
