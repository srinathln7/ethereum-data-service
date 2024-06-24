package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	v1 "ethereum-data-service/api/v1"
	"ethereum-data-service/internal/client"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/services/pub"
	"ethereum-data-service/internal/services/sub"
)

// Service represents a service that can be started and stopped.
type Service struct {
	Name   string
	Start  func(*sync.WaitGroup, chan struct{})
	Active bool
}

func main() {
	// Define command-line flags
	pubFlag := flag.Bool("pub", false, "Start publisher service")
	subFlag := flag.Bool("sub", false, "Start subscriber service")
	apiFlag := flag.Bool("api-server", false, "Start API server")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize clients
	client, err := client.InitClient()
	if err != nil {
		log.Fatalf("Failed to initialize clients: %v", err)
	}

	// Define services
	services := []Service{
		{
			Name:   "BlockNotification",
			Start:  func(wg *sync.WaitGroup, shutdown chan struct{}) { pub.RunBlockNotifier(client, cfg, shutdown) },
			Active: *pubFlag,
		},
		{
			Name:   "BlockSubscriber",
			Start:  func(wg *sync.WaitGroup, shutdown chan struct{}) { sub.RunBlockSubscriber(client.Redis, cfg, shutdown) },
			Active: *subFlag,
		},
		{
			Name:   "API Server",
			Start:  func(wg *sync.WaitGroup, shutdown chan struct{}) { v1.RunAPIServer(client.Redis, cfg, shutdown) },
			Active: *apiFlag,
		},
	}

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
			go func(s Service) {
				defer wg.Done()
				log.Printf("Started %s", s.Name)
				s.Start(&wg, shutdown)
			}(service)
		}
	}

	// Wait for termination signal
	sig := <-sigCh
	log.Printf("Received signal %v. Initiating graceful shutdown...", sig)

	// Signal graceful shutdown to all services
	close(shutdown)

	// Wait for goroutines to finish with a timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All services shut down gracefully.")
	case <-time.After(5 * time.Second):
		log.Println("Shutdown timed out after 5 seconds.")
	}

	log.Println("Shutdown complete.")
}
