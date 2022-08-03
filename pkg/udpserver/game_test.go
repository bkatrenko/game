package udpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzState_getSendData(f *testing.F) {
	f.Add("testGameID", "player1ID", "player2ID")
	f.Fuzz(func(t *testing.T, testGameID, player1ID, player2ID string) {
		state := State{
			ID: testGameID,
			Player1: Rect{
				ID: player1ID,
			},
			Player2: Rect{
				ID: player2ID,
			},
			CameFrom: player1ID,
		}

		sendData := state.getSendData()

		want := State{}
		want.Player2 = Rect{ID: player2ID}
		want.CameFrom = player1ID

		assert.Equal(t, want, sendData)
	})
}

func FuzzState_getCurrentPlayer(f *testing.F) {
	f.Add("player1ID", "player2ID")
	f.Fuzz(func(t *testing.T, player1ID, player2ID string) {
		state := State{
			Player1: Rect{
				ID: player1ID,
			},
			Player2: Rect{
				ID: player2ID,
			},
			CameFrom: player1ID,
		}

		currentPlayer := state.GetCurrentPlayer()
		assert.Equal(t, Rect{ID: player1ID}, currentPlayer)

		state.CameFrom = player2ID
		currentPlayer = state.GetCurrentPlayer()
		assert.Equal(t, Rect{ID: player2ID}, currentPlayer)
	})
}
