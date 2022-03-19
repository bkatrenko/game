package model

const (
	BallDiameter = 25
	PlaneWidth   = 35
	PlaneHeight  = 35

	GoalWidth  = 30
	GoalHeight = 90

	Player1GoalX = 0
	Player1GoalY = ScreenWidth/2 - GoalHeight*1.4

	Player2GoalX = ScreenWidth - GoalWidth
	Player2GoalY = ScreenWidth/2 - GoalHeight*1.4
)

type (
	State struct {
		ID       string
		CameFrom string

		Player1 Rect
		Player2 Rect
		Ball    Rect

		MessageType int8
		Message     string

		Player1Score  int8
		Player1Locked bool

		Player2Score  int8
		Player2Locked bool
	}

	JoinGame struct {
		GameID       string
		PlayerNumber int8
		PlayedID     string
	}
)

func (s *State) GetCurrentPlayer() Rect {
	if s.CameFrom == s.Player1.ID {
		return s.Player1
	}

	if s.CameFrom == s.Player2.ID {
		return s.Player2
	}

	panic("player ID is not consistent with any player")
}

func (s *State) SetCurrentPlayer(player Rect) {
	if s.CameFrom == s.Player1.ID {
		s.Player1 = player
		return
	}

	if s.CameFrom == s.Player2.ID {
		s.Player2 = player
		return
	}

	panic("player ID is not consistent with any player")
}
