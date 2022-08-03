package desktop

import (
	"context"
	"encoding/json"
	"fmt"
	"game/client"
	"game/model"
	"image/color"
	"sync"
	"time"

	"github.com/bkatrenko/game/configs/desktop"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	UDPClient  *client.UDP
	HTTPClient *client.HTTP

	mu    sync.Mutex
	State model.State

	Field Field
}

func NewGame(udpClient *client.UDP, state model.State, field Field) *Game {
	return &Game{
		UDPClient: udpClient,
		State:     state,
		Field:     field,
	}
}

func Join(config desktop.Config) (model.State, error) {
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
	g.Field.drawField(screen)

	text.Draw(screen, fmt.Sprint(g.State.Player1Score), player1Score, player1ScoresPositionX, player1ScoresPositionY, color.White)
	text.Draw(screen, fmt.Sprint(g.State.Player2Score), player2Score, player2ScoresPositionX, player2ScoresPositionY, color.White)

	drawObject(screen, g.State.Player1.Vector, model.PlaneDiameter, player1Color)
	drawObject(screen, g.State.Player2.Vector, model.PlaneDiameter, player2Color)
	drawObject(screen, g.State.Ball.Vector, model.BallDiameter, ballColor)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Ball\n(x:%0.2f, y:%0.2f)\n speedX: %0.2f speedY: %0.2f", g.State.Ball.Vector.X, g.State.Ball.Vector.Y, g.State.Player1.Speed.X, g.State.Player1.Speed.Y))
}

func (g *Game) Update() error {
	dx, dy := ebiten.Wheel()
	g.mu.Lock()

	player := g.State.GetCurrentPlayer()
	player.RestrictSpeedLimit()
	player.UpdateXY(float32(dx), float32(dy), model.ScreenHeight, model.ScreenWidth)
	player.CalculateSpeed()
	g.State.SetCurrentPlayer(player)

	g.mu.Unlock()

	g.sendUpdate()
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return model.ScreenWidth, model.ScreenHeight
}

func (g *Game) sendUpdate() {
	currentState, err := json.Marshal(g.State.SendPlayerPos())
	if err != nil {
		println("error while marshal state", err.Error())
		return
	}

	remoteState, err := g.UDPClient.Send(context.Background(), currentState)
	if err != nil {
		println("error while send state", err.Error())
		return
	}

	var state model.State
	if err := json.Unmarshal(remoteState, &state); err != nil {
		println("error while marshal state", err.Error())
		return
	}

	if state.MessageType == model.MessageTypeError {
		println("error from server received:", state.Message)
		return
	}

	g.mu.Lock()
	if g.State.Player1.ID != state.CameFrom {
		g.State.Player1 = state.Player1
	}

	if g.State.Player2.ID != state.CameFrom {
		g.State.Player2 = state.Player2
	}

	g.State.Ball = state.Ball

	g.State.Player1Score = state.Player1Score
	g.State.Player2Score = state.Player2Score

	g.mu.Unlock()
}

func (g *Game) Run() {
	ticker := time.NewTicker(serverUpdateAfter)

	go func() {
		for {
			<-ticker.C
			g.sendUpdate()
		}
	}()
}
