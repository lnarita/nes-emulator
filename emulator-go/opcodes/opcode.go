package opcodes

import (
	"fmt"

	"students.ic.unicamp.br/goten/processor"
)

type Variation struct {
	opcode         byte
	addressingMode processor.AddressMode
	cycles         int
}

type OpCode interface {
	getVariations() []Variation
	Exec(*processor.Console, *Variation) (int, LoggingStruct)
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

type LoggingStruct struct {
	low            *int
	high           *int
	addr           *int
	data           *int
	name           string
	opcode         int
	addressingMode *processor.AddressMode
}

func strAddr(logging LoggingStruct) string {
	if logging.addressingMode != nil {
		if logging.low != nil && logging.high != nil {
			return fmt.Sprintf("%2x %2x", logging.low, (*logging.high >> 8))
		} else if logging.low != nil {
			return fmt.Sprintf("%2x", logging.low)
		}
	}
	return ""
}

func strAddr2(logging LoggingStruct) string {
	if logging.addressingMode != nil {
		if logging.addr != nil && logging.data != nil {
			return fmt.Sprintf("%v %v %v", logging.name, logging.addr, logging.data)
		} else if logging.addr != nil {
			return fmt.Sprintf("%v %v", logging.name, logging.addr)
		}
	}
	return fmt.Sprintf("%v", logging.name)
}

func PrintOpCode(logging LoggingStruct) string {
	return fmt.Sprintf("%2x %5v %30v", logging.opcode, strAddr(logging), strAddr2(logging))
}
