package client

import (
	"bytes"
	"fmt"
	"os/exec"

	"go-redis/interface/resp"
	"go-redis/resp/parser"
)

type Client struct{}

func MakeClient(peer string) *Client {
	return &Client{}
}

func (client *Client) Start() []byte {
	out, err := exec.Command("redis-cli", "-p", "6378", "ping").Output()
	if err != nil {
		fmt.Println("执行出错:", err)
		return nil
	}

	return out
}

func (client *Client) Send(port string, args [][]byte) resp.Reply {
	cmd := []string{"-p", port}

	for _, arg := range args {
		cmd = append(cmd, string(arg))
	}

	out, err := exec.Command("redis-cli", cmd...).Output()

	fmt.Println("---->", string(out))
	if err != nil {
		fmt.Println("执行出错:", err)
		return nil
	}

	ch := parser.ParserTream(bytes.NewReader(out))

	select {
	case payload := <-ch:
		if payload.Err != nil {
			fmt.Println("解析出错:", payload.Err)
			return nil
		}
		return payload.Data
	}
}

func (client *Client) Close() error {
	return nil
}
