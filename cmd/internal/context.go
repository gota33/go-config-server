package internal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func NewAppContext() (ctx context.Context, cancel func()) {
	ctx, cancel = context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Try exit...")
		cancel()

		<-c
		fmt.Println("Force exit")
		os.Exit(0)
	}()
	return
}
