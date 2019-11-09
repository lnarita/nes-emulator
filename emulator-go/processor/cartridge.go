package processor

import (
	"errors"
	"fmt"
)

// INESHeader the nes cartridge Header contents
type INESHeader struct {
	PrgRomSize int
	ChrRomSize int
	Flags6     byte
	Flags7     byte
	PrgRAMSize int
	Flags9     byte
	Flags10    byte
}

func (header INESHeader) String() string {
	return fmt.Sprintf("INES { \"PRG-size\": %d, \"CHR-size\": %d, \"RAM-size\": %d, \"flag6\": %08b, \"flag7\": %08b, \"flag9\": %08b, \"flag10\": %08b }",
		header.PrgRomSize, header.ChrRomSize, header.PrgRAMSize, header.Flags6, header.Flags7, header.Flags9, header.Flags10)
}

// Cartridge the cartridge contents
type Cartridge struct {
	Header    INESHeader
	Trainer   []byte
	ROM       []byte
	CHR       []byte
	RAM       []byte
	Mapper    byte    // Mapper type
	Mirroring *Mirror // Mirroring mode. 0: horizontal; 1: vertical; 2: 4-screen VRAM
	Battery   bool    // Battery present
}

func (cartridge Cartridge) String() string {
	var mirroring string
	switch cartridge.Mirroring.id {
	case 0x0:
		mirroring = "horizontal"
	case 0x1:
		mirroring = "vertical"
	default:
		mirroring = "4-screen VRAM"
	}
	return fmt.Sprintf("Cartridge { \"Header\": %s, \"Trainer\": [% X], \"PRG\": [% X], \"CHR\": [% X], \"RAM\": [% X], \"Mapper\": %d, \"Mirroring\": %s, \"Battery\": %t}",
		cartridge.Header, cartridge.Trainer, cartridge.ROM, cartridge.CHR, cartridge.RAM, cartridge.Mapper, mirroring, cartridge.Battery)
}

const HeaderSize = 16
const TrainerSize = 512
const SaveRamSize = (0x7FFF - 0x6000) + 1

func CartridgeFromBytes(content []byte) (*Cartridge, error) {
	if len(content) > 16 {
		if "NES\x1A" == string(content[:4]) {
			// 4: Size of PRG ROM in 16 KB units
			prgRomSize := content[4]
			// 5: Size of CHR ROM in 8 KB units (Value 0 means the board uses CHR RAM)
			chrRomSize := content[5]
			flags6 := content[6]
			//76543210
			//||||||||
			//|||||||+- Mirroring: 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
			//|||||||              1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
			//||||||+-- 1: Cartridge contains Battery-backed PRG RAM ($6000-7FFF) or other persistent memory
			//|||||+--- 1: 512-byte Trainer at $7000-$71FF (stored before PRG data)
			//||||+---- 1: Ignore mirroring control or above mirroring bit; instead provide four-screen VRAM
			//++++----- Lower nybble of Mapper number
			mirroring := flags6 & 0x08
			if mirroring != 0 {
				mirroring = 0x02
			} else {
				mirroring = flags6 & 0x01
			}
			hasTrainerData := (flags6 & 0x04) != 0
			hasBattery := (flags6 & 0x02) != 0
			// # 7: Flags 7
			flags7 := content[7]
			//76543210
			//||||||||
			//|||||||+- VS Unisystem
			//||||||+-- PlayChoice-10 (8KB of Hint Screen data stored after CHR data)
			//||||++--- If equal to 2, flags 8-15 are in NES 2.0 format
			//++++----- Upper nybble of Mapper number
			mapperLow := flags6 & 0x0F
			mapperHigh := flags7 & 0xF0
			mapper := mapperHigh | mapperLow
			// 8: Size of PRG RAM in 8 KB units (Value 0 infers 8 KB for compatibility; see PRG RAM circuit)
			prgRAMSize := content[8]
			if prgRAMSize == 0 {
				prgRAMSize = 1
			}
			// 9: Flags 9
			flags9 := content[9]
			// 10: Flags 10 (unofficial)
			flags10 := content[10]
			// 11-15: Zero filled
			header := INESHeader{PrgRomSize: int(prgRomSize) * 16 * KB, ChrRomSize: int(chrRomSize) * 8 * KB, PrgRAMSize: int(prgRAMSize) * 8 * KB, Flags6: flags6, Flags7: flags7, Flags9: flags9, Flags10: flags10}

			//    Header (16 bytes)
			//    Trainer, if present (0 or 512 bytes)
			//    PRG ROM data (16384 * x bytes)
			//    CHR ROM data, if present (8192 * y bytes)
			//    PlayChoice INST-ROM, if present (0 or 8192 bytes)
			//    PlayChoice PROM, if present (16 bytes data, 16 bytes CounterOut) (this is often missing, see PC10 ROM-Images for details)
			trainerSize := 0
			if hasTrainerData {
				trainerSize = TrainerSize
			}
			offset := HeaderSize
			trainer := content[offset:(trainerSize + offset)]
			offset += len(trainer)
			prgRom := content[offset:(header.PrgRomSize + offset)]
			offset += header.PrgRomSize
			var chr []byte
			if header.ChrRomSize == 0 {
				chr = make([]byte, 8192)
			} else {
				chr = content[offset:(header.ChrRomSize + offset)]
			}
			sram := make([]byte, SaveRamSize)

			var mirror *Mirror
			switch mirroring {
			case 0:
				mirror = horizontal
			case 1:
				mirror = vertical
			case 2:
				mirror = single0
			case 3:
				mirror = single1
			default:
				mirror = four
			}

			return &Cartridge{Header: header, Trainer: trainer, ROM: prgRom, CHR: chr, Mapper: mapper, RAM: sram, Mirroring: mirror, Battery: hasBattery}, nil
		}
	}

	return nil, errors.New("invalid iNES Header")
}

type Mirror struct {
	id          byte
	lookupTable [4]uint16
}

func (mirror *Mirror) calcAddress(address uint16) uint16 {
	addr := (address - 0x2000) % 0x1000
	table := addr / 0x0400
	offset := addr % 0x0400
	return 0x2000 + mirror.lookupTable[table]*0x0400 + offset
}

var (
	horizontal = &Mirror{id: 0, lookupTable: [4]uint16{0, 0, 1, 1}}
	vertical   = &Mirror{id: 1, lookupTable: [4]uint16{0, 1, 0, 1}}
	single0    = &Mirror{id: 2, lookupTable: [4]uint16{0, 0, 0, 0}}
	single1    = &Mirror{id: 3, lookupTable: [4]uint16{1, 1, 1, 1}}
	four       = &Mirror{id: 4, lookupTable: [4]uint16{0, 1, 2, 3}}
)
