package main

import (
	"fmt"
	"io/ioutil"
	//"log"
	"os"
	//"time"

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

	cpu := processor.Setup(mem, true)

	ppu := &processor.PPU{}
	console := processor.Console{CPU: cpu, PPU: ppu, Memory: mem}

	running := true

	// initUI()

	for running {
		state := opcodes.State{PC: console.CPU.PC, SP: console.CPU.SP, A: console.CPU.A, X: console.CPU.X, Y: console.CPU.Y, Flags: console.CPU.Flags, Cycle: console.CPU.Cycle}
		if state.PC == 0xE928 {
			fmt.Println("ugh")
		}
		//start := time.Now()
		decoded := fetchAndDecodeInstruction(&console)
		console.CPU.PC++

		state.OpCodeName = decoded.Opc.GetName()
		state.OpCode = decoded.Variation

		cycle := decoded.Opc.Exec(&console, &decoded.Variation, &state)
		for i:=0;i<cycle;i++ {
			console.CPU.Cycle++
			if console.CPU.Cycle % 1000 == 0 {
				fmt.Printf("ping\n")
			}
		}
		if decoded.Opc.GetName() == "BRK" {
			running = false
		}
		//elapsed := time.Since(start)
		//log.Printf("elapsed: %0.15f - expected: %0.15f\n", elapsed.Seconds(), processor.CyclePeriod*float64(cycles))
		fmt.Printf("%s\n", state)
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