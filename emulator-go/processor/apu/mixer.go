package apu

type mixer struct {
	pulseTable []float64
	tndTable   []float64
}

func (m *mixer) init() {
	for i := 0; i < 203; i++ {
		var entry float64 = 163.67 / (24329.0/float64(i) + 100.0)
		m.tndTable = append(m.tndTable, entry)
	}

	for i := 0; i < 31; i++ {
		var entry float64 = 95.52 / (8128.0/float64(i) + 100.0)
		m.pulseTable = append(m.pulseTable, entry)
	}
}
