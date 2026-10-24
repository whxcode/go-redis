package aof

import (
	"log"
	"os"

	"go-redis/config"
	"go-redis/database"
)

type CmdLine = [][]byte

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type AofHandler struct {
	database    *database.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int
}

func MakeAofHandler(db *database.Database) (af *AofHandler, err error) {
	aofFile, err := os.OpenFile(config.Properties.AppendFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}

	af = &AofHandler{
		database:    db,
		aofFilename: config.Properties.AppendFilename,
		aofFile:     aofFile,
	}

	return af, nil
}

// AddAof 将命令添加到 AOF 队列中数据
func (a *AofHandler) AddAof(dbIndex int, cmd CmdLine) {
	if !config.Properties.AppendOnly || a.aofChan == nil {
		log.Printf("AOF is not enabled, skipping AddAof for command: %v", cmd)
		return
	}

	p := &payload{
		cmdLine: cmd,
		dbIndex: dbIndex,
	}

	a.aofChan <- p
}

// handleAof payload (set k,v) <- aofChan (洛盘)
func (a *AofHandler) handlerAof() {
}
