package main

import (
	"context"
	"encoding/json"
	"fmt"
	"game/client"
	"game/model"
	"image/color"
	_ "image/png"
	"log"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	serverUpdateAfter = time.Millisecond * 10
)

var (
	player1Score font.Face
	player2Score font.Face

	pointerImage  = ebiten.NewImage(model.PlaneWidth, model.PlaneHeight)
	ballImage     = ebiten.NewImage(model.BallDiameter, model.BallDiameter)
	opponentImage = ebiten.NewImage(model.PlaneWidth, model.PlaneHeight)
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	player1Score, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    26,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	player2Score, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    26,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	pointerImage.Fill(color.RGBA{0xff, 0, 0, 0xff})
	ballImage.Fill(color.RGBA{0xff, 0xff, 0, 0xff})
	opponentImage.Fill(color.RGBA{0, 0, 0xff, 0xff})
}

type (
	Game struct {
		UDPClient  *client.UDPClient
		HTTPClient *client.HTTPClient

		mu    sync.Mutex
		state model.State

		field Field
	}
)

func (g *Game) Update() error {
	dx, dy := ebiten.Wheel()

	g.mu.Lock()

	player := g.state.GetCurrentPlayer()

	player.RestrictSpeedLimit()
	player.UpdateXY(float32(dx), float32(dy), model.ScreenHeight, model.ScreenWidth)
	player.CalculateSpeed()

	g.state.SetCurrentPlayer(player)

	g.mu.Unlock()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.field.drawField(screen)

	text.Draw(screen, fmt.Sprint(g.state.Player1Score), player1Score, model.ScreenWidth/2-70, 20, color.White)
	text.Draw(screen, fmt.Sprint(g.state.Player2Score), player2Score, model.ScreenWidth/2+50, 20, color.White)

	opPlayer1 := &ebiten.DrawImageOptions{}
	opPlayer1.GeoM.Translate(float64(g.state.Player1.Vector.X), float64(g.state.Player1.Vector.Y))
	screen.DrawImage(pointerImage, opPlayer1)

	opPlayer2 := &ebiten.DrawImageOptions{}
	opPlayer2.GeoM.Translate(float64(g.state.Player2.Vector.X), float64(g.state.Player2.Vector.Y))
	screen.DrawImage(opponentImage, opPlayer2)

	opBall := &ebiten.DrawImageOptions{}
	opBall.GeoM.Translate(float64(g.state.Ball.Vector.X), float64(g.state.Ball.Vector.Y))
	screen.DrawImage(ballImage, opBall)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Ball\n(x:%0.2f, y:%0.2f)\n speedX: %0.2f speedY: %0.2f", g.state.Ball.Vector.X, g.state.Ball.Vector.Y, g.state.Player1.Speed.X, g.state.Player1.Speed.Y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return model.ScreenWidth, model.ScreenHeight
}

func (g *Game) run() {
	ticker := time.NewTicker(serverUpdateAfter)

	go func() {
		for {
			<-ticker.C

			currentState, err := json.Marshal(g.state.SendPlayerPos())
			if err != nil {
				println("error while marshal state", err.Error())
				continue
			}

			remoutState, err := g.UDPClient.Send(context.Background(), currentState)
			if err != nil {
				println("error while send state", err.Error())
				continue
			}

			var state model.State
			if err := json.Unmarshal(remoutState, &state); err != nil {
				println("error while marshal state", err.Error())
				continue
			}

			if state.MessageType == model.MessageTypeError {
				println("error from server received:", state.Message)
				continue
			}

			g.mu.Lock()
			if g.state.Player1.ID != state.CameFrom {
				g.state.Player1 = state.Player1
			}

			if g.state.Player2.ID != state.CameFrom {
				g.state.Player2 = state.Player2
			}

			g.state.Ball = state.Ball

			g.state.Player1Score = state.Player1Score
			g.state.Player2Score = state.Player2Score

			g.mu.Unlock()
		}

	}()
}

func main() {
	config, err := getConfig()
	if err != nil {
		panic(err)
	}

	udpClient, err := client.NewUDPClient(config.UDPServerHostPort)
	if err != nil {
		panic(err)
	}

	httpClient, err := client.NewHTTPlient(config.HTTPServerHostPort)
	if err != nil {
		panic(err)
	}

	state, err := httpClient.Join(model.JoinGame{
		GameID:       config.GameID,
		PlayerNumber: config.PlayerNumber,
		PlayedID:     config.PlayedID,
	})
	if err != nil {
		panic(err)
	}

	state.Player1.Image = pointerImage
	state.Player2.Image = opponentImage
	state.Ball.Image = ballImage
	state.CameFrom = config.PlayedID

	g := &Game{
		UDPClient: udpClient,
		state:     state,
		field: Field{
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
		},
	}

	go g.run()

	ebiten.SetWindowSize(model.ScreenWidth, model.ScreenHeight)
	ebiten.SetWindowTitle("Wheel (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
