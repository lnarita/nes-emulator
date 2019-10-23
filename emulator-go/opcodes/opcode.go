package opcodes

import "students.ic.unicamp.br/goten/processor"

type Variation struct {
	opcode        byte
	addresingMode processor.AddressMode
	cycles        int
}

type OpCode interface {
	getVariations() []Variation
	exec(*processor.CPU, *processor.Memory, *processor.AddressMode)
	getName() string
}
