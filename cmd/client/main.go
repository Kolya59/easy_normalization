package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kolya59/easy_normalization/pkg/client"
)

func main() {
	// Graceful shutdown
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM)
	signal.Notify(sigint, syscall.SIGINT)
	done := make(chan interface{})

	go client.Start(done)

	// Wait interrupt signal
	<-sigint
	close(done)
}
