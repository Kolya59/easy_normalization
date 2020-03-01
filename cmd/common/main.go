package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kolya59/easy_normalization/cmd/server"
)

func main() {
	// Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGINT)
	done := make(chan interface{})

	go server.Start(done)

	// Wait interrupt signal
	select {
	case <-sigint:
		close(done)
	}
}
