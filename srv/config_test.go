package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testUDPAddress  = "udp://test.local"
	testHTTPAddress = "http://test.local"
)

func TestNewConfigFromEnv(t *testing.T) {
	err := os.Setenv(UDPAddress, testUDPAddress)
	assert.Nil(t, err)
	err = os.Setenv(HTTPAddress, testHTTPAddress)
	assert.Nil(t, err)

	config, err := newConfigFromEnv()
	assert.Nil(t, err)
	assert.Equal(t, testUDPAddress, config.listenAddressUDP)
	assert.Equal(t, testHTTPAddress, config.listenAddressHTTP)
}
