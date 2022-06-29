package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const ()

func init() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	config, err := newConfigFromEnv()
	if err != nil {
		panic(err)
	}

	log.Info().
		Str("http_address", config.listenAddressHTTP).
		Str("udp_address", config.listenAddressUDP).
		Msg("start hockey server")

	p := newProc()
	udpServer := newUDPServer(config.listenAddressUDP, p, newCompressor())
	httpServer := newHTTPServer(config, p)

	ctx, cancelFunc := context.WithCancel(context.Background())

	//go p.StartGameEngine(ctx)
	go udpServer.Run(ctx)
	go httpServer.Run(context.Background())

	waitForInterruption()
	log.Info().Msg("got interruption signal")
	if err := httpServer.Stop(ctx); err != nil && err != http.ErrServerClosed {
		log.Info().Str("error", err.Error()).Msg("error while stop an HTTP server")
		panic(err)
	}
	log.Info().Msg("HTTP server closed")
	cancelFunc()
}

func waitForInterruption() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt
}
