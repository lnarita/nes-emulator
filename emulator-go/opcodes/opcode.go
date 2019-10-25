package opcodes

import (
	"students.ic.unicamp.br/goten/processor"
)

type Variation struct {
	opcode        byte
	addresingMode processor.AddressMode
	cycles        int
}

type OpCode interface {
	getVariations() []Variation
	exec(*processor.CPU, *processor.Memory, *Variation)
	getName() string
}

type mapValue struct {
	opcode    OpCode
	variation Variation
}

func getAllOpCodes() map[byte]mapValue {
	opcodeArrays := [][]OpCode{
		ArithmeticAndLogicalOpCodes,
		FlagOpcodes,
		JumpOpCodes,
		MoveOpCodes,
	}
	opcodeMap := map[byte]mapValue{}

	for _, opcodeList := range opcodeArrays {
		for _, opcode := range opcodeList {
			variations := opcode.getVariations()
			for _, variation := range variations {
				opcodeMap[variation.opcode] = mapValue{opcode, variation}
			}
		}
	}

	return opcodeMap
}

var AllOpCodes = getAllOpCodes()
