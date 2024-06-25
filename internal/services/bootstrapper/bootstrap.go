package bootstrapper

import (
	"context"
	"ethereum-data-service/internal/client"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/model"
	"ethereum-data-service/internal/storage"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

// RunBootstrapSvc initializes the Bootstrap Service, which fetches the most recent blocks from Ethereum and stores them in Redis.
// It creates a context for the operations, handles OS signals for graceful shutdown, calculates execution time, and shuts down automatically after completion.
func RunBootstrapSvc(client *client.Client, cfg *config.Config) {
	ethClient, rdb := client.ETH_HTTPS, client.REDIS

	// Create a common context instance with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.BOOTSTRAP_TIMEOUT)
	defer cancel()

	startTime := time.Now()

	err := loadRecentBlockData(ctx, ethClient, rdb, cfg)
	if err != nil {
		log.Printf("error running bootstrapper service: %v", err)
	}

	// Calculate and log total execution time
	totalTime := time.Since(startTime)
	log.Printf("Bootstrapper successfully completed in %s", totalTime)

	log.Println("Shutting down bootstraper gracefully...")
	os.Exit(0)

}

// loadRecentBlockData fetches the most recent blocks from Ethereum and loads them into the Redis server.
// It retrieves the latest block number, iterates over the number of blocks to sync, and stores each block in Redis with a calculated expiration time.
func loadRecentBlockData(ctx context.Context, ethClient *ethclient.Client, rdb *redis.Client, cfg *config.Config) error {

	latestBlockNumber, err := ethClient.BlockNumber(ctx)
	if err != nil {
		return err
	}
	latestBlockBigInt := new(big.Int).SetUint64(latestBlockNumber)

	log.Printf("Fetching the last %d Ethereum blocks (least-to-most recent) and loading them to Redis...\n", cfg.NUM_BLOCKS_TO_SYNC)
	log.Printf("Latest Block height is: %d \n", latestBlockBigInt.Int64())

	// Start loading from the oldest block i.e. `latestBlock - 50`
	n := cfg.NUM_BLOCKS_TO_SYNC
	for i := n; i >= 0; i-- {
		blockNumber := new(big.Int).Sub(latestBlockBigInt, big.NewInt(int64(i)))
		block, err := ethClient.BlockByNumber(ctx, blockNumber)
		if err != nil {
			return err
		}

		blockDataInBytes, err := model.FormatBlockData(ethClient, block)
		if err != nil {
			return err
		}

		// Calculate expiry time for each block data
		// For example: Assume blocks are [B0....B49]. B0 data will expire after the first 13 seconds, B1 after the next 13 seconds, and so on.
		// On average, it takes 12 seconds to produce a new block in Ethereum. We set the expiry to 13 seconds to provide a safety margin.
		expiryTime := time.Duration((n - i + 1)) * cfg.ETH_AVG_BLOCK_TIME

		err = storage.AddBlockDataToDB(ctx, rdb, blockDataInBytes, expiryTime)
		if err != nil {
			return err
		}
	}

	log.Printf("Successfully loaded %d blocks to Redis", cfg.NUM_BLOCKS_TO_SYNC)
	return nil
}
