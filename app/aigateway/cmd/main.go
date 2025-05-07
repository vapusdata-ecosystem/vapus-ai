package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	pkgs "github.com/vapusdata-ecosystem/vapusai/aigateway/pkgs"
	"github.com/vapusdata-ecosystem/vapusai/aigateway/server"
	// "github.com/lucas-clemente/quic-go/http3"
)

func main() {
	// Initialize the fibr server for webapp
	aigwateway := server.NewAIGateway()
	sigChan := make(chan os.Signal, 1)
	err := aigwateway.Listen(fmt.Sprintf(":%d", pkgs.NetworkConfigManager.AIGateway.Port))
	if err != nil {
		pkgs.DmLogger.Err(err).Msg("error while starting the server")
	} else {
		pkgs.DmLogger.Info().Msg("server started successfully at port " + fmt.Sprintf(":%d", pkgs.NetworkConfigManager.AIGateway.Port))
	}
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal. Gracefully stopping the AI Gateway server...")
		server.Shutdown()
	}()
}
