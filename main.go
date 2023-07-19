package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"lsp.com/server/server"
)

func main() {
	// Create a new LSP server
	server := server.NewLSPServer()

	// Run the server in second routine
	go server.Run(context.Background())

	// Wait for termination signal
	waitForTerminationSignal()
}

func waitForTerminationSignal() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	log.Println("Shutting down...")
}
