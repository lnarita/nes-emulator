package processor

type APU struct {
	// APU Registers per channel
	pulse1        pulseRegisters    // $4000-$4003
	pulse2        pulseRegisters    // $4004-$4007
	triangle      triangleRegisters // $4008-$400B
	noise         noiseRegisters    // $400C-$400F
	dmc           dmcRegisters      // $4010-$4013
	channelEnable uint16            // $4015
	frameCounter  uint16            // $4017

}

type pulseRegisters struct {
	timer         uint16
	lengthCounter uint16
	envelope      uint16
	sweep         uint16
}
type triangleRegisters struct {
	timer         uint16
	lengthCounter uint16
	linearCounter uint16
}

type noiseRegisters struct {
	timer                       uint16
	lengthCounter               uint16
	envelope                    uint16
	linearFeedbackShiftRegister uint16
}

type dmcRegisters struct {
	timer        uint16
	memoryReader uint16
	sampleBuffer uint16
	outputUnit   uint16
}
