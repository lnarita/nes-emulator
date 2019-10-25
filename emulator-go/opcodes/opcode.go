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
	exec(*processor.CPU, *processor.Memory, *processor.AddressMode)
	getName() string
}

func getAllOpCodes() map[byte]OpCode {
	opcodeArrays := [][]OpCode{
		ArithmeticAndLogicalOpCodes,
		FlagOpcodes,
		JumpOpCodes,
		MoveOpCodes,
	}
	opcodeMap := map[byte]OpCode{}

	for _, opcodeList := range opcodeArrays {
		for _, opcode := range opcodeList {
			variations := opcode.getVariations()
			for _, variation := range variations {
				opcodeMap[variation.opcode] = opcode
			}
		}
	}

	return opcodeMap
}

var AllOpCodes = getAllOpCodes()
