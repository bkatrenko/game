package main

import (
	"encoding/json"
	"net/http"
)

func (s *httpServer) handleJoin(w http.ResponseWriter, r *http.Request) {
	var joinMessage JoinGame

	if err := json.NewDecoder(r.Body).Decode(&joinMessage); err != nil {
		println("error while decode join request:", err.Error())
		s.writeError(w, "can't decode join request", http.StatusBadRequest)
		return
	}

	state, err := s.proc.Join(joinMessage)
	if err != nil {
		println("error while join game:", err.Error())
		s.writeError(w, "can't join game", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(state); err != nil {
		println("error while encode join response:", err.Error())
	}
}
