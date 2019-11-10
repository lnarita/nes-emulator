package ui

import (
	"log"

	"github.com/gordonklaus/portaudio"
)

type Audio struct {
	stream         *portaudio.Stream
	SampleRate     float64
	outputChannels int
	channel        chan float64
}

func NewAudio() *Audio {
	a := Audio{}
	return &a
}

func (a *Audio) Start(c chan float64) error {
	a.channel = c
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return err
	}
	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)
	stream, err := portaudio.OpenStream(parameters, a.Callback)
	if err != nil {
		return err
	}
	if err := stream.Start(); err != nil {
		return err
	}
	a.stream = stream
	a.SampleRate = parameters.SampleRate
	log.Println(parameters.SampleRate)
	a.outputChannels = parameters.Output.Channels
	return nil
}

func (a *Audio) Stop() error {
	return a.stream.Close()
}

func (a *Audio) Callback(out []float32) {
	var output float32
	for i := range out {
		if i%a.outputChannels == 0 {
			select {
			case sample := <-a.channel:
				output = float32(sample)
			default:
				output = 0
			}
		}
		out[i] = output
	}
}
