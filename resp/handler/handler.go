package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strings"

	"go-redis/database"
	"go-redis/interface/databaseImp"
	"go-redis/interface/tcp"
	"go-redis/lib/sync/atomic"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"go-redis/tools"
)

type RespHandler struct {
	activeConn tools.SyncSet[*connection.Connection]
	db         databaseImp.Database
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	// var db databaseImp.Database = database.MakeEchoDatabase()
	var db databaseImp.Database = database.MakeDatabase()

	return &RespHandler{
		db: db,
	}
}

// Close implements [tcp.Handler].
func (h *RespHandler) Close() error {
	log.Println("handler shutting down...")
	h.closing.Set(true)
	h.activeConn.Range(func(conn *connection.Connection) bool {
		_ = conn.Close()
		return true
	})

	h.db.Close()

	return nil
}

// Handle implements [tcp.Handler].
func (h *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if h.closing.Get() {
		_ = conn.Close()
		return
	}

	client := connection.MakeConn(conn)

	h.activeConn.Store(client)

	ch := parser.ParserTream(conn)

	for payload := range ch {
		// err
		if payload.Err != nil {
			if payload.Err == io.EOF ||
				payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				h.closeClient(client)
				return
			}

			// protocol error
			errReply := reply.MakeStandardErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				h.closeClient(client)
				log.Println("write error, close the connection")
				return
			}

			continue
		}
		// nil exec
		if payload.Data == nil {
			continue
		}

		var args [][]byte = make([][]byte, 0)

		ref := reflect.TypeOf(payload.Data)

		switch payload.Data.(type) {
		case *reply.StatusReply:
			statusReply := payload.Data.(*reply.StatusReply)
			args = append(args, []byte(statusReply.Status))
		case *reply.MultiBulkReply:
			{
				args = payload.Data.(*reply.MultiBulkReply).Args
			}
		default:
			fmt.Printf("payload.Data type: %v\n", ref)
		}

		result := h.db.Exec(client, args)

		if result != nil {
			client.Write(result.ToBytes())
		} else {
			client.Write(reply.MakeUnknownErrReply().ToBytes())
		}
	}
}

func (h *RespHandler) closeClient(client *connection.Connection) {
	fmt.Printf("connection closed: %s\n", client.RemoteAddr().String())
	client.Close()
	h.db.AfterClientClose(client)
	h.activeConn.Delete(client)
}

var _ tcp.Handler = &RespHandler{}
