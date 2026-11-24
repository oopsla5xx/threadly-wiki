package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"threadly-api-v1/internal/app"
)

func main() {
	shutdownSignalContext, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	if err := app.Run(shutdownSignalContext); err != nil {
		stop()
		log.Fatalf("api: %v", err)
	}
}
