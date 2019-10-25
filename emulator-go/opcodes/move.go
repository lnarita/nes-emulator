package opcodes

import "students.ic.unicamp.br/goten/processor"

type lda struct{}

func (o lda) exec(console *processor.Console, variation *Variation) {

}

func (o lda) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xA1, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0xA5, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xA9, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xAD, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0xB1, addressingMode: processor.IndirectY, cycles: 5},
		Variation{opcode: 0xB5, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0xB9, addressingMode: processor.AbsoluteY, cycles: 4},
		Variation{opcode: 0xBd, addressingMode: processor.AbsoluteX, cycles: 4},
	}
}

func (o lda) getName() string {
	return "LDA"
}

type sta struct{}

func (o sta) exec(console *processor.Console, variation *Variation) {

}

func (o sta) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x81, addressingMode: processor.IndirectX, cycles: 6},
		Variation{opcode: 0x85, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x89, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x9D, addressingMode: processor.IndirectY, cycles: 6},
		Variation{opcode: 0x91, addressingMode: processor.ZeroPageX, cycles: 4},
		Variation{opcode: 0x95, addressingMode: processor.AbsoluteY, cycles: 5},
		Variation{opcode: 0x99, addressingMode: processor.AbsoluteX, cycles: 5},
	}
}

func (o sta) getName() string {
	return "STA"
}

type ldx struct{}

func (o ldx) exec(console *processor.Console, variation *Variation) {

}

func (o ldx) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xA2, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xA6, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xAE, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0xB6, addressingMode: processor.ZeroPageY, cycles: 4},
		Variation{opcode: 0xBE, addressingMode: processor.AbsoluteY, cycles: 4},
	}
}

func (o ldx) getName() string {
	return "LDX"
}

type stx struct{}

func (o stx) exec(console *processor.Console, variation *Variation) {

}

func (o stx) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x86, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0x8E, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x96, addressingMode: processor.ZeroPageY, cycles: 4},
	}
}

func (o stx) getName() string {
	return "STX"
}

type ldy struct{}

func (o ldy) exec(console *processor.Console, variation *Variation) {

}

func (o ldy) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xA0, addressingMode: processor.Immediate, cycles: 2},
		Variation{opcode: 0xA4, addressingMode: processor.ZeroPage, cycles: 3},
		Variation{opcode: 0xAC, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0xB4, addressingMode: processor.ZeroPageY, cycles: 4},
		Variation{opcode: 0xBC, addressingMode: processor.AbsoluteY, cycles: 4},
	}
}

func (o ldy) getName() string {
	return "LDY"
}

type sty struct{}

func (o sty) exec(console *processor.Console, variation *Variation) {

}

func (o sty) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x84, addressingMode: processor.ZeroPage, cycles: 2},
		Variation{opcode: 0x8C, addressingMode: processor.Absolute, cycles: 4},
		Variation{opcode: 0x94, addressingMode: processor.ZeroPageX, cycles: 4},
	}
}

func (o sty) getName() string {
	return "STY"
}

type tax struct{}

func (o tax) exec(console *processor.Console, variation *Variation) {

}

func (o tax) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xAA, addressingMode: nil, cycles: 2},
	}
}

func (o tax) getName() string {
	return "TAX"
}

type txa struct{}

func (o txa) exec(console *processor.Console, variation *Variation) {

}

func (o txa) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x8A, addressingMode: nil, cycles: 2},
	}
}

func (o txa) getName() string {
	return "TXA"
}

type tay struct{}

func (o tay) exec(console *processor.Console, variation *Variation) {

}

func (o tay) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xA8, addressingMode: nil, cycles: 2},
	}
}

func (o tay) getName() string {
	return "TAY"
}

type tya struct{}

func (o tya) exec(console *processor.Console, variation *Variation) {

}

func (o tya) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x98, addressingMode: nil, cycles: 2},
	}
}

func (o tya) getName() string {
	return "TYA"
}

type tsx struct{}

func (o tsx) exec(console *processor.Console, variation *Variation) {

}

func (o tsx) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xBA, addressingMode: nil, cycles: 2},
	}
}

func (o tsx) getName() string {
	return "TSX"
}

type txs struct{}

func (o txs) exec(console *processor.Console, variation *Variation) {

}

func (o txs) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x9A, addressingMode: nil, cycles: 2},
	}
}

func (o txs) getName() string {
	return "TXS"
}

type pla struct{}

func (o pla) exec(console *processor.Console, variation *Variation) {

}

func (o pla) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x68, addressingMode: nil, cycles: 4},
	}
}

func (o pla) getName() string {
	return "PLA"
}

type pha struct{}

func (o pha) exec(console *processor.Console, variation *Variation) {

}

func (o pha) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x48, addressingMode: nil, cycles: 3},
	}
}

func (o pha) getName() string {
	return "PHA"
}

type plp struct{}

func (o plp) exec(console *processor.Console, variation *Variation) {

}

func (o plp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x28, addressingMode: nil, cycles: 4},
	}
}

func (o plp) getName() string {
	return "plp"
}

type php struct{}

func (o php) exec(console *processor.Console, variation *Variation) {

}

func (o php) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x08, addressingMode: nil, cycles: 3},
	}
}

func (o php) getName() string {
	return "php"
}

var MoveOpCodes = []OpCode{
	lda{},
	sta{},
	ldx{},
	stx{},
	ldy{},
	sty{},
	tax{},
	txa{},
	tay{},
	tya{},
	tsx{},
	txs{},
	pla{},
	pha{},
	plp{},
	php{},
}
