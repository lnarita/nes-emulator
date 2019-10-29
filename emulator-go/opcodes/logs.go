package opcodes

import (
	"fmt"
	"strings"
)

type State struct {
	PC             uint16
	SP             uint16
	A              byte
	X              byte
	Y              byte
	Flags          byte
	Address        []interface{}
	Data           byte
	HasData        bool
	Cycle          int
	PPUCycle       int
	OpCodeName     string
	OpCode         Variation
	Parameter1     byte
	Parameter2     byte
	ParameterCount byte
}

func (s State) String() string {
	var hexDump string
	switch s.ParameterCount {
	case 1:
		hexDump = fmt.Sprintf("%02X", s.Parameter1)
	case 2:
		hexDump = fmt.Sprintf("%02X %02X", s.Parameter1, s.Parameter2)
	default:
		hexDump = ""
	}
	var data string
	if s.OpCode.addressingMode != nil {
		addrFormatStr, _ := s.OpCode.addressingMode.AddressFormatString()
		var addressStr string
		var dataStr string
		addressStr = fmt.Sprintf(addrFormatStr, s.Address...)
		if s.HasData && s.OpCode.addressingMode != Accumulator && s.OpCode.addressingMode != Immediate {
			dataStr = fmt.Sprintf("= %02X", s.Data)
		} else {
			dataStr = ""
		}
		data = fmt.Sprintf("%s %s %s", s.OpCodeName, addressStr, dataStr)
	} else {
		data = fmt.Sprintf("%s", s.OpCodeName)
	}
	var opcode string
	if strings.HasPrefix(s.OpCodeName, "*") {
		opcode = fmt.Sprintf("%02X %-6s%-31s", s.OpCode.opcode, hexDump, data)
	} else {
		opcode = fmt.Sprintf("%02X %-6s %-30s", s.OpCode.opcode, hexDump, data)
	}
	state := fmt.Sprintf("A:%02X X:%02X Y:%02X P:%02X SP:%02X", s.A, s.X, s.Y, s.Flags, s.SP&0x00FF)
	return fmt.Sprintf("%04X  %40s  %s  CYC:%d", s.PC, opcode, state, s.Cycle)
}
