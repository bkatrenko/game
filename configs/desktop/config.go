package desktop

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	UDPServerHostPort  = "UDP_SERVER_HOST_PORT"
	HTTPServerHostPort = "HTTP_SERVER_HOST_PORT"

	PlayerID     = "PLAYER_ID"
	GameID       = "GAME_ID"
	PlayerNumber = "PLAYER_NUMBER"
)

type Config struct {
	UDPServerHostPort,
	HTTPServerHostPort,
	PlayedID,
	GameID string
	PlayerNumber int8
}

func GetConfigFromEnv() (Config, error) {
	config, err := parseConfig()
	if err != nil {
		return Config{}, fmt.Errorf("error while parse config: %w", err)
	}

	return config, validateConfig(config)
}

func parseConfig() (Config, error) {
	playerNumber, err := strconv.ParseInt(os.Getenv(PlayerNumber), 10, 8)
	if err != nil {
		return Config{}, fmt.Errorf("invalid player number: %w", err)
	}

	return Config{
		UDPServerHostPort:  os.Getenv(UDPServerHostPort),
		HTTPServerHostPort: os.Getenv(HTTPServerHostPort),

		PlayedID:     os.Getenv(PlayerID),
		GameID:       os.Getenv(GameID),
		PlayerNumber: int8(playerNumber),
	}, nil
}

func validateConfig(config Config) error {
	if config.PlayedID == "" {
		return errors.New("player ID can't be empty")
	}

	if config.GameID == "" {
		return errors.New("game ID can't be empty")
	}

	if config.UDPServerHostPort == "" {
		return errors.New("bad UDP game server host | port")
	}

	if config.HTTPServerHostPort == "" {
		return errors.New("bad HTTP game server host | port")
	}

	if config.PlayerNumber != 0 && config.PlayerNumber != 1 {
		return errors.New("expect player number to be 0 or 1")
	}

	return nil
}
