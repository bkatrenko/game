package main

import (
	"context"
	"encoding/json"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// server wraps all the UDP echo server functionality.
// ps.: the server is capable of answering to a single
// client at a time.
func (s *server) runHTTPServer(ctx context.Context) error {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(30 * time.Second))

	r.Post("/game/join", func(w http.ResponseWriter, r *http.Request) {
		var joinMessage JoinGame

		if err := json.NewDecoder(r.Body).Decode(&joinMessage); err != nil {
			println("error while decode join request:", err.Error())
			s.writeError(w, "can't decode join request", http.StatusBadRequest)
			return
		}

		state, err := s.proc.join(joinMessage)
		if err != nil {
			println("error while join game:", err.Error())
			s.writeError(w, "can't join game", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(state); err != nil {
			println("error while encode join response:", err.Error())
		}
	})

	return http.ListenAndServe(s.listenAddressHTTP, r)
}

func (s *server) writeError(w http.ResponseWriter, text string, status int) {
	w.WriteHeader(status)

	_, err := w.Write([]byte(text))
	if err != nil {
		println("error while send response:", err.Error())
	}
}
