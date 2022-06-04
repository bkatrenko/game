package main

import (
	"context"

	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

const (
	gameJoinRoute = "/game/join"

	DefaultIdleTimeout       = time.Second * 60
	DefaultReadTimeout       = time.Second * 15
	DefaultWriteTimeout      = time.Second * 15
	DefaultReadHeaderTimeout = time.Second

	DefaultMaxHeaderBytes = 1024
)

type httpServer struct {
	proc   Processor
	server *http.Server
}

func newHTTPServer(config config, proc Processor) *httpServer {
	return &httpServer{
		proc: proc,
		server: &http.Server{
			Addr:              config.listenAddressHTTP,
			ReadTimeout:       DefaultReadTimeout,
			ReadHeaderTimeout: DefaultReadHeaderTimeout,
			WriteTimeout:      DefaultWriteTimeout,
			IdleTimeout:       DefaultIdleTimeout,
			MaxHeaderBytes:    DefaultMaxHeaderBytes,
		},
	}
}

func (s *httpServer) Run(ctx context.Context) {
	s.server.Handler = s.router()
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func (s *httpServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *httpServer) router() http.Handler {
	r := chi.NewRouter()

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(30 * time.Second))
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post(gameJoinRoute, s.handleJoin)
	return r
}

func (s *httpServer) writeError(w http.ResponseWriter, text string, status int) {
	w.WriteHeader(status)

	_, err := w.Write([]byte(text))
	if err != nil {
		log.Err(err).Msg("error while send response")
	}
}
