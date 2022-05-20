package main

import (
	"game/client"
	"game/model"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
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

func init() {
	initScoresText()
	ebiten.SetWindowSize(model.ScreenWidth, model.ScreenHeight)
	ebiten.SetWindowTitle("Simple online game")
}

func main() {
	config, err := getConfigFromEnv()
	if err != nil {
		panic(err)
	}

	state, err := join(config)
	if err != nil {
		panic(err)
	}

	g := &Game{
		UDPClient: client.NewUDPClient(config.UDPServerHostPort),
		state:     state,
		field:     newField(),
	}

	go g.run()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
