// connection 表示websocket连接的封装，提供了对底层连接的管理和操作。
package connection

import (
	"net"
	"sync"
	"time"

	"go-redis/interface/resp"
	"go-redis/lib/sync/wait"
)

type Connection struct {
	conn     net.Conn
	wait     wait.Wait
	mu       sync.Mutex
	selected int
}

var _ resp.Connection = (*Connection)(nil)

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Close() error {
	c.wait.WaitWithTimeout(10 * time.Second)

	return c.conn.Close()
}

// GetDBIndex implements [resp.Connection].
func (c *Connection) GetDBIndex() int {
	return c.selected
}

// SelectDB implements [resp.Connection].
func (c *Connection) SelectDB(db int) {
	c.selected = db
}

// Write implements [resp.Connection].
func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	c.mu.Lock()
	c.wait.Add(1)
	defer c.mu.Unlock()
	defer c.wait.Done()

	n, err := c.conn.Write(bytes)
	if err != nil {
		return err
	}

	if n != len(bytes) {
		return &net.OpError{
			Op:   "write",
			Net:  "tcp",
			Addr: c.conn.RemoteAddr(),
			Err:  err,
		}
	}

	return nil
}

func MakeConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}
