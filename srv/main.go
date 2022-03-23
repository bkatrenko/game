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
	udpTimeout    = time.Millisecond * 200
)

func main() {
	config, err := newConfigFromEnv()
	if err != nil {
		panic(err)
	}
	p := newProc()
	server := newServer(config, p)
	go p.startModifier()
	go server.runUDPServer(context.Background())
	server.runHTTPServer(context.Background())
}
