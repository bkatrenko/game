package main

import (
	"fmt"
	"sync"
	"time"
)

type proc struct {
	games sync.Map
}

func newProc() *proc {
	return &proc{
		games: sync.Map{},
	}
}

func (p *proc) handle(upcomingState State) (State, error) {
	currentGame, ok := p.getGame(upcomingState.ID)
	if !ok {
		return State{}, fmt.Errorf("can't fine the game with ID: %s", upcomingState.ID)
	}

	currentCame := p.modifyState(currentGame, upcomingState)
	p.loadGame(currentCame)

	currentCame.CameFrom = upcomingState.CameFrom
	return currentCame, nil
}

func (p *proc) join(joinRequest JoinGame) (State, error) {
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

func (p *proc) startModifier() {
	ticker := time.NewTicker(time.Millisecond * 20)

	for {
		<-ticker.C

		currentState, ok := p.getGame("1")
		if !ok {
			continue
		}

		currentState.Ball.RestrictSpeedLimit()
		currentState.Ball.SlowDown()
		currentState.Ball.UpdateXY(currentState.Ball.Speed.X, currentState.Ball.Speed.Y, ScreenHeight, ScreenWidth)

		//if currentState.Player1.HasCollisionWith(currentState.Ball) {
		//currentState.Ball.AddSpeed(currentState.Player1.Speed.X, currentState.Player1.Speed.Y)
		//}

		//if currentState.Player2.HasCollisionWith(currentState.Ball) {
		//currentState.Ball.AddSpeed(currentState.Player2.Speed.X, currentState.Player2.Speed.Y)
		//}

		if currentState.Ball.HasCollisionWith(currentState.Player1) {
			currentState.Ball.ReflectFrom(currentState.Player1)
			currentState.Ball.Vector.X += currentState.Ball.Speed.X / 2
			currentState.Ball.Vector.Y += currentState.Ball.Speed.Y / 2

		}

		// if currentState.Ball.HasCollisionWith(currentState.Player2) {
		// 	currentState.Ball.ReflectFrom(currentState.Player2)
		// 	//currentState.Ball.Vector.X += currentState.Ball.Speed.X
		// 	//currentState.Ball.Vector.Y += currentState.Ball.Speed.Y
		// }

		currentState = p.checkPlayer1Goal(currentState)
		currentState = p.checkPlayer2Goal(currentState)

		p.loadGame(currentState)

	}
}

func (p *proc) modifyState(currentState, upcomingState State) State {
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

	currentState = p.checkPlayer1Goal(currentState)
	currentState = p.checkPlayer2Goal(currentState)

	return currentState
}

func (p *proc) checkPlayer1Goal(state State) State {
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

func (p *proc) checkPlayer2Goal(state State) State {
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

func (p *proc) getGame(id string) (State, bool) {
	state, ok := p.games.Load(id)
	if !ok {
		return State{}, false
	}

	typedState, ok := state.(State)
	if !ok {
		return State{}, false
	}

	return typedState, true
}

func (p *proc) loadGame(state State) {
	p.games.Store(state.ID, state)
}

func (p *proc) startStateFromJoin(joinRequest JoinGame) State {
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

	return state
}
