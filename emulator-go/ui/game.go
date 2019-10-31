package ui

import (
	"errors"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
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

func game(screen *ebiten.Image) error {

	if ebiten.IsDrawingSkipped() {
		return nil
	}
	go readController()

	text.Draw(screen, "GAME", uiFont, 100, 75, color.White)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("Esc key pressed, shutting down")
	}

	return nil
}
