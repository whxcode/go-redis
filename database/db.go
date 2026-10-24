package database

import (
	"bytes"
	"fmt"
	"log"

	"go-redis/datastruct/dict"
	"go-redis/interface/databaseImp"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

type (
	CmdLine  = [][]byte
	ExecFunc func(db *DB, args CmdLine) resp.Reply
)

type DB struct {
	index int
	data  dict.Dict
}

// AfterClientClose implements [databaseImp.Database].
func (d *DB) AfterClientClose(c resp.Connection) {
	// panic("unimplemented")
}

// Close implements [databaseImp.Database].
func (d *DB) Close() {
	log.Println("close db")
}

// Exec 执行命令
// eg: SET K V -> arity = 3
func (d *DB) Exec(client resp.Connection, cmd CmdLine) resp.Reply {
	// PING、SET、SETNX
	cmdName := string(bytes.ToLower(cmd[0]))
	cmdExecutor, ok := cmdTable[cmdName]
	// log.Printf("Exec cmdName=%s, cmdExecutor=%v, ok=%v", cmdName, cmdExecutor, ok)

	if !ok {
		return reply.MakeErrReply("ERR unknown command[ '" + cmdName + "']")
	}

	if !validateArity(cmdExecutor.arity, cmd) {
		return reply.MakeArgNumErrReply(fmt.Sprintf("%s,except[%d],actual[%d]", cmdName, cmdExecutor.arity, len(cmd)))
	}

	return cmdExecutor.executor(d, cmd[1:])
}

// SET K V -> arity = 3
// EXISTS k1,k2,... arity = -2 代表可用超过 2；参数是变长的
func validateArity(arity int, cmd CmdLine) bool {
	argNum := len(cmd)

	if arity >= 0 {
		return argNum == arity
	}

	return argNum >= -arity
}

func MakeDB() *DB {
	return &DB{
		index: 0,
		data:  dict.MakeSyncDict(),
	}
}

func (db *DB) GetEntity(key string) (*databaseImp.DataEntity, bool) {
	raw, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}

	entity, _ := raw.(*databaseImp.DataEntity)

	return entity, true
}

func (db *DB) PutEntity(key string, entity *databaseImp.DataEntity) int {
	return db.data.Put(key, entity)
}

// PutIfExists 存在则更新，不存在则不做任何操作
func (db *DB) PutIfExists(key string, entity *databaseImp.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

// PutIfAbsent 不存在则更新，存在则不做任何操作
func (db *DB) PutIfAbsent(key string, entity *databaseImp.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

func (db *DB) Remove(key string) int {
	return db.data.Remove(key)
}

// Removies 移动多个 key
func (db *DB) Removies(keys ...string) (deleted int) {
	for _, key := range keys {
		deleted += db.data.Remove(key)
	}

	return deleted
}

// Flush 清空数据库
func (db *DB) Flush() {
	db.data.Clear()
}

var _ databaseImp.Database = (*DB)(nil)
