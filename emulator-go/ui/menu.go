package ui

import (
	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

var (
	uiImage       *ebiten.Image
	uiFont        font.Face
	uiFontMHeight int
)

func init() {
	tt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}

	uiFont = truetype.NewFace(tt, &truetype.Options{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	b, _, _ := uiFont.GlyphBounds('M')
	uiFontMHeight = (b.Max.Y - b.Min.Y).Ceil()
}

func menu(screen *ebiten.Image) error {

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	text.Draw(screen, "GOTEN", uiFont, 100, 75, color.White)
	text.Draw(screen, "A NES EMULATOR", uiFont, 70, 90, color.White)
	text.Draw(screen, "PRESS ENTER TO BEGIN", uiFont, 50, 125, color.White)

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		currentScreen = gameScreen
	}

	return nil
}
