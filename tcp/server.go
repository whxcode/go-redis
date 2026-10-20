package tcp

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go-redis/interface/tcp"
)

type Config struct {
	Address string
}

func ListenAndServerWithSignal(cfg *Config, handler tcp.Handler) error {
	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan

		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	log.Printf("server is listening at %s\n", cfg.Address)

	ListenAndServer(listener, cfg, handler, closeChan)

	return nil
}

func ListenAndServer(
	listener net.Listener,
	cfg *Config, handler tcp.Handler, closeChan <-chan struct{},
) error {
	// 监听系统关闭信号
	go func() {
		<-closeChan
		log.Println("server is shutting down...")
		listener.Close()
		_ = handler.Close()
		os.Exit(-1)
	}()

	defer func() {
		listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()

	var waitDone sync.WaitGroup

	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}

		log.Printf("client %s connected\n", conn.RemoteAddr().String())
		waitDone.Add(1)

		go func() {
			defer waitDone.Done()
			handler.Handle(ctx, conn)
		}()

	}

	waitDone.Wait()

	return nil
}
