package database

import (
	"go-redis/interface/databaseImp"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

// get k1
func get(db *DB, args CmdLine) resp.Reply {
	entity, ok := db.GetEntity(string(args[0]))

	if !ok {
		return reply.MakeNullBulkReply()
	}

	bytes := entity.Data.([]byte)

	return reply.MakeBulkReply(bytes)
}

// set k1 v1
func set(db *DB, args CmdLine) resp.Reply {
	key := string(args[0])

	db.PutEntity(key, &databaseImp.DataEntity{
		Data: args[1],
	})

	return reply.MakeOkReply()
}

// setx k1 v1
func setx(db *DB, args CmdLine) resp.Reply {
	key := string(args[0])

	i := db.PutIfAbsent(key, &databaseImp.DataEntity{
		Data: args[1],
	})

	return reply.MakeIntReply(i)
}

// getset k v 获取旧值并设置新值
func getset(db *DB, args CmdLine) resp.Reply {
	key := string(args[0])

	ov, ok := db.GetEntity(key)

	db.PutEntity(key, &databaseImp.DataEntity{
		Data: args[1],
	})

	if !ok {
		return reply.MakeNullBulkReply()
	}

	return reply.MakeBulkReply(ov.Data.([]byte))
}

// strlen k 获取字符串长度
func strlen(db *DB, args CmdLine) resp.Reply {
	key := string(args[0])

	entity, ok := db.GetEntity(key)

	if !ok {
		return reply.MakeNullBulkReply()
	}
	v := entity.Data.([]byte)

	return reply.MakeIntReply(int(len(v)))
}

func init() {
	RegisterCommand("get", get, 2)
	RegisterCommand("set", set, 3)
	RegisterCommand("setx", setx, 3)
	RegisterCommand("getset", setx, 3)
	RegisterCommand("strlen", setx, 2)
}
