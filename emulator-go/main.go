package main

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"students.ic.unicamp.br/goten/opcodes"
	"students.ic.unicamp.br/goten/processor"
	"students.ic.unicamp.br/goten/ui"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	args := os.Args
	log.Printf("%s", args)
	fileName := args[1]

	emulate(fileName)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func emulate(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	check(err)

	car, err := processor.CartridgeFromBytes(data)
	check(err)

	mem := processor.Load(car)
	cpu := processor.Setup(mem, true)
	ppu := &processor.PPU{}
	controller1 := &processor.Controller{}
	controller2 := &processor.Controller{}

	console := processor.Console{CPU: cpu, PPU: ppu, Memory: mem, Controller1: controller1, Controller2: controller2}
	ppu.Console = &console

	go ui.InitUI()
	state := opcodes.State{}
	for {

		start := time.Now()

		state.ClearState()
		state.PC = console.CPU.PC
		state.SP = console.CPU.SP
		state.A = console.CPU.A
		state.X = console.CPU.X
		state.Y = console.CPU.Y
		state.Flags = console.CPU.Flags
		state.Cycle = console.CPU.Cycle

		decoded := fetchAndDecodeInstruction(&console)
		console.CPU.PC++

		state.OpCodeName = decoded.Opc.GetName()
		state.OpCode = decoded.Variation

		cycle := decoded.Opc.Exec(&console, &decoded.Variation, &state)

		for i := 0; i < cycle; i++ {
			console.Tick()
		}

		if state.OpCodeName == "BRK" {
			break
		}

		elapsed := time.Since(start).Seconds()
		expected := processor.CyclePeriod * float64(cycle)
		log.Printf("%s\n", state)
		//log.Printf("%s\n| elapsed: %0.15f - expected: %0.15f\n", state, elapsed, expected)
		if elapsed >= expected {
			log.Printf("<<<<< (%4s) - %0.15f >>>>>\n", state.OpCodeName, elapsed/expected)
		} else {
			time.Sleep(time.Duration(expected-elapsed) * time.Second)
		}
		time.Sleep(500000)
	}

}

func fetchAndDecodeInstruction(console *processor.Console) opcodes.MapValue {
	instruction := console.FetchData(console.CPU.PC)
	decoded := opcodes.AllOpCodes[instruction]
	return decoded
}
