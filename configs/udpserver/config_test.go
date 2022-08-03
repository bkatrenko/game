package udpserver

import (
	"os"
	"testing"

	"github.com/bkatrenko/game/configs/udpserver"
	"github.com/stretchr/testify/assert"
)

const (
	testUDPAddress  = "udp://test.local"
	testHTTPAddress = "http://test.local"
)

func TestNewConfigFromEnv(t *testing.T) {
	err := os.Setenv(udpserver.UDPAddress, testUDPAddress)
	assert.Nil(t, err)
	err = os.Setenv(udpserver.HTTPAddress, testHTTPAddress)
	assert.Nil(t, err)

	config, err := udpserver.NewConfigFromEnv()
	assert.Nil(t, err)
	assert.Equal(t, testUDPAddress, config.ListenAddressUDP)
	assert.Equal(t, testHTTPAddress, config.ListenAddressHTTP)
}
