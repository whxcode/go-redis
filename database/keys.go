package database

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"

	"github.com/IGLOU-EU/go-wildcard/v2"
)

// del key [key ...]
func del(db *DB, args CmdLine) resp.Reply {
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = string(arg)
	}

	return reply.MakeIntReply(db.Removies(keys...))
}

// exists key [key ...]
func exists(db *DB, args CmdLine) resp.Reply {
	i := 0

	for _, arg := range args {
		key := string(arg)
		if _, exists := db.GetEntity(key); exists {
			i++
		}
	}

	return reply.MakeIntReply(i)
}

// exists key [key ...]
func flushdb(db *DB, args CmdLine) resp.Reply {
	db.Flush()

	return reply.MakeOkReply()
}

// type key
func eType(db *DB, args CmdLine) resp.Reply {
	key := string(args[0])

	entity, exists := db.GetEntity(key)

	if !exists {
		return reply.MakeNullBulkReply()
	}

	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
		// case *Set:
		return reply.MakeStatusReply("set")
		// case *List:
		return reply.MakeStatusReply("list")
	default:
		return reply.MakeUnknownErrReply()
	}
}

// rename key newkey 交换
func rename(db *DB, args CmdLine) resp.Reply {
	src := string(args[0])
	dist := string(args[1])

	entity, exists := db.GetEntity(src)

	if !exists {
		return reply.MakeErrReply("no such key")
	}

	db.PutEntity(dist, entity)
	db.Removies(src)

	return reply.MakeOkReply()
}

// renamex key newkey 如果 newkey 已经存在，则不执行任何操作
func renamex(db *DB, args CmdLine) resp.Reply {
	src := string(args[0])
	dist := string(args[1])

	_, exists := db.GetEntity(dist)

	if exists {
		return reply.MakeIntReply(0)
	}

	entity, exists := db.GetEntity(src)

	if !exists {
		return reply.MakeErrReply("no such key")
	}

	db.PutEntity(dist, entity)
	db.Removies(src)

	return reply.MakeIntReply(1)
}

// keys *
func keys(db *DB, args CmdLine) resp.Reply {
	newArgs := make([][]byte, 0, 10)

	db.data.ForEach(func(key string, _ any) bool {
		if wildcard.MatchFromByte(args[0], []byte(key)) {
			newArgs = append(newArgs, []byte(key))
		}

		return true
	})

	return reply.MakeMultiBulkReply(newArgs)
}

func init() {
	RegisterCommand("del", del, -1)
	RegisterCommand("exists", exists, -2)
	RegisterCommand("flushdb", flushdb, 1)
	RegisterCommand("type", eType, 2)
	RegisterCommand("rename", rename, 3)
	RegisterCommand("renamex", renamex, 3)
	RegisterCommand("keys", keys, 2)
}
