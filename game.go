package main

import (
	"context"
	"encoding/json"
	"fmt"
	"game/client"
	"game/model"
	"image/color"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	UDPClient  *client.UDP
	HTTPClient *client.HTTP

	mu    sync.Mutex
	state model.State

	field Field
}

func join(config Config) (model.State, error) {
	state, err := client.NewHTTPClient(config.HTTPServerHostPort).Join(model.JoinGame{
		GameID:       config.GameID,
		PlayerNumber: config.PlayerNumber,
		PlayedID:     config.PlayedID,
	})
	if err != nil {
		return model.State{}, fmt.Errorf("error while join the game: %w", err)
	}

	state.CameFrom = config.PlayedID
	return state, nil
}

// Draw function do all work on drawing everything on user's screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.field.drawField(screen)

	text.Draw(screen, fmt.Sprint(g.state.Player1Score), player1Score, player1ScoresPositionX, player1ScoresPositionY, color.White)
	text.Draw(screen, fmt.Sprint(g.state.Player2Score), player2Score, player2ScoresPositionX, player2ScoresPositionY, color.White)

	drawObject(screen, g.state.Player1.Vector, model.PlaneDiameter, player1Color)
	drawObject(screen, g.state.Player2.Vector, model.PlaneDiameter, player2Color)
	drawObject(screen, g.state.Ball.Vector, model.BallDiameter, ballColor)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Ball\n(x:%0.2f, y:%0.2f)\n speedX: %0.2f speedY: %0.2f", g.state.Ball.Vector.X, g.state.Ball.Vector.Y, g.state.Player1.Speed.X, g.state.Player1.Speed.Y))
}

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

			remoteState, err := g.UDPClient.Send(context.Background(), currentState)
			if err != nil {
				println("error while send state", err.Error())
				continue
			}

			var state model.State
			if err := json.Unmarshal(remoteState, &state); err != nil {
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
