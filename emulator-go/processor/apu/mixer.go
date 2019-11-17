package apu

import "log"

type mixer struct {
	pulseTable []float32
	tndTable   []float32
}

func (m *mixer) init() {
	m.pulseTable = []float32{}
	m.tndTable = []float32{}
	for i := 0; i < 203; i++ {
		var entry float32 = 163.67 / (24329.0/float32(i) + 100.0)
		m.tndTable = append(m.tndTable, entry)
	}

	for i := 0; i < 31; i++ {
		var entry float32 = 95.52 / (8128.0/float32(i) + 100.0)
		m.pulseTable = append(m.pulseTable, entry)
	}
}

func (m *mixer) output(apu *APU) float32 {

	p1 := apu.pulse1.outputValue()
	if p1 != 0 {
		log.Println(p1)
	}
	p2 := apu.pulse2.outputValue()
	t := apu.triangle.outputValue()
	n := apu.noise.outputValue()
	dmc := apu.dmc.outputValue()

	ptableValue := m.pulseTable[p1+p2]
	tndValue := m.tndTable[3*t+2*n+dmc]

	return ptableValue + tndValue
}
