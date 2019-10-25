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
	fmt.Println(console)

	fmt.Println(opcodes.AllOpCodes)
}
