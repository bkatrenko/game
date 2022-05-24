package main

const (
	BallDiameter = 20
	PlaneWidth   = 25
	PlaneHeight  = 25

	GoalWidth  = 30
	GoalHeight = 90

	Player1GoalY = ScreenWidth/2 - GoalHeight*1.4
	Player2GoalY = ScreenWidth/2 - GoalHeight*1.4
)

type (
	// State describe game's state of the world with locations, ball and scores.
	State struct {
		// ID is a current game identifier, should be a string.
		// TODO: it would be nice to make some limit about an ID length
		ID string `json:"id,omitempty"`
		// CameFrom contain the ID of the user who send current state
		CameFrom string `json:"fr,omitempty"`
		// Player1 contain player #1 data: id, location, speed, etc
		Player1 Rect `json:"p1,omitempty"`
		// Player2 contain player #2 data: id, location, speed, etc
		Player2 Rect `json:"p2,omitempty"`
		// Ball contains the data about ball: location, speed, etc
		Ball Rect `json:"b,omitempty"`
		// MessageType define the type of the message that server could send or receive:
		// It could be data message, notification to a client/server or any kind of
		// information that should be transferred in the game
		MessageType int8 `json:"mt,omitempty"`
		// Message contains text (information from server that should be shown to a user)
		// or debug data
		Message string `json:"m,omitempty"`
		// Player1Score contains scores of player #1
		Player1Score int8 `json:"s1,omitempty"`
		// Player1Locked responsible for locking score increasing if ball is inside the
		// gate of player #1
		Player1Locked bool `json:"-"`
		// Player2Score contains scores of player #2
		Player2Score int8 `json:"s2,omitempty"`
		// Player2Locked responsible for locking score increasing if ball is inside the
		// gate of player #2
		Player2Locked bool `json:"-"`
	}

	JoinGame struct {
		GameID       string
		PlayerNumber int8
		PlayedID     string
	}
)

// getSendData does a simple clean up before send it to a client.
// - We don't need to send to a client his own location
// - We don't need to send a Game ID while it is knowing for a client
func (s *State) getSendData() State {
	s.ID = ""

	if s.CameFrom == s.Player1.ID {
		s.Player1 = Rect{}
		return *s
	}

	s.Player2 = Rect{}
	return *s
}

// GetCurrentPlayer return state sender's data
func (s *State) GetCurrentPlayer() Rect {
	if s.CameFrom == s.Player1.ID {
		return s.Player1
	}

	return s.Player2
}

// SetCurrentPlayer set new sender's state
func (s *State) SetCurrentPlayer(player Rect) {
	if s.CameFrom == s.Player1.ID {
		s.Player1 = player
		return
	}

	s.Player2 = player
}
