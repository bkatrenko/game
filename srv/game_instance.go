package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	GameUpdatePeriod = time.Millisecond * 40
)

type (
	GameInstance struct {
		state       State
		updatesChan chan GameInstanceUpdate
		ticker      *time.Ticker
	}

	GameInstanceUpdate struct {
		state        State
		responseChan chan State
	}
)

func NewGameInstance(joinRequest JoinGame) *GameInstance {
	state := State{
		ID: joinRequest.GameID,
		Player1: Rect{
			Width:  PlaneHeight,
			Height: PlaneHeight,
			Vector: NewVector(50.0, ScreenHeight/2),

			PrevX:      0.0,
			PrevY:      0.0,
			SpeedLimit: 6.0,
			Speed:      NewVector(0, 0),
		},
		Player2: Rect{
			Width:      PlaneHeight,
			Height:     PlaneHeight,
			Vector:     NewVector(ScreenWidth-PlaneWidth, ScreenHeight/2),
			PrevX:      0.0,
			PrevY:      0.0,
			SpeedLimit: 6.0,
			Speed:      NewVector(0, 0),
		},
		Ball: Rect{
			Width:  BallDiameter,
			Height: BallDiameter,
			Vector: NewVector(200, 200),

			PrevX:      0.0,
			PrevY:      0.0,
			SpeedLimit: 6.0,
			Speed:      NewVector(0, 0),
		},
	}

	if joinRequest.PlayerNumber == 0 {
		state.Player1.ID = joinRequest.PlayedID
	}

	if joinRequest.PlayerNumber == 1 {
		state.Player2.ID = joinRequest.PlayedID
	}

	return &GameInstance{
		state:       state,
		updatesChan: make(chan GameInstanceUpdate),
		ticker:      time.NewTicker(GameUpdatePeriod),
	}
}

func (gi *GameInstance) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case newRequest := <-gi.updatesChan:
				gi.state = gi.modifyState(gi.state, newRequest.state)
				gi.state.CameFrom = newRequest.state.CameFrom
				newRequest.responseChan <- gi.state

			case <-gi.ticker.C:
			case <-ctx.Done():
				log.Debug().Str("game_id", gi.state.ID).Msg("stop game instance: context is done")
				return
			}

			gi.state.Ball.RestrictSpeedLimit()
			gi.state.Ball.SlowDown()
			gi.state.Ball.UpdateXY(gi.state.Ball.Speed.X, gi.state.Ball.Speed.Y, ScreenHeight, ScreenWidth)

			if gi.state.Ball.HasCollisionWith(gi.state.Player1) {
				gi.state.Ball.ReflectFrom(gi.state.Player1)
				gi.state.Ball.Vector.X += gi.state.Ball.Speed.X / 2
				gi.state.Ball.Vector.Y += gi.state.Ball.Speed.Y / 2

			}

			if gi.state.Ball.HasCollisionWith(gi.state.Player2) {
				gi.state.Ball.ReflectFrom(gi.state.Player2)
				gi.state.Ball.Vector.X += gi.state.Ball.Speed.X
				gi.state.Ball.Vector.Y += gi.state.Ball.Speed.Y
			}

			gi.state = gi.checkPlayer1Goal(gi.state)
			gi.state = gi.checkPlayer2Goal(gi.state)
		}
	}()
}

func (gi *GameInstance) getState() State {
	return gi.state
}

func (gi *GameInstance) getUpdateChan() chan GameInstanceUpdate {
	return gi.updatesChan
}

func (gi *GameInstance) checkPlayer1Goal(state State) State {
	if state.Ball.Vector.X <= 0+GoalWidth &&
		state.Ball.Vector.Y >= Player1GoalY &&
		state.Ball.Vector.Y+BallDiameter <= Player1GoalY+GoalHeight {

		if !state.Player1Locked {
			state.Player2Score++
			state.Player1Locked = true
		}

		return state
	}

	state.Player1Locked = false
	return state
}

func (gi *GameInstance) checkPlayer2Goal(state State) State {
	if state.Ball.Vector.X+BallDiameter >= ScreenWidth-GoalWidth &&
		state.Ball.Vector.Y >= Player2GoalY &&
		state.Ball.Vector.Y+BallDiameter <= Player2GoalY+GoalHeight {

		if !state.Player2Locked {
			state.Player1Score++
			state.Player2Locked = true
		}

		return state
	}

	state.Player2Locked = false
	return state
}

func (gi *GameInstance) modifyState(currentState, upcomingState State) State {
	if upcomingState.CameFrom == currentState.Player1.ID {
		currentState.Player1 = upcomingState.Player1
	}

	if upcomingState.CameFrom == currentState.Player2.ID {
		currentState.Player2 = upcomingState.Player2
	}

	if currentState.Player1.HasCollisionWith(currentState.Ball) {
		currentState.Ball.AddSpeed(currentState.Player1.Speed.X, currentState.Player1.Speed.Y)
	}

	if currentState.Player2.HasCollisionWith(currentState.Ball) {
		currentState.Ball.AddSpeed(currentState.Player2.Speed.X, currentState.Player2.Speed.Y)
	}

	currentState = gi.checkPlayer1Goal(currentState)
	currentState = gi.checkPlayer2Goal(currentState)

	return currentState
}
