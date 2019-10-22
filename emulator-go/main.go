package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"students.ic.unicamp.br/goten/cpu"
)

func main() {
	args := os.Args

	fmt.Println(args)
	emulate(args[1])
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func emulate(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	check(err)

	car, err := cpu.CartridgeFromBytes(data)
	check(err)
	fmt.Println(car)

}
