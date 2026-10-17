package main

import (
	"time"

	"go-redis/lib/sync/wait"
)

func main() {
	w := &wait.Wait{}
	w.Add(1)

	w.WaitWithTimeout(time.Second * 2)

	time.Sleep(time.Second * 5)
	w.Add(-1)
	time.Sleep(time.Second * 3)
}
