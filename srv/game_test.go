package main

import (
	"reflect"
	"testing"
)

func TestState_getSendData(t *testing.T) {
	type fields struct {
		ID            string
		CameFrom      string
		Player1       Rect
		Player2       Rect
		Ball          Rect
		MessageType   int8
		Message       string
		Player1Score  int8
		Player1Locked bool
		Player2Score  int8
		Player2Locked bool
	}
	tests := []struct {
		name   string
		fields fields
		want   State
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				ID:            tt.fields.ID,
				CameFrom:      tt.fields.CameFrom,
				Player1:       tt.fields.Player1,
				Player2:       tt.fields.Player2,
				Ball:          tt.fields.Ball,
				MessageType:   tt.fields.MessageType,
				Message:       tt.fields.Message,
				Player1Score:  tt.fields.Player1Score,
				Player1Locked: tt.fields.Player1Locked,
				Player2Score:  tt.fields.Player2Score,
				Player2Locked: tt.fields.Player2Locked,
			}
			if got := s.getSendData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("State.getSendData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_GetCurrentPlayer(t *testing.T) {
	type fields struct {
		ID            string
		CameFrom      string
		Player1       Rect
		Player2       Rect
		Ball          Rect
		MessageType   int8
		Message       string
		Player1Score  int8
		Player1Locked bool
		Player2Score  int8
		Player2Locked bool
	}
	tests := []struct {
		name   string
		fields fields
		want   Rect
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				ID:            tt.fields.ID,
				CameFrom:      tt.fields.CameFrom,
				Player1:       tt.fields.Player1,
				Player2:       tt.fields.Player2,
				Ball:          tt.fields.Ball,
				MessageType:   tt.fields.MessageType,
				Message:       tt.fields.Message,
				Player1Score:  tt.fields.Player1Score,
				Player1Locked: tt.fields.Player1Locked,
				Player2Score:  tt.fields.Player2Score,
				Player2Locked: tt.fields.Player2Locked,
			}
			if got := s.GetCurrentPlayer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("State.GetCurrentPlayer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestState_SetCurrentPlayer(t *testing.T) {
	type fields struct {
		ID            string
		CameFrom      string
		Player1       Rect
		Player2       Rect
		Ball          Rect
		MessageType   int8
		Message       string
		Player1Score  int8
		Player1Locked bool
		Player2Score  int8
		Player2Locked bool
	}
	type args struct {
		player Rect
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &State{
				ID:            tt.fields.ID,
				CameFrom:      tt.fields.CameFrom,
				Player1:       tt.fields.Player1,
				Player2:       tt.fields.Player2,
				Ball:          tt.fields.Ball,
				MessageType:   tt.fields.MessageType,
				Message:       tt.fields.Message,
				Player1Score:  tt.fields.Player1Score,
				Player1Locked: tt.fields.Player1Locked,
				Player2Score:  tt.fields.Player2Score,
				Player2Locked: tt.fields.Player2Locked,
			}
			s.SetCurrentPlayer(tt.args.player)
		})
	}
}
