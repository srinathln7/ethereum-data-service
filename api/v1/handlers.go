package v1

import (
	"ethereum-data-service/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// setupHandlers registers all endpoint handlers.
func setupHandlers(router *gin.Engine, rdb *redis.Client) {

	// Default home route
	router.GET("/", listRoutes(router))

	// Application specific
	router.GET("/events", getEvents(rdb))
	router.GET("/block", getBlock(rdb))
	router.GET("/tx", getTransaction(rdb))

	// Handle favicon.ico request without logging
	router.GET("/favicon.ico", handleFavicon)
}

// listRoutes handles the root endpoint to list all registered routes.
func listRoutes(router *gin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		routes := router.Routes()
		var routePaths []string
		for _, route := range routes {
			routePaths = append(routePaths, route.Path)
		}
		c.JSON(http.StatusOK, gin.H{"routes": routePaths})
	}
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
