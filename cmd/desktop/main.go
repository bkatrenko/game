package main

import (
	"game/client"
	"game/model"
	_ "image/png"
	"log"

	configs "github.com/bkatrenko/game/configs/desktop"
	"github.com/bkatrenko/game/pkg/desktop"
	"github.com/hajimehoshi/ebiten/v2"
)

func init() {
	desktop.InitScoresText()
	ebiten.SetWindowSize(model.ScreenWidth, model.ScreenHeight)
	ebiten.SetWindowTitle("Simple online game")
}

func main() {
	config, err := configs.GetConfigFromEnv()
	if err != nil {
		panic(err)
	}

	state, err := desktop.Join(config)
	if err != nil {
		panic(err)
	}

	g := desktop.NewGame(client.NewUDPClient(config.UDPServerHostPort), state, desktop.NewField())

	go g.Run()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
