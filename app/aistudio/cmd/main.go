package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	server "github.com/vapusdata-ecosystem/vapusdata/aistudio/server"
)

func main() {
	// This is the main function to start the platform server
	grpcserver := server.GrpcServer()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start a goroutine that waits for a termination signal.
	go func() {
		<-sigChan
		log.Println("Received shutdown signal. Gracefully stopping the gRPC server...")
		server.Shutdown(grpcserver)
	}()
	grpcserver.Run()
}
