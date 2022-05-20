package client

import (
	"bytes"
	"context"
	"game/model"
	"io"
	"net"
	"time"
)

const (
	maxBufferSize = 1024
	udpTimeout    = time.Millisecond * 15
)

type (
	UDP struct {
		conn *net.UDPConn
		addr *net.UDPAddr
		out  chan []byte
	}
)

func NewUDPClient(address string) *UDP {
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		panic(err)
	}

	client := &UDP{
		out:  make(chan []byte),
		addr: raddr,
	}

	client.connect()
	return client
}

func (c *UDP) connect() {
	conn, err := net.DialUDP("udp", nil, c.addr)
	if err != nil {
		panic(err)
	}
	c.conn = conn
}

func (c *UDP) Send(ctx context.Context, update []byte) ([]byte, error) {
	update = model.Compress(update)
	_, err := io.Copy(c.conn, bytes.NewBuffer(update))
	if err != nil {
		println("error while copy data:", err.Error())
		return nil, err
	}
	buffer := make([]byte, maxBufferSize)

	deadline := time.Now().Add(udpTimeout)
	err = c.conn.SetReadDeadline(deadline)
	if err != nil {
		println("connection read deadline:", err.Error())
		return nil, err
	}
	nRead, _, err := c.conn.ReadFrom(buffer)
	if err != nil {
		println("error while read from UDP", err)
		return nil, err
	}

	buffer, _ = model.Decompress(buffer[:nRead])
	return buffer, nil
}
