package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.0/glfw"
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

	// initUI()

	for running {
		decoded := fetchAndDecodeInstruction(&console)
		_, opcodeLogging := decoded.Opc.Exec(&console, &decoded.Variation)
		console.CPU.PC++
		if decoded.Opc.GetName() == "BRK" {
			running = false
		}
		fmt.Printf("%04X  %s  %s  CYC:%d\n", console.CPU.PC, opcodes.PrintOpCode(opcodeLogging), console.CPU.String(), console.CPU.Cycle)
		fmt.Print(opcodes.PrintOpCode(opcodeLogging))

	}

}

func fetchAndDecodeInstruction(console *processor.Console) opcodes.MapValue {
	instruction := console.Memory.FetchData(console.CPU.PC)
	decoded := opcodes.AllOpCodes[instruction]
	return decoded
}

const (
	width  = 256
	height = 240
	scale  = 3
	title  = "GOTEN NES EMULATOR"
)

func initUI() {
	// initialize glfw
	success := glfw.Init()
	if !success {
		panic("Failed to init glfw")
	}
	defer glfw.Terminate()

	// create window
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, glfwerr := glfw.CreateWindow(width*scale, height*scale, title, nil, nil)
	if glfwerr != nil {
		panic(glfwerr)
	}
	window.MakeContextCurrent()

	// initialize gl
	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Enable(gl.TEXTURE_2D)

	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
