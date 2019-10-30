package ui

import (
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	width  = 256
	height = 240
	scale  = 3
	title  = "GOTEN NES EMULATOR"
	fps    = 60
)

// update is called every frame (1/60 [s]).
func update(screen *ebiten.Image) error {

	// Write your game's logical update.

	if ebiten.IsDrawingSkipped() {
		// When the game is running slowly, the rendering result
		// will not be adopted.
		return nil
	}

	// Write your game's rendering.

	return nil
}

func InitUI() {
	// Call ebiten.Run to start your game loop.
	if err := ebiten.Run(update, width, height, scale, title); err != nil {
		log.Fatal(err)
	}
}
