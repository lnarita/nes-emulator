package ui

import (
	"errors"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

var (
	pixels = make([]byte, width*height*4)
)

func readController() {
	buttons := [8]bool{
		ebiten.IsKeyPressed(ebiten.KeyA),     // 0 - A
		ebiten.IsKeyPressed(ebiten.KeyS),     // 1 - B
		ebiten.IsKeyPressed(ebiten.KeyZ),     // 2 - SELECT
		ebiten.IsKeyPressed(ebiten.KeyX),     // 3 - START
		ebiten.IsKeyPressed(ebiten.KeyUp),    // 4 - UP
		ebiten.IsKeyPressed(ebiten.KeyDown),  // 5 - DOWN
		ebiten.IsKeyPressed(ebiten.KeyLeft),  // 6 - LEFT
		ebiten.IsKeyPressed(ebiten.KeyRight), // 7 - RIGHT
	}
	console.Controller1.SetButtons(buttons)
}

// Draw paints current game state.
func draw() {
	for i, v := range console.PPU.Pixels {
		pixels[4*i] = v.R
		pixels[4*i+1] = v.G
		pixels[4*i+2] = v.B
		pixels[4*i+3] = v.A

	}
}

func game(screen *ebiten.Image) error {

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	go readController()

	text.Draw(screen, "GAME", uiFont, 100, 75, color.White)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("Esc key pressed, shutting down")
	}

	draw()
	screen.ReplacePixels(pixels)

	return nil
}
