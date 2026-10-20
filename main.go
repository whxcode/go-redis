package main

import (
	"fmt"
	"log"
	"os"

	"go-redis/config"
	"go-redis/tcp"
)

const configFile string = "redis.json"

var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6378,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)

	return !os.IsNotExist(err) && !info.IsDir()
}

func main() {
	if fileExists(configFile) {
		config.SetupConfig(configFile)
	} else {
		config.Properties = defaultProperties
	}

	h := &tcp.EchoHandler{}

	err := tcp.ListenAndServerWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, h)

	log.Printf("%v", err)
}
