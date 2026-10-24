package main

import (
	"fmt"
	"log"
	"runtime/debug"

	"go-redis/config"
	"go-redis/resp/handler"
	"go-redis/tcp"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			// %v 打印错误原因，%s 打印完整的堆栈信息（包含文件和行号）
			log.Printf("捕获到 Panic: %v\n堆栈信息如下:\n%s", r, debug.Stack())
		}
	}()

	h := handler.MakeHandler()

	err := tcp.ListenAndServerWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, h)

	log.Printf("%v", err)
}
