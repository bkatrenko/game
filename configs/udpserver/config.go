package udpserver

import (
	"errors"
	"os"
)

const (
	// UDPAddress is an address that will be listened by UDP server part of the application,
	// should be not empty in any case
	UDPAddress = "UDP_ADDRESS"
	// HTTPAddress is an address where HTTP server will listen for a new request, mostly for
	// the health checks and join game requests
	HTTPAddress = "HTTP_ADDRESS"
)

// config responsible for storing necessary data that will be used to start an application:
// like host:port we listen, timeouts or any parsed environment variable we need to pass into the
// application
type Config struct {
	ListenAddressUDP  string
	ListenAddressHTTP string
}

// newConfigFromEnv will parse current environment variables and fill config struct with the data
// from there
func NewConfigFromEnv() (Config, error) {
	newConfig := Config{
		ListenAddressUDP:  os.Getenv(UDPAddress),
		ListenAddressHTTP: os.Getenv(HTTPAddress),
	}

	if newConfig.ListenAddressHTTP == "" {
		return Config{}, errors.New("empty HTTP address")
	}

	if newConfig.ListenAddressUDP == "" {
		return Config{}, errors.New("empty UDP address")
	}

	return newConfig, nil
}
