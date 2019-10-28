package opcodes

import "students.ic.unicamp.br/goten/processor"

type blp struct{}

func (o blp) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o blp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x10, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o blp) GetName() string {
	return "BLP"
}

type bmi struct{}

func (o bmi) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o bmi) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x30, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bmi) GetName() string {
	return "BMI"
}

type bvc struct{}

func (o bvc) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o bvc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x50, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bvc) GetName() string {
	return "BVC"
}

type bvs struct{}

func (o bvs) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o bvs) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x70, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bvs) GetName() string {
	return "BVS"
}

type bcc struct{}

func (o bcc) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o bcc) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x90, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bcc) GetName() string {
	return "bcc"
}

type bcs struct{}

func (o bcs) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o bcs) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xB0, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bcs) GetName() string {
	return "bcs"
}

type bne struct{}

func (o bne) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o bne) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xD0, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o bne) GetName() string {
	return "BNE"
}

type beq struct{}

func (o beq) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o beq) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0xF0, addressingMode: processor.Relative, cycles: 2},
	}
}

func (o beq) GetName() string {
	return "BEQ"
}

type brk struct{}

func (o brk) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o brk) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x00, addressingMode: nil, cycles: 7},
	}
}

func (o brk) GetName() string {
	return "BRK"
}

type rti struct{}

func (o rti) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o rti) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x00, addressingMode: nil, cycles: 6},
	}
}

func (o rti) GetName() string {
	return "RTI"
}

type jsr struct{}

func (o jsr) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o jsr) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x20, addressingMode: processor.Absolute, cycles: 6},
	}
}

func (o jsr) GetName() string {
	return "JSR"
}

type rts struct{}

func (o rts) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o rts) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x60, addressingMode: nil, cycles: 6},
	}
}

func (o rts) GetName() string {
	return "JSR"
}

type jmp struct{}

func (o jmp) Exec(console *processor.Console, variation *Variation) int {
	return 0

}

func (o jmp) getVariations() []Variation {
	return []Variation{
		Variation{opcode: 0x4C, addressingMode: processor.Absolute, cycles: 3},
		Variation{opcode: 0x6C, addressingMode: processor.Indirect, cycles: 5},
	}
}

func (o jmp) GetName() string {
	return "JMP"
}

var JumpOpCodes = []OpCode{
	blp{},
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
