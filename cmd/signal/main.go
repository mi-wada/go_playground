package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HelloServer struct {
	inverval   time.Duration
	ctx        context.Context
	shutdowned chan struct{}
}

func NewHelloServer(interval time.Duration, ctx context.Context) *HelloServer {
	if ctx == nil {
		ctx = context.Background()
	}
	return &HelloServer{inverval: interval, ctx: ctx, shutdowned: make(chan struct{})}
}

func (s *HelloServer) Start() {
	ticker := time.NewTicker(s.inverval)
	for {
		select {
		case <-s.ctx.Done():
			log.Println("Shutting down...")
			time.Sleep(2 * time.Second)
			ticker.Stop()
			s.shutdowned <- struct{}{}
			return
		case <-ticker.C:
			fmt.Println("Hello")
		}
	}
}

func (s *HelloServer) Shutdown() <-chan struct{} {
	return s.shutdowned
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	srv := NewHelloServer(1*time.Second, ctx)
	go func() {
		srv.Start()
	}()

	<-ctx.Done()
	stop()
	<-srv.Shutdown()
}
