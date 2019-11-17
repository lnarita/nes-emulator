package apu

type mixer struct {
	pulseTable []float64
	tndTable   []float64
}

func (m *mixer) init() {
	m.pulseTable = []float64{}
	m.tndTable = []float64{}
	for i := 0; i < 203; i++ {
		var entry float64 = 163.67 / (24329.0/float64(i) + 100.0)
		m.tndTable = append(m.tndTable, entry)
	}

	for i := 0; i < 31; i++ {
		var entry float64 = 95.52 / (8128.0/float64(i) + 100.0)
		m.pulseTable = append(m.pulseTable, entry)
	}
}

func (m *mixer) output(apu *APU) float64 {

	p1 := byte(10) //apu.pulse1.outputValue()
	// log.Println(p1)
	p2 := apu.pulse2.outputValue()
	t := apu.triangle.outputValue()
	n := apu.noise.outputValue()
	dmc := apu.dmc.outputValue()

	ptableValue := m.pulseTable[p1+p2]
	tndValue := m.tndTable[3*t+2*n+dmc]

	return ptableValue + tndValue
}
