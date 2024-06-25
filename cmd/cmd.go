package cmd

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	v1 "ethereum-data-service/api/v1"
	"ethereum-data-service/internal/client"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/services/bootstrapper"
	"ethereum-data-service/internal/services/pub"
	"ethereum-data-service/internal/services/sub"
)

var (
	clientInstance *client.Client
	cfg            *config.Config
)

// Init initializes the configuration and clients.
func Init() {
	var err error
	// Load configuration
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize clients
	clientInstance, err = client.InitClient()
	if err != nil {
		log.Fatalf("Failed to initialize clients: %v", err)
	}
}

// Service represents a service that can be started and stopped.
type Service struct {
	Name   string
	Start  func(*sync.WaitGroup, chan struct{})
	Active bool
}

// ParseFlags parses the command-line flags and returns the active services.
func ParseFlags() []Service {
	// Define command-line flags
	bootFlag := flag.Bool("bootstrap", false, "Start Bootstrap service")
	pubFlag := flag.Bool("pub", false, "Start Publisher service")
	subFlag := flag.Bool("sub", false, "Start Subscriber service")
	apiFlag := flag.Bool("api-server", false, "Start HTTP API server")
	flag.Parse()

	// Define services
	return []Service{
		{
			Name: "Bootstrapper",
			Start: func(wg *sync.WaitGroup, shutdown chan struct{}) {
				bootstrapper.RunBootstrapSvc(clientInstance, cfg)
			},
			Active: *bootFlag,
		},
		{
			Name: "BlockNotification",
			Start: func(wg *sync.WaitGroup, shutdown chan struct{}) {
				pub.RunBlockNotifierSvc(clientInstance, cfg, shutdown)
			},
			Active: *pubFlag,
		},
		{
			Name: "BlockSubscriber",
			Start: func(wg *sync.WaitGroup, shutdown chan struct{}) {
				sub.RunBlockSubscriberSvc(clientInstance.REDIS, cfg, shutdown)
			},
			Active: *subFlag,
		},
		{
			Name: "API Server",
			Start: func(wg *sync.WaitGroup, shutdown chan struct{}) {
				v1.RunAPIServer(clientInstance.REDIS, cfg, shutdown)
			},
			Active: *apiFlag,
		},
	}
}

// HandleShutdown handles the graceful shutdown process.
func HandleShutdown(wg *sync.WaitGroup, sigCh chan os.Signal, shutdown chan struct{}) {
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
	case <-time.After(cfg.DEFAULT_TIMEOUT):
		log.Println("Shutdown timed out after 5 seconds.")
	}

	log.Println("Shutdown complete.")
}
