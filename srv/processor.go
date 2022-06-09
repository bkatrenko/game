package main

import (
	"context"
	"fmt"
	"sync"
)

type (
	Processor interface {
		HandleIncomingWorldState(ctx context.Context, upcomingState State) (State, error)
		Join(ctx context.Context, joinRequest JoinGame) (State, error)
	}

	processorImpl struct {
		games sync.Map
	}
)

func newProc() Processor {
	return &processorImpl{
		games: sync.Map{},
	}
}

func (p *processorImpl) HandleIncomingWorldState(ctx context.Context, incomingState State) (State, error) {
	updatesChan, ok := p.getGame(incomingState.ID)
	if !ok {
		return State{}, fmt.Errorf("can't fine the game with ID: %s", incomingState.ID)
	}

	responseChan := make(chan State)
	updatesChan <- GameInstanceUpdate{
		state:        incomingState,
		responseChan: responseChan,
	}
	currentGame := <-responseChan
	currentGame.CameFrom = incomingState.CameFrom

	return currentGame, nil
}

func (p *processorImpl) Join(ctx context.Context, joinRequest JoinGame) (State, error) {
	gameInstance := NewGameInstance(joinRequest)
	gameInstance.Start(context.Background())

	p.loadGame(joinRequest.GameID, gameInstance.getUpdateChan())
	return gameInstance.getState(), nil
}

func (p *processorImpl) getGame(id string) (chan GameInstanceUpdate, bool) {
	updatesChan, ok := p.games.Load(id)
	if !ok {
		return nil, false
	}

	typedUpdatesChan, ok := updatesChan.(chan GameInstanceUpdate)
	if !ok {
		return nil, false
	}

	return typedUpdatesChan, true
}

func (p *processorImpl) loadGame(id string, updatesChan chan GameInstanceUpdate) {
	p.games.Store(id, updatesChan)
}
