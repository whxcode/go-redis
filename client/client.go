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

func (client *Client) Send(port string, originMsg []byte) resp.Reply {
	c := `printf '` + string(originMsg) + `' | nc -q 0 127.0.0.1 ` + port
	cmd := exec.Command("bash", "-c", c)
	out, err := cmd.CombinedOutput()
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
