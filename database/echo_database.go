package database

import (
	"fmt"

	"go-redis/interface/databaseImp"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

type EchoDatabase struct{}

// AfterClientClose implements [databaseImp.Database].
func (e *EchoDatabase) AfterClientClose(c resp.Connection) {
	// panic("EchoDatabase unimplemented AfterClientClose")
}

// Close implements [databaseImp.Database].
func (e *EchoDatabase) Close() {
	// panic("EchoDatabase unimplemented Close")
}

// Exec implements [databaseImp.Database].
func (e *EchoDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	fmt.Println("EchoDatabase.Exec() called with args:", args)
	for _, arg := range args {
		fmt.Printf("%s、", string(arg))
	}
	fmt.Println()

	return reply.MakeMultiBulkReply(args)
}

func MakeEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

var _ databaseImp.Database = (*EchoDatabase)(nil)
