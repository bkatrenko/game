package main

import (
	"context"
	"time"
)

const (
	// maxBufferSize specifies the size of the buffers that
	// are used to temporarily hold data from the UDP packets
	// that we receive.
	maxBufferSize = 1024
	udpTimeout    = time.Millisecond * 100
)

func main() {
	config, err := newConfigFromEnv()
	if err != nil {
		panic(err)
	}
	p := newProc()
	udpServer := newUDPServer(config, p, newCompressor())
	httpServer := newHTTPServer(config, p)

	go p.StartGameEngine()
	go udpServer.run(context.Background())

	httpServer.run(context.Background())
}
