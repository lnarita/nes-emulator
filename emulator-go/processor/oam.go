package processor

type OAM struct {
	data [0x100]byte
}

func (oam *OAM) Read(address uint16) byte {
	return oam.data[address%0x100]
}

func (oam *OAM) Write(address uint16, data byte) {
	index := address % 0x100
	if index%4 == 2 {
		v := data % 0xE3
		oam.data[index] = v
	} else {
		oam.data[index] = data
	}
}
