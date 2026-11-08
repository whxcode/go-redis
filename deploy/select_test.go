package deploy

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestSelect(t *testing.T) {
	c1 := make(chan string, 1)
	c2 := make(chan string)

	go func() {
		time.Sleep(time.Second * 3)
		c1 <- "whx"
	}()

	go func() {
		time.Sleep(time.Second * 2)
		c := <-c2
		fmt.Println("customer c2...", c)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)

	defer cancel()
	cancel()

	select {
	case msg := <-c1:
		fmt.Println("获取msg", msg)
	case c2 <- "32":
		// fmt.Println("发送msg")
	case <-time.After(time.Second * 10):
		fmt.Println("超时")
	case <-ctx.Done():
		fmt.Println("ctx timeout", ctx.Err())
		// default:
		//	fmt.Println("default")

	}
}
