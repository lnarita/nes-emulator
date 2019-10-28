package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"students.ic.unicamp.br/goten/opcodes"
	"students.ic.unicamp.br/goten/processor"
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

	cpu := processor.Setup(mem)

	ppu := &processor.PPU{}
	console := processor.Console{CPU: cpu, PPU: ppu, Memory: mem}

	running := true

	for running {
		decoded := fetchAndDecodeInstruction(&console)
		decoded.Opc.Exec(&console, &decoded.Variation)
		console.CPU.PC++
		if decoded.Opc.GetName() == "BRK" {
			running = false
		}

	}

}

func fetchAndDecodeInstruction(console *processor.Console) opcodes.MapValue {
	instruction := console.Memory.FetchData(console.CPU.PC)
	decoded := opcodes.AllOpCodes[instruction]
	return decoded
}
