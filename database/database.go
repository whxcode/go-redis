package database

import (
	"strconv"

	"go-redis/interface/databaseImp"
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

type Database struct {
	dbs [16]*DB
}

// AfterClientClose implements [databaseImp.Database].
func (d *Database) AfterClientClose(c resp.Connection) {
	// panic("unimplemented")
}

// Close implements [databaseImp.Database].
func (d *Database) Close() {
	// panic("unimplemented")
}

// Exec implements [databaseImp.Database].
func (d *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	cmdName := string(args[0])

	if cmdName == "select" {
		return selectdb(client, args[1:])
	}

	if cmdName == "showselect" {
		return showSelectdb(client, args[1:])
	}

	db := d.dbs[client.GetDBIndex()]

	return db.Exec(client, args)
}

var _ databaseImp.Database = (*Database)(nil)

func MakeDatabase() *Database {
	dbs := &Database{}
	for i := range 16 {
		dbs.dbs[i] = MakeDB()
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

	conn.SelectDB(i)

	return reply.MakeOkReply()
}

func showSelectdb(conn resp.Connection, args CmdLine) resp.Reply {
	return reply.MakeIntReply(conn.GetDBIndex())
}
