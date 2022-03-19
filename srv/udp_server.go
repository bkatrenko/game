package main

import (
	"context"
	"encoding/json"
	"fmt"
	"game/model"
	"net"
	"time"
)

type (
	server struct {
		proc              proc
		listenAddressUDP  string
		listenAddressHTTP string
	}
)

func newServer(config config, proc proc) server {
	return server{
		proc:              proc,
		listenAddressUDP:  config.listenAddressUDP,
		listenAddressHTTP: config.listenAddressHTTP,
	}
}

// server wraps all the UDP echo server functionality.
// ps.: the server is capable of answering to a single
// client at a time.
func (s *server) runUDPServer(ctx context.Context) error {
	// ListenPacket provides us a wrapper around ListenUDP so that
	// we don't need to call `net.ResolveUDPAddr` and then subsequentially
	// perform a `ListenUDP` with the UDP address.
	//
	// The returned value (PacketConn) is pretty much the same as the one
	// from ListenUDP (UDPConn) - the only difference is that `Packet*`
	// methods and interfaces are more broad, also covering `ip`.
	pc, err := net.ListenPacket("udp", s.listenAddressUDP)
	if err != nil {
		return err
	}

	// `Close`ing the packet "connection" means cleaning the data structures
	// allocated for holding information about the listening socket.
	defer pc.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	// Given that waiting for packets to arrive is blocking by nature and we want
	// to be able of canceling such action if desired, we do that in a separate
	// go routine.
	go func() {
		for {
			// By reading from the connection into the buffer, we block until there's
			// new content in the socket that we're listening for new packets.
			//
			// Whenever new packets arrive, `buffer` gets filled and we can continue
			// the execution.
			//
			// note.: `buffer` is not being reset between runs.
			//	  It's expected that only `n` reads are read from it whenever
			//	  inspecting its contents.
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			// fmt.Printf("packet-received: bytes=%d from=%s\n",
			// 	n, addr.String())

			var state model.State
			if err := json.Unmarshal(buffer[:n], &state); err != nil {
				println("error while unmarshal state:", err.Error())
				continue
			}

			state, err = s.proc.handle(state)
			if err != nil {
				state.Message = err.Error()
				state.MessageType = model.MessageTypeError
			}

			updatedStatem, err := json.Marshal(state)
			if err != nil {
				println("error while marshal state:", err.Error())
				continue
			}

			// Setting a deadline for the `write` operation allows us to not block
			// for longer than a specific timeout.
			//
			// In the case of a write operation, that'd mean waiting for the send
			// queue to be freed enough so that we are able to proceed.
			deadline := time.Now().Add(udpTimeout)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			// Write the packet's contents back to the client.
			n, err = pc.WriteTo(updatedStatem, addr)
			if err != nil {
				doneChan <- err
				return
			}

			//fmt.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
		println("error error while handle TCP:", err.Error())
	}

	return nil
}
