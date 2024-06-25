package v1

import (
	"context"
	"ethereum-data-service/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RunAPIServer: Initializes and runs the API server with graceful shutdown.
func RunAPIServer(rdb *redis.Client, cfg *config.Config, shutdown <-chan struct{}) {

	// Set Gin mode to release for production
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router without default middleware
	router := gin.New()

	// Add only necessary middleware
	router.Use(gin.Logger(), gin.Recovery())

	// Define the endpoints and their handlers
	setupHandlers(router, rdb)

	srv := &http.Server{
		Addr:    ":" + cfg.API_PORT,
		Handler: router,
	}

	// Run server in a goroutine so it doesn't block
	log.Printf("Listening on port:%s \n", cfg.API_PORT)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-shutdown
	log.Println("Shutting down API server...")

	// Create a deadline to wait for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt a graceful server shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("API Server forced to shutdown: %v", err)
	}

	log.Println("API server exited gracefully")
}
