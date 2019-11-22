package apu

type frameCounterRegister struct {
	mode       byte
	irqInhibit bool
}

func (f *frameCounterRegister) read() byte {
	result := f.mode << 7
	if f.irqInhibit {
		result |= 0b01000000
	}

	return result
}

func (f *frameCounterRegister) write(data byte) {
	f.mode = data >> 7
	f.irqInhibit = (data & 0b01000000) == 0b01000000
}
