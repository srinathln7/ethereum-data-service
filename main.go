package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"ethereum-data-service/cmd"
)

func init() {
	// Initialize configuration and clients
	cmd.Init()
}

func main() {
	// Parse command-line flags and get active services
	services := cmd.ParseFlags()

	// Create a channel to receive OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal graceful shutdown
	shutdown := make(chan struct{})

	// WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup

	// Start active services
	for _, service := range services {
		if service.Active {
			wg.Add(1)
			go func(s cmd.Service) {
				defer wg.Done()
				log.Printf("Started %s service", s.Name)
				s.Start(&wg, shutdown)
			}(service)
		}
	}

	// Handle graceful shutdown
	cmd.HandleShutdown(&wg, sigCh, shutdown)
}
