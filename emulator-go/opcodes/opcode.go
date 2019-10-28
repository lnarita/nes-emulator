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
	Exec(*processor.Console, *Variation) int
	GetName() string
}

type MapValue struct {
	Opc       OpCode
	Variation Variation
}

func getAllOpCodes() map[byte]MapValue {
	opcodeArrays := [][]OpCode{
		ArithmeticAndLogicalOpCodes,
		FlagOpcodes,
		JumpOpCodes,
		MoveOpCodes,
		UnofficialOpcodes,
	}
	opcodeMap := map[byte]MapValue{}

	for _, opcodeList := range opcodeArrays {
		for _, opcode := range opcodeList {
			variations := opcode.getVariations()
			for _, variation := range variations {
				opcodeMap[variation.opcode] = MapValue{Opc: opcode, Variation: variation}
			}
		}
	}

	return opcodeMap
}

var AllOpCodes = getAllOpCodes()

func (o OpcCode) ToString() {

}
