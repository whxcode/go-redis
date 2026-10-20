package tcp

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (c *EchoClient) Close() error {
	c.Waiting.WaitWithTimeout(1 * time.Second)
	return c.Conn.Close()
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func (h *EchoHandler) Size() int {
	return 0
}

func (h *EchoHandler) Handle(ctx context.Context, conn net.Conn) {
	conn.Write([]byte("Welcome to the echo server!\n"))

	if h.closing.Get() {
		_ = conn.Close()
		return
	}

	client := &EchoClient{
		Conn: conn,
	}

	h.activeConn.Store(client, struct{}{})

	rwDer := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	// 以 \n 作为结束标志
	for {
		str, err := rwDer.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Printf("client %s closed the connection\n", conn.RemoteAddr().String())
				h.activeConn.Delete(client)
			} else {
				log.Printf("client %s disconnected\n", conn.RemoteAddr().String())
			}
			return

		}

		upperStr := strings.ToUpper(str)

		client.Waiting.Add(1)
		n, err := rwDer.WriteString(upperStr)
		rwDer.Flush()
		client.Waiting.Done()

		if err != nil {
			log.Printf("client %s write error: %v\n", conn.RemoteAddr().String(), err)
		}
		if n != len(upperStr) {
			log.Printf("client %s write error: wrote %d bytes, expected %d bytes\n", conn.RemoteAddr().String(), n, len(upperStr))
		}

	}
}

func (h *EchoHandler) Close() error {
	if h.closing.Get() {
		log.Printf("EchoHandler is already closing\n")
		return nil
	}

	h.closing.Set(true)

	h.activeConn.Range(func(key, _ interface{}) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})

	return nil
}
