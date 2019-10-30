package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"students.ic.unicamp.br/goten/opcodes"
	"students.ic.unicamp.br/goten/processor"
	"students.ic.unicamp.br/goten/ui"
)

func main() {
	args := os.Args
	fmt.Println(args)
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
	console := processor.Console{CPU: cpu, PPU: ppu, Memory: mem}

	go ui.InitUI()
	for {

		start := time.Now()
		state := opcodes.State{PC: console.CPU.PC, SP: console.CPU.SP, A: console.CPU.A, X: console.CPU.X, Y: console.CPU.Y, Flags: console.CPU.Flags, Cycle: console.CPU.Cycle}

		decoded := fetchAndDecodeInstruction(&console)
		console.CPU.PC++

		state.OpCodeName = decoded.Opc.GetName()
		state.OpCode = decoded.Variation

		cycle := decoded.Opc.Exec(&console, &decoded.Variation, &state)
		for i := 0; i < cycle; i++ {
			console.CPU.Cycle++
		}
		if decoded.Opc.GetName() == "BRK" {
			break
		}
		elapsed := time.Since(start).Seconds()
		expected := processor.CyclePeriod * float64(cycle)
		fmt.Printf("%s\n", state)
		//fmt.Printf("%s\n| elapsed: %0.15f - expected: %0.15f\n", state, elapsed, expected)
		if elapsed >= expected {
			fmt.Printf("<<<<< (%s) - %0.15f >>>>>\n", state.OpCodeName, elapsed/expected)
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
