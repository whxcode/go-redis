package client

import (
	"fmt"
	"os/exec"
)

type Client struct{}

func MakeClient(peer string) *Client {
	return &Client{}
}

func (client *Client) Start() {
	out, err := exec.Command("redis-cli", "-p", "6378", "ping").Output()
	if err != nil {
		fmt.Println("执行出错:", err)
		return
	}
	fmt.Println(string(out))
}

func (client *Client) Close() error {
	return nil
}
