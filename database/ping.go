package database

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

func Ping(db *DB, args CmdLine) resp.Reply {
	return reply.MakePongReply()
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
