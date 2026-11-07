package main

import (
	"flag"
	"fmt"
	"log"
	"runtime/debug"

	"go-redis/config"
	"go-redis/resp/handler"
	"go-redis/tcp"
)

func main() {
	okBytes := []byte("+OK\r\n")
	fmt.Println(okBytes)
	fmt.Println(string(okBytes))

	// 定义 -p 参数，默认值 6379，描述信息
	port := flag.Int("p", 6379, "Redis port")

	// 解析命令行参数
	flag.Parse()

	// 使用参数值（注意 port 是指针，需要解引用）

	config.Properties.Port = *port
	config.Properties.Self = fmt.Sprintf("127.0.0.1:%d", *port)

	fmt.Printf("使用端口: [%v]\n", config.Properties.Self)

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
