package processor

import "students.ic.unicamp.br/goten/common"

type Memory struct {
	RAM []byte
	ROM *Cartridge
}

func Load(cartridge *Cartridge) Memory {
	return Memory{RAM: make([]byte, 2*common.KB), ROM: cartridge}
}

func (mem Memory) Fetch(address int) byte {
}
