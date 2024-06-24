package pubsub

import (
	"context"
	"ethereum-data-service/internal/client"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/model"
	"ethereum-data-service/pkg/enum"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

// RunBlockNotificationService: Listens for new incoming blocks from the Ethereum blockchain,
// extracts and formats the block as per the required format, and then publishes it to the Redis channel.
func RunBlockNotificationService() {
	log.Println("Starting Block Notification Service")

	// Load the config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Start a new ethereum client
	ethClient, err := client.NewETHClient(enum.WSS)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Start a new Redis client
	rdb, err := client.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	// Create a common context instance
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals for graceful shutdown
	go handleSignals(ethClient, rdb, cancel)

	log.Println("Listening for new blocks from the Ethereum Blockchain...")
	err = listenForBlocks(ctx, ethClient, rdb, cfg)
	if err != nil {
		log.Fatalf("Error in block listener: %v", err)
	}
}

func listenForBlocks(ctx context.Context, ethClient *ethclient.Client, rdb *redis.Client, cfg *config.Config) error {

	headers := make(chan *types.Header)
	sub, err := ethClient.SubscribeNewHead(ctx, headers)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return err
		case header := <-headers:
			err := handleNewHeader(ctx, ethClient, rdb, cfg, header)
			if err != nil {
				// DO NOT: return error here but continue the service
				log.Printf("Error handling new header: %v", err)
			}
		}
	}
}

func handleNewHeader(ctx context.Context, ethClient *ethclient.Client, rdb *redis.Client, cfg *config.Config, header *types.Header) error {
	blockNumber := header.Number
	block, err := ethClient.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return err
	}

	blockDataInBytes, err := model.FormatBlockData(ethClient, block)
	if err != nil {
		return err
	}

	// Publish the formated blockdata to the redis channel
	if err := rdb.Publish(ctx, cfg.RedisPubSubCh, blockDataInBytes).Err(); err != nil {
		return err
	}

	log.Printf("Published new block %d to Redis channel: %s\n", blockNumber, cfg.RedisPubSubCh)
	return nil
}

func handleSignals(ethClient *ethclient.Client, rdb *redis.Client, cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal
	sig := <-sigCh
	log.Printf("Received signal %v. Initating graceful shut down...", sig)

	// Trigger cancellation of context
	ethClient.Close()
	rdb.Close()
	cancel()
}
