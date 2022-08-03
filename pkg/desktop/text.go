package desktop

import (
	"game/model"
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	//	serverUpdateAfter = time.Millisecond * 10
	serverUpdateAfter = time.Millisecond * 50

	scoresDPI      = 72
	scoresTextSize = 26

	player1ScoresPositionX = model.ScreenWidth/2 - 70
	player1ScoresPositionY = 20

	player2ScoresPositionX = model.ScreenWidth/2 + 50
	player2ScoresPositionY = 20
)

// Global variables are here.
// Usually it is better to NEVER use global variables while in this case we are keeping some
// kind of "global" state of the application that makes an app quiet impossible to properly test.
// In exactly this case my motivation is to define font & colors that should be initialised only once and
// should be never changed after:
// I don't want to put the variable into "Game" object 'cause it will make Game object dirty.
// In case we use global variables, it is always better to keep them in "main.go" or cmd package -
// it would be much more cleaner for maintainers.
var (
	player1Score font.Face
	player2Score font.Face

	player1Color = color.RGBA{0xff, 0, 0, 0xff}
	player2Color = color.RGBA{0, 0, 0xff, 0xff}
	ballColor    = color.RGBA{0xff, 0xff, 0xff, 0}
)

func InitScoresText() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	player1Score, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    scoresTextSize,
		DPI:     scoresDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	player2Score, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    scoresTextSize,
		DPI:     scoresDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}
