package ui

import (
	"errors"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

func readController() {
	// 0 - A; 1 - B; 2 - Select; 3 - Start; 4 - Up; 5 - Down; 6 - Left; 7 - Right
	buttons := [8]bool{
		ebiten.IsKeyPressed(ebiten.KeyA),
		ebiten.IsKeyPressed(ebiten.KeyS),
		ebiten.IsKeyPressed(ebiten.KeyZ),
		ebiten.IsKeyPressed(ebiten.KeyX),
		ebiten.IsKeyPressed(ebiten.KeyUp),
		ebiten.IsKeyPressed(ebiten.KeyDown),
		ebiten.IsKeyPressed(ebiten.KeyLeft),
		ebiten.IsKeyPressed(ebiten.KeyRight),
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
