package database

import "strings"

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
