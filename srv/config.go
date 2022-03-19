package main

import (
	"errors"
	"os"
)

const (
	UDPAddress  = "UDP_ADDRESS"
	HTTPAddress = "HTTP_ADDRESS"
)

type config struct {
	listenAddressUDP  string
	listenAddressHTTP string
}

func newConfigFromEnv() (config, error) {
	newConfig := config{
		listenAddressUDP:  os.Getenv(UDPAddress),
		listenAddressHTTP: os.Getenv(HTTPAddress),
	}

	if newConfig.listenAddressHTTP == "" {
		return config{}, errors.New("empty HTTP address")
	}

	if newConfig.listenAddressUDP == "" {
		return config{}, errors.New("empty UDP address")
	}

	return newConfig, nil
}
