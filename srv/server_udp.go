package main

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

type (
	udpServer struct {
		processor        Processor
		compressor       compressor
		listenAddressUDP string
	}
)

func newUDPServer(config config, processor Processor, compressor compressor) udpServer {
	return udpServer{
		processor:        processor,
		compressor:       compressor,
		listenAddressUDP: config.listenAddressUDP,
	}
}

func (s *udpServer) Run(ctx context.Context) {
	pc, err := net.ListenPacket("udp", s.listenAddressUDP)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buffer := make([]byte, maxBufferSize)
	for {
		if err := pc.SetReadDeadline(time.Now().Add(udpTimeout)); err != nil {
			continue
		}

		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			continue
		}

		raw := make([]byte, n)
		copy(raw, buffer[:n])

		go func() {
			state, err := s.decode(raw)
			if err != nil {
				log.Err(err).Msg("error while decode input data")
				state.Message = err.Error()
				state.MessageType = MessageTypeError
			}

			state, err = s.processor.HandleIncomingWorldState(context.Background(), state)
			if err != nil {
				log.Err(err).Msg("handler error")
				state.Message = err.Error()
				state.MessageType = MessageTypeError
				return
			}

			updatedState, err := s.encode(state)
			if err != nil {
				log.Err(err).Msg("error while encode state")
				state.Message = err.Error()
				state.MessageType = MessageTypeError
			}

			err = pc.SetWriteDeadline(time.Now().Add(udpTimeout))
			if err != nil {
				log.Err(err).Msg("error while set write deadline")
				return
			}

			_, err = pc.WriteTo(updatedState, addr)
			if err != nil {
				log.Err(err).Msg("error while write state into UDP")
				return
			}
		}()

		select {
		case <-ctx.Done():
			log.Info().Msg("stop UDP server: context is done")
			return
		default:
		}
	}
}

func (s *udpServer) decode(raw []byte) (State, error) {
	decompressed, err := s.compressor.Decompress(raw)
	if err != nil {
		return State{}, err
	}

	var state State
	if err := json.Unmarshal(decompressed, &state); err != nil {
		return State{}, err
	}

	return state, nil
}

func (s *udpServer) encode(state State) ([]byte, error) {
	updatedState, err := json.Marshal(state.getSendData())
	if err != nil {
		return nil, err
	}

	updatedState = s.compressor.Compress(updatedState)
	return updatedState, nil
}
