package ui

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten"
	"students.ic.unicamp.br/goten/processor"
)

const (
	width  = 256
	height = 240
	scale  = 3
	title  = "GOTEN NES EMULATOR"
	fps    = 60

	menuScreen = iota
	gameScreen = iota
)

var (
	frame   *image.RGBA
	console *processor.Console
)

var keyboardKeys = [][]string{
	{"Esc", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=", "\\", "`", " "},
	{"Tab", "Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "[", "]", "BS"},
	{"Ctrl", "A", "S", "D", "F", "G", "H", "J", "K", "L", ";", "'", "Enter"},
	{"Shift", "Z", "X", "C", "V", "B", "N", "M", ",", ".", "/", " "},
	{" ", "Alt", "Space", " ", " "},
	{},
	{"", "Up", ""},
	{"Left", "Down", "Right"},
}

var currentScreen = menuScreen

type rand struct {
	x, y, z, w uint32
}

func (r *rand) next() uint32 {
	// math/rand is too slow to keep 60 FPS on web browsers.
	// Use Xorshift instead: http://en.wikipedia.org/wiki/Xorshift
	t := r.x ^ (r.x << 11)
	r.x, r.y, r.z = r.y, r.z, r.w
	r.w = (r.w ^ (r.w >> 19)) ^ (t ^ (t >> 8))
	return r.w
}

var theRand = &rand{12345678, 4185243, 776511, 45411}

// update is called every frame (1/60 [s]).
func update(screen *ebiten.Image) error {
	if currentScreen == gameScreen {
		return game(screen)
	}
	return menu(screen)

}

// InitUI starts the interface of the program and opens a window
func InitUI(c *processor.Console) {
	console = c
	frame = image.NewRGBA(image.Rect(0, 0, width, height))

	// Call ebiten.Run to start your game loop.
	if err := ebiten.Run(update, width, height, scale, title); err != nil {
		log.Fatal(err)
	}
}
