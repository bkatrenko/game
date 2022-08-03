package udpserver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testJoinGame = JoinGame{
		GameID:       "1",
		PlayerNumber: 1,
		PlayedID:     "id",
	}
)

func TestNewGameInstance(t *testing.T) {
	testGameInstance := getDefaultState(testJoinGame)
	testGameInstance.Player2.ID = "id"

	type args struct {
		joinRequest JoinGame
	}
	tests := []struct {
		name    string
		args    args
		want    *GameInstance
		wantErr bool
	}{
		{
			name: "empty game ID",
			args: args{
				joinRequest: JoinGame{
					GameID:       "",
					PlayerNumber: 0,
					PlayedID:     "id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty player ID",
			args: args{
				joinRequest: JoinGame{
					GameID:       "1",
					PlayerNumber: 0,
					PlayedID:     "",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid player number",
			args: args{
				joinRequest: JoinGame{
					GameID:       "1",
					PlayerNumber: 99,
					PlayedID:     "id",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "successful game creation",
			args: args{
				joinRequest: JoinGame{
					GameID:       "1",
					PlayerNumber: 1,
					PlayedID:     "id",
				},
			},
			want: &GameInstance{
				state:       testGameInstance,
				updatesChan: make(chan GameInstanceUpdate),
				ticker:      time.NewTicker(GameUpdatePeriod),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGameInstance(tt.args.joinRequest)
			if tt.want != nil {
				assert.Equal(t, tt.want.state, got.state)
			}

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestGameInstance_Start(t *testing.T) {
	defaultState := getDefaultState(testJoinGame)
	defaultState.Player1.ID = "id1"
	defaultState.Player2.ID = "id2"

	instance := GameInstance{
		state:       defaultState,
		updatesChan: make(chan GameInstanceUpdate),
		ticker: &time.Ticker{
			C: make(<-chan time.Time),
		},
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	instance.Start(ctx)

	update := GameInstanceUpdate{
		state: State{
			ID: testJoinGame.GameID,
		},
		responseChan: make(chan State),
	}
	instance.updatesChan <- update
	response := <-update.responseChan

	assert.NotEmpty(t, response)
	cancelFunc()
}
