package cmd

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	v1 "ethereum-data-service/api/v1"

	"ethereum-data-service/internal/client"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/services/bootstrapper"
	"ethereum-data-service/internal/services/pub"
	"ethereum-data-service/internal/services/sub"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Initialize clients
	clientInstance, err = client.InitClient()
	if err != nil {
		log.Fatalf("failed to initialize clients: %v", err)
	}
}

var RootCmd = &cobra.Command{
	Use:   "ethereum_api_service",
	Short: "vc-CLI",
	Run: func(cmd *cobra.Command, args []string) {
		color.HiCyan("************************ Welcome to the VC-ETHEREUM DATA API SERVICE CLI *****************")
		color.HiCyan("Please use any of the following sub-commands 'bootstrap', 'api-server', 'pub' or 'sub'")
		color.HiCyan("To start the BlockBootstrapper service: go run main.go bootstrap`")
		color.HiCyan("To start the BlockSubscription service: `go run main.go sub`")
		color.HiCyan("To start the BlockNotification service `go run main.go pub`")
		color.HiCyan("To start the HTTP API server: `go run main.go api-server`")
	},
}

func init() {
	cobra.OnInitialize(Init)

	RootCmd.AddCommand(bootstrapCmd)
	RootCmd.AddCommand(pubCmd)
	RootCmd.AddCommand(subCmd)
	RootCmd.AddCommand(apiServerCmd)
}

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Start BlockBootstrap service",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		shutdown := make(chan struct{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			bootstrapper.RunBootstrapSvc(clientInstance, cfg)
		}()
		handleShutdown(&wg, shutdown)
	},
}

var pubCmd = &cobra.Command{
	Use:   "pub",
	Short: "Start BlockNotification service",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		shutdown := make(chan struct{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			pub.RunBlockNotifierSvc(clientInstance, cfg, shutdown)
		}()
		handleShutdown(&wg, shutdown)
	},
}

var subCmd = &cobra.Command{
	Use:   "sub",
	Short: "Start BlockSubscriber service",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		shutdown := make(chan struct{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			sub.RunBlockSubscriberSvc(clientInstance.REDIS, cfg, shutdown)
		}()
		handleShutdown(&wg, shutdown)
	},
}

var apiServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "Start HTTP-API server",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		shutdown := make(chan struct{})
		wg.Add(1)
		go func() {
			defer wg.Done()
			v1.RunAPIServer(clientInstance.REDIS, cfg, shutdown)
		}()
		handleShutdown(&wg, shutdown)
	},
}

func handleShutdown(wg *sync.WaitGroup, shutdown chan struct{}) {
	// Setup a signal handler to capture interrupt and termination signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-done
	log.Printf("Received signal %v. Initiating graceful shutdown...", sig)

	// Signal graceful shutdown to all services
	close(shutdown)

	// Wait for goroutines to finish with a timeout
	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case <-doneChan:
		log.Println("All running services shut down gracefully.")
	case <-time.After(cfg.DEFAULT_TIMEOUT):
		log.Println("Shutdown timed out after 5 seconds.")
	}

	log.Println("Shutdown complete.")
}
