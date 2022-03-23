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
	maxBufferSize           = 1024
	udpTimeout              = time.Millisecond * 30
	maxConnectionLivingTime = time.Second * 30
)

type (
	UDPClient struct {
		conn      *net.UDPConn
		addr      *net.UDPAddr
		out       chan []byte
		connected time.Time
	}
)

func NewUDPClient(address string) (*UDPClient, error) {
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	client := &UDPClient{
		out:  make(chan []byte),
		addr: raddr,
	}

	client.makeConnection()
	return client, nil
}

func (c *UDPClient) makeConnection() {
	conn, err := net.DialUDP("udp", nil, c.addr)
	if err != nil {
		panic(err)
	}
	c.conn = conn
	c.connected = time.Now()
}

func (c *UDPClient) Send(ctx context.Context, update []byte) ([]byte, error) {
	// if time.Since(c.connected) > maxConnectionLivingTime {
	// 	fmt.Println("reconnect")
	// 	c.makeConnection()
	// }

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
