package main

import (
	"fmt"
	"io/ioutil"
	"os"

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
	fmt.Println(car)

	mem := processor.Load(car)
	fmt.Println(mem)

	cpu := processor.Setup(mem)
	fmt.Println(cpu)
}
