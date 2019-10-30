package ui

import (
	"errors"
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

func game(screen *ebiten.Image) error {

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	text.Draw(screen, "GAME", uiFont, 100, 75, color.White)

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errors.New("Esc key pressed, shutting down")
	}

	return nil
}
