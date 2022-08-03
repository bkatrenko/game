package udpserver

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	GameUpdatePeriod = time.Millisecond * 40

	Player1Number = iota - 1
	Player2Number
)

type (
	GameInstance struct {
		state       State
		updatesChan chan GameInstanceUpdate
		ticker      *time.Ticker
	}

	GameInstanceUpdate struct {
		newPlayer    JoinGame
		state        State
		responseChan chan State
	}
)

func NewGameInstance(joinRequest JoinGame) (*GameInstance, error) {
	if err := joinRequest.validate(); err != nil {
		return nil, err
	}

	state := getDefaultState(joinRequest)

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
	}, nil
}

func (gi *GameInstance) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case update := <-gi.updatesChan:
				gi.handleUpdate(update)
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
				gi.state.Ball.Vector.X += gi.state.Ball.Speed.X
				gi.state.Ball.Vector.Y += gi.state.Ball.Speed.Y

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

func (gi *GameInstance) handleUpdate(update GameInstanceUpdate) {
	switch {
	case update.newPlayer.GameID != "":
		gi.state = gi.addPlayer(gi.state, update.newPlayer)
		gi.state.CameFrom = update.newPlayer.PlayedID
		update.responseChan <- gi.state
		return
	case update.state.ID != "":
		gi.state = gi.modifyState(gi.state, update.state)
		gi.state.CameFrom = update.state.CameFrom
		update.responseChan <- gi.state
	default:
		panic("can't handle the update: event seems to be empty")
	}
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

func (gi *GameInstance) addPlayer(currentState State, joinRequest JoinGame) State {
	switch joinRequest.PlayerNumber {
	case Player1Number:
		currentState.Player1.ID = joinRequest.PlayedID
		return currentState
	case Player2Number:
		currentState.Player2.ID = joinRequest.PlayedID
		return currentState
	default:
		panic("wrong join request number: should be 0 or 1, only two players allowed")
	}
}

func getDefaultState(joinRequest JoinGame) State {
	return State{
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
			Vector: NewVector(ScreenWidth/2, ScreenHeight/2),

			PrevX:      0.0,
			PrevY:      0.0,
			SpeedLimit: 6.0,
			Speed:      NewVector(0, 0),
		},
	}
}
