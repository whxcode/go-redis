package database

import (
	"strconv"

	"go-redis/aof"
	"go-redis/config"
	"go-redis/interface/databaseImp"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
	"go-redis/tools"
)

type Database struct {
	dbs       [16]*DB
	aof       *aof.AofHandler
	dbIndex   int
	collector bool
}

// AfterClientClose implements [databaseImp.Database].
func (d *Database) AfterClientClose(c resp.Connection) {
	// panic("unimplemented")
}

// Close implements [databaseImp.Database].
func (d *Database) Close() {
	// panic("unimplemented")
}

var whiteCmds tools.SyncSet[string] = tools.SyncSet[string]{}

func init() {
	whiteCmds.Store("set")
}

// Exec implements [databaseImp.Database].
func (d *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	cmdName := string(args[0])

	if cmdName == "select" {
		d.dbIndex, _ = strconv.Atoi(string(args[0]))
		return selectdb(client, args[1:])
	}

	if cmdName == "showselect" {
		return showSelectdb(client, args[1:])
	}

	dbIndex := d.dbIndex

	if client != nil {
		dbIndex = client.GetDBIndex()
	}

	db := d.dbs[dbIndex]

	// 记录落盘
	if d.collector && whiteCmds.Has(cmdName) {
		d.aof.AddAof(dbIndex, args)
	}

	return db.Exec(client, args)
}

var _ databaseImp.Database = (*Database)(nil)

func MakeDatabase() *Database {
	dbs := &Database{}
	for i := range 16 {
		dbs.dbs[i] = MakeDB()
	}

	if config.Properties.AppendOnly {
		aof, err := aof.MakeAofHandler(dbs)
		if err != nil {
			panic(err)
		}

		dbs.aof = aof
		dbs.collector = false
		dbs.aof.LoadAof()
		dbs.collector = true
	}

	return dbs
}

func selectdb(conn resp.Connection, args CmdLine) resp.Reply {
	if len(args) != 1 {
		return reply.MakeErrReply("ERR wrong number of arguments for 'select' command")
	}

	i, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeErrReply("ERR invalid DB index")
	}

	if i >= 16 {
		return reply.MakeErrReply("ERR DB index is out of range support [0,16)")
	}

	if conn != nil {
		conn.SelectDB(i)
	}

	return reply.MakeOkReply()
}

func showSelectdb(conn resp.Connection, args CmdLine) resp.Reply {
	return reply.MakeIntReply(conn.GetDBIndex())
}
