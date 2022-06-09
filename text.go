package main

import (
	"game/model"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	//	serverUpdateAfter = time.Millisecond * 10
	serverUpdateAfter = time.Millisecond * 20

	scoresDPI      = 72
	scoresTextSize = 26

	player1ScoresPositionX = model.ScreenWidth/2 - 70
	player1ScoresPositionY = 20

	player2ScoresPositionX = model.ScreenWidth/2 + 50
	player2ScoresPositionY = 20
)

func initScoresText() {
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
