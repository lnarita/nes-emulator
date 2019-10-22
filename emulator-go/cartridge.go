package main

import (
	"errors"
	"fmt"
)

// NESHeader the nes cartridge header contents
type NESHeader struct {
	prgRomSize byte
	chrRomSize byte
	flags6     byte
	flags7     byte
	prgRAMSize byte
	flags9     byte
	flags10    byte
}

// Cartridge the cartridge contents
type Cartridge struct {
	header NESHeader
	prgRom []byte
}

func cartridgeFromBytes(content []byte) (*Cartridge, error) {
	if len(content) > 16 {
		if "NES\x1A" == string(content[:4]) {
			fmt.Println("valid")
			// 		# 4: Size of PRG ROM in 16 KB units
			prgRomSize := content[4]
			// 		# 5: Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
			chrRomSize := content[5]
			flags6 := content[6]
			// 		# 7: Flags 7
			flags7 := content[7]
			// 		# 8: Size of PRG RAM in 8 KB units (Value 0 infers 8 KB for compatibility; see PRG RAM circuit)
			prgRAMSize := content[8]
			// 		# 9: Flags 9
			flags9 := content[9]
			// 		# 10: Flags 10 (unofficial)
			flags10 := content[10]
			// 		# 11-15: Zero filled
			header := NESHeader{prgRomSize: prgRomSize, chrRomSize: chrRomSize, flags6: flags6, flags7: flags7,
				prgRAMSize: prgRAMSize, flags9: flags9, flags10: flags10}

			prgRom := content[16:(header.prgRomSize + 16)]

			return &Cartridge{header: header, prgRom: prgRom}, nil
		}
	}

	return nil, errors.New("can't work with 42")
}
