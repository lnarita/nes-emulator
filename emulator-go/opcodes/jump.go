package opcodes

import "students.ic.unicamp.br/goten/processor"

type blp struct{}

func (o blp) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o blp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x10, addresingMode: processor.Relative, cycles: 2},
	}
}

func (o blp) getName() string {
	return "BLP"
}

type bmi struct{}

func (o bmi) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o bmi) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x30, addresingMode: processor.Relative, cycles: 2},
	}
}

func (o bmi) getName() string {
	return "BMI"
}

type bvc struct{}

func (o bvc) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o bvc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x50, addresingMode: processor.Relative, cycles: 2},
	}
}

func (o bvc) getName() string {
	return "BVC"
}

type bvs struct{}

func (o bvs) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o bvs) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x70, addresingMode: processor.Relative, cycles: 2},
	}
}

func (o bvs) getName() string {
	return "BVS"
}

type bcc struct{}

func (o bcc) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o bcc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x90, addresingMode: processor.Relative, cycles: 2},
	}
}

func (o bcc) getName() string {
	return "bcc"
}

type bcs struct{}

func (o bcs) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o bcs) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xB0, addresingMode: processor.Relative, cycles: 2},
	}
}

func (o bcs) getName() string {
	return "bcs"
}

type bne struct{}

func (o bne) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o bne) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xD0, addresingMode: processor.Relative, cycles: 2},
	}
}

func (o bne) getName() string {
	return "BNE"
}

type beq struct{}

func (o beq) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o beq) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xF0, addresingMode: processor.Relative, cycles: 2},
	}
}

func (o beq) getName() string {
	return "BEQ"
}

type brk struct{}

func (o brk) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o brk) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x00, addresingMode: nil, cycles: 7},
	}
}

func (o brk) getName() string {
	return "BRK"
}

type rti struct{}

func (o rti) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o rti) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x00, addresingMode: nil, cycles: 6},
	}
}

func (o rti) getName() string {
	return "RTI"
}

type jsr struct{}

func (o jsr) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o jsr) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x20, addresingMode: processor.Absolute, cycles: 6},
	}
}

func (o jsr) getName() string {
	return "JSR"
}

type rts struct{}

func (o rts) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o rts) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x60, addresingMode: nil, cycles: 6},
	}
}

func (o rts) getName() string {
	return "JSR"
}

type jmp struct{}

func (o jmp) exec(cpu *processor.CPU, memory *processor.Memory, addressMode *processor.AddressMode) {

}

func (o jmp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x4C, addresingMode: processor.Absolute, cycles: 3},
		Variation{opcode: 0x6C, addresingMode: processor.Indirect, cycles: 5},
	}
}

func (o jmp) getName() string {
	return "JMP"
}

var JumpOpCodes []OpCode {
	bpl{},
	bmi{},
	bvc{},
	bvs{},
	bcc{},
	bcs{},
	bne{},
	beq{},
	brk{},
	rti{},
	jsr{},
	rts{},
	jmp{},
}