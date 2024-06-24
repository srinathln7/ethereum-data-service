package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"ethereum-data-service/internal/client"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/services/pub"
	"ethereum-data-service/internal/services/sub"
)

func main() {
	// Define command-line flags
	pubFlag := flag.Bool("pub", false, "Start publisher service")
	subFlag := flag.Bool("sub", false, "Start subscriber service")
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

	// Create a channel to receive OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Channel to signal graceful shutdown
	shutdown := make(chan struct{})

	// WaitGroup to wait for goroutines to finish
	var wg sync.WaitGroup

	// Start services based on flags
	if *pubFlag {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pub.RunBlockNotifier(client, cfg, shutdown)
		}()
	}

	if *subFlag {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub.RunBlockSubscriber(client.Redis, cfg, shutdown)
		}()
	}

	// Wait for termination signal
	sig := <-sigCh
	log.Printf("Received signal %v. Initiating graceful shutdown...", sig)

	// Signal graceful shutdown to both services
	close(shutdown)

	// Wait for goroutines to finish with a timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All goroutines finished successfully.")
	case <-time.After(5 * time.Second):
		log.Println("Shutdown timed out after 5 seconds.")
	}

	log.Println("Shutdown complete.")
}
