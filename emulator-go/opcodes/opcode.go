package opcodes

import (
	"students.ic.unicamp.br/goten/processor"
)

type Variation struct {
	opcode         byte
	addressingMode processor.AddressMode
	cycles         int
}

type OpCode interface {
	getVariations() []Variation
	exec(*processor.Console, *Variation)
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
		UnofficialOpcodes,
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
