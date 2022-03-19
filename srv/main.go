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
	udpTimeout    = time.Millisecond * 50
)

func main() {
	config, err := newConfigFromEnv()
	if err != nil {
		panic(err)
	}

	server := newServer(config, newProc())
	go server.runUDPServer(context.Background())
	server.runHTTPServer(context.Background())
}
