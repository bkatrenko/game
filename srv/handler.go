package main

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (s *httpServer) handleJoin(w http.ResponseWriter, r *http.Request) {
	var joinMessage JoinGame

	if err := json.NewDecoder(r.Body).Decode(&joinMessage); err != nil {
		log.Err(err).Msg("error while decode join request")
		s.writeError(w, "can't decode join request", http.StatusBadRequest)
		return
	}

	state, err := s.proc.Join(r.Context(), joinMessage)
	if err != nil {
		log.Err(err).Msg("error while join game")
		s.writeError(w, "can't join game", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(state); err != nil {
		log.Err(err).Msg("error while encode join response")
	}
}
