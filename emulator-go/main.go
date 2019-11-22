package main

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/gordonklaus/portaudio"
	"students.ic.unicamp.br/goten/opcodes"
	"students.ic.unicamp.br/goten/processor"
	"students.ic.unicamp.br/goten/processor/apu"
	"students.ic.unicamp.br/goten/ui"
)

var start time.Time
var cycles int = 0

func main() {
	// we need a parallel OS thread to avoid audio stuttering
	runtime.GOMAXPROCS(1)

	portaudio.Initialize()
	defer portaudio.Terminate()

	audio := ui.NewAudio()
	if err := audio.Start(); err != nil {
		log.Fatalln(err)
	}
	defer audio.Stop()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	args := os.Args
	//log.Printf("%s", args)
	fileName := args[1]

	data, err := ioutil.ReadFile(fileName)
	check(err)

	car, err := processor.CartridgeFromBytes(data)
	check(err)

	mem := processor.Load(car)
	cpu := processor.Setup(mem, false)
	ppu := &processor.PPU{}
	apu := &apu.APU{}
	apu.SampleRate = processor.CPUFrequency / audio.SampleRate
	apu.Init(audio.Channel)
	controller1 := &processor.Controller{}
	controller2 := &processor.Controller{}

	console := processor.Console{Cartridge: car, CPU: cpu, PPU: ppu, APU: apu, Memory: mem, Controller1: controller1, Controller2: controller2}
	ppu.Console = &console
	ppu.Reset()

	go emulate(&console)
	ui.InitUI(&console)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func emulate(console *processor.Console) {
	state := opcodes.State{}
	for {
		if console.CPU.Stall > 0 {
			console.CPU.Stall--
			console.Tick()
			continue
		}

		if cycles == 0 {
			start = time.Now()
		}

		state.ClearState()
		state.PC = console.CPU.PC
		state.SP = console.CPU.SP
		state.A = console.CPU.A
		state.X = console.CPU.X
		state.Y = console.CPU.Y
		state.Flags = console.CPU.Flags
		state.Cycle = console.CPU.Cycle

		decoded := fetchAndDecodeInstruction(console)
		console.CPU.PC++

		state.OpCodeName = decoded.Opc.GetName()
		state.OpCode = decoded.Variation

		cycle := decoded.Opc.Exec(console, &decoded.Variation, &state)
		cycles += cycle
		//log.Printf("%s\n", state)
		console.CheckInterrupts()
		//if console.CPU.Cycle == 89346 {
		//	log.Printf("0000")
		//}
		for i := 0; i < cycle; i++ {
			console.Tick()
		}

		if state.OpCodeName == "BRK" {
			break
		}

		elapsed := time.Since(start).Seconds()
		expected := processor.CyclePeriod * float64(cycles)
		toSleep := (expected - elapsed) * 1000000
		if toSleep > 100 {
			time.Sleep(time.Duration(toSleep) * time.Microsecond)
			cycles = 0
		}

	}

}

func fetchAndDecodeInstruction(console *processor.Console) opcodes.MapValue {
	instruction := console.FetchData(console.CPU.PC)
	decoded := opcodes.AllOpCodes[instruction]
	return decoded
}
