package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	configs "github.com/bkatrenko/game/configs/udpserver"
	"github.com/bkatrenko/game/pkg/udpserver"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const ()

func init() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func main() {
	config, err := configs.NewConfigFromEnv()
	if err != nil {
		panic(err)
	}

	log.Info().
		Str("http_address", config.ListenAddressHTTP).
		Str("udp_address", config.ListenAddressUDP).
		Msg("start hockey server")

	p := udpserver.NewProc()
	udpServer := udpserver.NewUDPServer(config.ListenAddressUDP, p, udpserver.NewCompressor())
	httpServer := udpserver.NewHTTPServer(config, p)

	ctx, cancelFunc := context.WithCancel(context.Background())

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
