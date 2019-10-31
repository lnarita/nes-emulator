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

var currentScreen = menuScreen

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
