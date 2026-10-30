package aof

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"

	"go-redis/config"
	"go-redis/interface/databaseImp"
	"go-redis/lib/utils"
	"go-redis/resp/reply"

	"go-redis/resp/parser"
)

type CmdLine = [][]byte

type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type AofHandler struct {
	database    databaseImp.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	currentDB   int
}

func MakeAofHandler(db databaseImp.Database) (a *AofHandler, err error) {
	aofFile, err := os.OpenFile(config.Properties.AppendFilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}

	a = &AofHandler{
		database:    db,
		aofFilename: config.Properties.AppendFilename,
		aofFile:     aofFile,
		aofChan:     make(chan *payload, 1024*10), // 10k
		currentDB:   -1,
	}

	// a.LoadAof()

	go a.handlerAof()

	return a, nil
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
	for payload := range a.aofChan {
		if payload.dbIndex != a.currentDB {
			args := reply.MakeMultiBulkReply(utils.ToCmdLine("select", strconv.Itoa(payload.dbIndex))).ToBytes()
			n, err := a.aofFile.Write(args)
			if err != nil {
				log.Printf("AOF write error: %v", err)
				continue
			}

			if n != len(args) {
				log.Printf("AOF write error: expected to write %d bytes, but wrote %d bytes", len(args), n)
				continue
			}

			a.currentDB = payload.dbIndex
		}

		args := reply.MakeMultiBulkReply(payload.cmdLine).ToBytes()
		n, err := a.aofFile.Write(args)
		if err != nil {
			log.Printf("AOF write error: %v", err)
			continue
		}

		if n != len(args) {
			log.Printf("AOF write error: expected to write %d bytes, but wrote %d bytes", len(args), n)
			continue
		}

		// 要把
	}
}

func (a *AofHandler) LoadAof() {
	// reader := bufio.NewReader(a.aofFile)
	ch := parser.ParserTream(a.aofFile)

	for payload := range ch {
		if payload.Err != nil {
			log.Printf("AOF load error: %v", payload.Err)
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

		a.database.Exec(nil, args)
	}
}
