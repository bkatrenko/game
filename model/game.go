package model

const (
	BallDiameter  = 20
	BallRadius    = 25
	PlaneDiameter = 25

	GoalWidth  = 30
	GoalHeight = 90

	Player1GoalX = 0
	Player1GoalY = ScreenWidth/2 - GoalHeight*1.4

	Player2GoalX = ScreenWidth - GoalWidth
	Player2GoalY = ScreenWidth/2 - GoalHeight*1.4
)

type (
	State struct {
		ID       string `json:"id"`
		CameFrom string `json:"fr"`

		Player1 Rect `json:"p1"`
		Player2 Rect `json:"p2"`
		Ball    Rect `json:"b"`

		MessageType int8   `json:"mt,omitempty"`
		Message     string `json:"m,omitempty"`

		Player1Score  int8 `json:"s1"`
		Player1Locked bool `json:"-"`

		Player2Score  int8 `json:"s2"`
		Player2Locked bool `json:"-"`
	}

	JoinGame struct {
		GameID       string
		PlayerNumber int8
		PlayedID     string
	}
)

func (s *State) SendPlayerPos() State {
	currentPlayer := s.GetCurrentPlayer()
	newState := State{
		ID:       s.ID,
		CameFrom: s.CameFrom,
	}

	if currentPlayer.ID == s.Player1.ID {
		newState.Player1 = s.Player1
	}

	if currentPlayer.ID == s.Player2.ID {
		newState.Player2 = s.Player2
	}

	return newState
}

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
