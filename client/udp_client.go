package client

import (
	"bytes"
	"context"
	"io"
	"net"
	"time"
)

const (
	// maxBufferSize specifies the size of the buffers that
	// are used to temporarily hold data from the UDP packets
	// that we send.
	maxBufferSize = 1024
	udpTimeout    = time.Millisecond * 50
)

type (
	UDPClient struct {
		addr *net.UDPAddr
		out  chan []byte
	}
)

func NewUDPClient(address string) (*UDPClient, error) {
	// Resolve the UDP address so that we can make use of DialUDP
	// with an actual IP and port instead of a name (in case a
	// hostname is specified).
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	return &UDPClient{
		out:  make(chan []byte),
		addr: raddr,
	}, nil
}

// client wraps the whole functionality of a UDP client that sends
// a message and waits for a response coming back from the server
// that it initially targetted.
func (c *UDPClient) Send(ctx context.Context, update []byte) ([]byte, error) {

	go func() {
		// Although we're not in a connection-oriented transport,
		// the act of `dialing` is analogous to the act of performing
		// a `connect(2)` syscall for a socket of type SOCK_DGRAM:
		// - it forces the underlying socket to only read and write
		//   to and from a specific remote address.
		conn, err := net.DialUDP("udp", nil, c.addr)
		if err != nil {
			println("error while make udp connection:", err.Error())
			return
		}

		// It is possible that this action blocks, although this
		// should only occur in very resource-intensive situations:
		// - when you've filled up the socket buffer and the OS
		//   can't dequeue the queue fast enough.
		_, err = io.Copy(conn, bytes.NewBuffer(update))
		if err != nil {
			conn.Close()
			println("error while copy data:", err.Error())
			return
		}

		//fmt.Printf("packet-written: bytes=%d\n", n)

		buffer := make([]byte, maxBufferSize)

		// Set a deadline for the ReadOperation so that we don't
		// wait forever for a server that might not respond on
		// a resonable amount of time.
		deadline := time.Now().Add(udpTimeout)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			conn.Close()
			println("error while set dealine:", err)
			return
		}

		nRead, _, err := conn.ReadFrom(buffer)
		if err != nil {
			conn.Close()
			println("error while read from URP", err)
			return
		}

		conn.Close()
		// fmt.Printf("packet-received: bytes=%d from=%s\n",
		// 	nRead, addr.String())

		c.out <- buffer[:nRead]
	}()

	res := <-c.out
	return res, nil
}
