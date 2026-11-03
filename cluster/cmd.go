package cluster

import (
	"strings"

	"go-redis/interface/databaseImp"
	"go-redis/interface/resp"
)

type (
	CmdLine  = [][]byte
	ExecFunc func(c *ClusterDatabase, client resp.Connection, db databaseImp.Database, args CmdLine) resp.Reply
)

var cmdTable = make(map[string]*command)

type command struct {
	executor ExecFunc
	arity    int
}

func RegisterCommand(name string, executor ExecFunc, arity int) {
	cmdName := strings.ToLower(name)

	cmdTable[cmdName] = &command{
		executor: executor,
		arity:    arity,
	}
}
