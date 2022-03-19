package main

import (
	"fmt"
	"game/model"
	"sync"

	"github.com/deeean/go-vector/vector2"
)

type proc struct {
	games sync.Map
}

func newProc() proc {
	return proc{
		games: sync.Map{},
	}
}

func (p *proc) handle(upcomingState model.State) (model.State, error) {
	currentGame, ok := p.getGame(upcomingState.ID)
	if !ok {
		return model.State{}, fmt.Errorf("can't fine the game with ID: %s", upcomingState.ID)
	}

	currentCame := p.modifyState(currentGame, upcomingState)
	p.loadGame(currentCame)

	currentCame.CameFrom = upcomingState.CameFrom
	return currentCame, nil
}

func (p *proc) join(joinRequest model.JoinGame) (model.State, error) {
	currentGame, ok := p.getGame(joinRequest.GameID)
	if !ok {
		currentGame = p.startStateFromJoin(joinRequest)
	}

	if joinRequest.PlayerNumber == 0 {
		currentGame.Player1.ID = joinRequest.PlayedID
	}

	if joinRequest.PlayerNumber == 1 {
		currentGame.Player2.ID = joinRequest.PlayedID
	}

	p.loadGame(currentGame)
	currentGame.CameFrom = joinRequest.PlayedID

	println("new game created, game ID:", joinRequest.GameID,
		"player ID:", joinRequest.PlayedID)
	return currentGame, nil
}

func (p *proc) modifyState(currentState, upcomingState model.State) model.State {
	if upcomingState.CameFrom == currentState.Player1.ID {
		currentState.Player1 = upcomingState.Player1
	}

	if upcomingState.CameFrom == currentState.Player2.ID {
		currentState.Player2 = upcomingState.Player2
	}

	currentState.Ball.RestrictSpeedLimit()
	currentState.Ball.SlowDown()
	currentState.Ball.UpdateXY(currentState.Ball.Speed.X, currentState.Ball.Speed.Y, model.ScreenHeight, model.ScreenWidth)

	if currentState.Player1.HasCollisionWith(currentState.Ball) {
		currentState.Ball.AddSpeed(currentState.Player1.Speed.X, currentState.Player1.Speed.Y)
	}

	if currentState.Player2.HasCollisionWith(currentState.Ball) {
		currentState.Ball.AddSpeed(currentState.Player2.Speed.X, currentState.Player2.Speed.Y)
	}

	currentState = p.checkPlayer1Goal(currentState)
	currentState = p.checkPlayer2Goal(currentState)

	return currentState
}

func (p *proc) checkPlayer1Goal(state model.State) model.State {
	if state.Ball.Vector.X <= 0+model.GoalWidth &&
		state.Ball.Vector.Y >= model.Player1GoalY &&
		state.Ball.Vector.Y+model.BallDiameter <= model.Player1GoalY+model.GoalHeight {

		if !state.Player1Locked {
			state.Player2Score++
			state.Player1Locked = true
		}

		return state
	}

	state.Player1Locked = false
	return state
}

func (p *proc) checkPlayer2Goal(state model.State) model.State {
	if state.Ball.Vector.X+model.BallDiameter >= model.ScreenWidth-model.GoalWidth &&
		state.Ball.Vector.Y >= model.Player2GoalY &&
		state.Ball.Vector.Y+model.BallDiameter <= model.Player2GoalY+model.GoalHeight {

		if !state.Player2Locked {
			state.Player1Score++
			state.Player2Locked = true
		}

		return state
	}

	state.Player2Locked = false
	return state
}

func (p *proc) getGame(id string) (model.State, bool) {
	state, ok := p.games.Load(id)
	if !ok {
		return model.State{}, false
	}

	typedState, ok := state.(model.State)
	if !ok {
		return model.State{}, false
	}

	return typedState, true
}

func (p *proc) loadGame(state model.State) {
	p.games.Store(state.ID, state)
}

func (p *proc) startStateFromJoin(joinRequest model.JoinGame) model.State {
	state := model.State{
		ID: joinRequest.GameID,
		Player1: model.Rect{
			Width:  model.PlaneHeight,
			Height: model.PlaneHeight,
			Vector: *vector2.New(0.0, model.ScreenHeight/2),

			PrevX:      0.0,
			PrevY:      0.0,
			SpeedLimit: 10.0,
			Speed:      *vector2.New(0, 0),
		},
		Player2: model.Rect{
			Width:      model.PlaneHeight,
			Height:     model.PlaneHeight,
			Vector:     *vector2.New(model.ScreenWidth-model.PlaneWidth, model.ScreenHeight/2),
			PrevX:      0.0,
			PrevY:      0.0,
			SpeedLimit: 10.0,
			Speed:      *vector2.New(0, 0),
		},
		Ball: model.Rect{
			Width:  model.BallDiameter,
			Height: model.BallDiameter,
			Vector: *vector2.New(200, 200),

			PrevX:      0.0,
			PrevY:      0.0,
			SpeedLimit: 10.0,
			Speed:      *vector2.New(0, 0),
		},
	}

	if joinRequest.PlayerNumber == 0 {
		state.Player1.ID = joinRequest.PlayedID
	}
	if joinRequest.PlayerNumber == 1 {
		state.Player2.ID = joinRequest.PlayedID
	}

	return state
}
