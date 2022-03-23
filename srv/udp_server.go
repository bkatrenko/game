package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"
	//"net"
	//"time"
)

type (
	server struct {
		proc              *proc
		listenAddressUDP  string
		listenAddressHTTP string
	}
)

func newServer(config config, proc *proc) server {
	return server{
		proc:              proc,
		listenAddressUDP:  config.listenAddressUDP,
		listenAddressHTTP: config.listenAddressHTTP,
	}
}

func (s *server) runUDPServer(ctx context.Context) error {
	pc, err := net.ListenPacket("udp", s.listenAddressUDP)
	if err != nil {
		return err
	}

	defer pc.Close()
	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				fmt.Println("error read from")
				return
			}

			decompressed, err := Decompress(buffer[:n])
			if err != nil {
				panic(err)
			}

			var state State
			if err := json.Unmarshal(decompressed, &state); err != nil {
				println("error while unmarshal state:", err.Error())
				return
			}

			state, err = s.proc.handle(state)
			if err != nil {
				fmt.Println("handler error:", err.Error())
				state.Message = err.Error()
				state.MessageType = MessageTypeError
				return
			}

			updatedStatem, err := json.Marshal(state.getSendData())
			if err != nil {
				println("error while marshal state:", err.Error())
				return
			}

			updatedStatem = Compress(updatedStatem)

			deadline := time.Now().Add(udpTimeout)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				fmt.Println(err)
				return
			}

			_, err = pc.WriteTo(updatedStatem, addr)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}()

	<-ctx.Done()
	fmt.Println(ctx.Err())
	return nil
}
