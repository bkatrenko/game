package main

import (
	"game/model"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Field struct {
	color     color.RGBA
	goalColor color.RGBA

	screenWidth  float64
	screenHeight float64

	centerHeight float64
	centerWidth  float64
}

func newField() Field {
	return Field{
		color: color.RGBA{
			R: 0,
			G: 255,
			B: 255,
			A: 255,
		},
		goalColor: color.RGBA{
			G: 255,
			B: 255,
			A: 255,
		},
		screenWidth:  model.ScreenWidth,
		screenHeight: model.ScreenHeight,
		centerHeight: 100,
		centerWidth:  100,
	}
}

func (f *Field) drawField(dst *ebiten.Image) {
	ebitenutil.DrawLine(dst, f.screenWidth/2, 0, f.screenWidth/2, f.screenHeight, f.color)
	ebitenutil.DrawLine(dst, 0, model.ScreenHeight/2, f.screenWidth, f.screenHeight/2, f.color)
	ebitenutil.DrawRect(dst, f.screenWidth/2-f.centerWidth/2, f.screenHeight/2-f.centerHeight/2, f.centerWidth, f.centerHeight, f.color)

	ebitenutil.DrawRect(dst, model.Player1GoalX, model.Player1GoalY, model.GoalWidth, model.GoalHeight, f.goalColor)
	ebitenutil.DrawRect(dst, model.Player2GoalX, model.Player2GoalY, model.GoalWidth, model.GoalHeight, f.goalColor)
}

func drawObject(screen *ebiten.Image, vector model.Vector, diameter int, color color.Color) {
	if screen == nil {
		panic("screen is nil: possible game is not initialised properly")
	}

	object := &ebiten.DrawImageOptions{}
	object.GeoM.Translate(float64(vector.X), float64(vector.Y))
	model.DrawCircle(screen, float64(vector.X), float64(vector.Y), diameter, color, true)
}
