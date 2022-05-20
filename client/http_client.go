package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"game/model"
	"net/http"
)

const (
	contentType = "application/json"
	joinRoute   = "/game/join"
)

type (
	HTTP struct {
		addr string
	}
)

func NewHTTPlient(address string) *HTTP {
	return &HTTP{
		addr: address,
	}
}

func (c *HTTP) Join(joinRequest model.JoinGame) (model.State, error) {
	joinRequestBytes, err := json.Marshal(joinRequest)
	if err != nil {
		return model.State{}, fmt.Errorf("error while marshal join request: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", c.addr, joinRoute), contentType, bytes.NewBuffer(joinRequestBytes))
	if err != nil {
		return model.State{}, fmt.Errorf("error while send join request: %w", err)
	}

	var result model.State
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return model.State{}, fmt.Errorf("error while decode join response: %w", err)
	}

	return result, nil
}
