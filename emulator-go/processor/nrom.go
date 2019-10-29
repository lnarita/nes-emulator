package processor

import "log"

type Mapper interface {
	Read(address uint16) byte
	Write(address uint16, data byte)
}

type NROM struct {
	*Cartridge
	banks uint16
}

func (m NROM) Read(address uint16) byte {
	switch {
	case address < 0x2000:
		return m.CHR[address]
	case address >= 0xC000:
		index := (m.banks-1)*0x4000 + (address - 0xC000)
		return m.ROM[index]
	case address >= 0x8000:
		index := address - 0x8000
		return m.ROM[index]
	case address >= 0x6000:
		index := address - 0x6000
		return m.RAM[index]
	default:
		log.Fatalf("unhandled NROM read at address: 0x%04X", address)
	}
	return 0
}

func (m NROM) Write(address uint16, data byte) {
	switch {
	case address < 0x2000:
		m.CHR[address] = data
	case address >= 0x8000:
		log.Printf("NROM write %02X (%d)", data, uint16(data)%m.banks)
	case address >= 0x6000:
		if m.Battery {
			index := address - 0x6000
			m.RAM[index] = data
		}
	default:
		log.Fatalf("unhandled NROM write at address: 0x%04X", address)
	}
}

func CreateMapper(cartridge *Cartridge) Mapper {
	switch cartridge.Mapper {
	case 0x00:
		return NROM{cartridge, uint16(len(cartridge.ROM) / 0x4000)}
	default:
		return nil
	}
}
