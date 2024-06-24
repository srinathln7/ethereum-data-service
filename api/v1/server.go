package v1

import (
	"context"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/storage"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RunAPIServer initializes and runs the API server with graceful shutdown.
func RunAPIServer(rdb *redis.Client, cfg *config.Config, shutdown <-chan struct{}) {

	// Set Gin mode to release for production
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router without default middleware
	router := gin.New()

	// Add only necessary middleware
	router.Use(gin.Logger(), gin.Recovery())

	// Define the endpoints and their handlers
	router.GET("/events", getEvents(rdb))
	router.GET("/block", getBlock(rdb))
	router.GET("/tx", getTransaction(rdb))

	srv := &http.Server{
		Addr:    ":" + cfg.API_PORT,
		Handler: router,
	}

	// Run server in a goroutine so it doesn't block
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

// getEvents handles the /events endpoint, retrieving events related to a specific address from Redis.
func getEvents(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		address := c.Query("address")
		if address == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "address query parameter is required"})
			return
		}

		events, err := storage.GetEventsByAddress(rdb, address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get events from Redis", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, events)
	}
}

// getBlock handles the /block endpoint, retrieving a block by its number from Redis.
func getBlock(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		blockNumber := c.Query("block_number")
		if blockNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "block_number query parameter is required"})
			return
		}

		block, err := storage.GetBlockByNumber(rdb, blockNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get block from Redis", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, block)
	}
}

// getTransaction handles the /transaction endpoint, retrieving a transaction by its hash from Redis.
func getTransaction(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		txHash := c.Query("tx_hash")
		if txHash == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tx_hash query parameter is required"})
			return
		}

		tx, err := storage.GetTransactionByHash(rdb, txHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get transaction from Redis", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, tx)
	}
}
