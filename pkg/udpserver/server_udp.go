package udpserver

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	// maxBufferSize specifies the size of the buffers that
	// are used to temporarily hold data from the UDP packets
	// that we receive.
	// For buffers it is always a bit challenging to choose the best fit for its size:
	// - It should be big enough to hols "every" message without an exceptions (otherwise
	// hard to debug issue will be created)
	// - Depends on buffer size/amount of them we should not affect performance in any way
	maxBufferSize = 1024
	// udpTimeout describe the time we will spend for writing and reading from the UDP
	// connection. For now it is only one value (so, equal for both reading and writing).
	// It is hardcoded inside UDP server to make simple configuration and not overload environment
	// variables with a lot of unused values that will be never changed.
	// However, it could be easily extracted and became configurable
	udpTimeout = time.Millisecond * 100
)

type (
	// UDPServer is responsible for receiving/responding UDP traffic.
	// It contains:
	// - processor: to perform a mutations on input state
	// - compressor: to decode/encode traffic to minimize the network latency
	// - listenAddressUDP: to know which port should be taken
	UDPServer struct {
		processor        Processor
		compressor       compressor
		listenAddressUDP string
	}
)

// newUDPServer is a simple constructor for udpServer structure
func NewUDPServer(udpAddress string, processor Processor, compressor compressor) UDPServer {
	return UDPServer{
		processor:        processor,
		compressor:       compressor,
		listenAddressUDP: udpAddress,
	}
}

// Run will start async UDP server that will be responsible for handling current world state,
// enrich it will an updated data and send to the client.
// Server also encoding and decoding the data.
func (s *UDPServer) Run(ctx context.Context) {
	pc, err := net.ListenPacket("udp", s.listenAddressUDP)
	if err != nil {
		panic(err)
	}
	defer pc.Close()

	buffer := make([]byte, maxBufferSize)
	for {
		if err := pc.SetReadDeadline(time.Now().Add(udpTimeout)); err != nil {
			log.Err(err).Msg("error while set read deadline")
			continue
		}

		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			// Here we use debug log while  this line gives an error every time application
			// throttling
			log.Debug().Str("error", err.Error()).Msg("error while read from connection")
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

// encode will convert input State into JSON format and compress it with zstd algorithm
// https://facebook.github.io/zstd/
func (s *UDPServer) encode(state State) ([]byte, error) {
	updatedState, err := json.Marshal(state.getSendData())
	if err != nil {
		return nil, err
	}

	updatedState = s.compressor.Compress(updatedState)
	return updatedState, nil
}

// decode do de-compress data from zstd encoding and unmarshal input data into
// the State structure
func (s *UDPServer) decode(raw []byte) (State, error) {
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
