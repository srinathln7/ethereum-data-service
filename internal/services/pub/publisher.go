package pub

import (
	"context"
	"ethereum-data-service/internal/client"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/model"
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
func RunBlockNotifier(client *client.Client, cfg *config.Config) {
	log.Println("Starting Block Notification Service")

	ethClient, rdb := client.WSSETH, client.Redis

	// Create a common context instance
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals for graceful shutdown
	go handleGracefulShutdown(ethClient, rdb, cancel)

	log.Println("Listening for new blocks from the Ethereum Blockchain...")
	err := listenForBlocks(ctx, ethClient, rdb, cfg)
	if err != nil {
		log.Fatalf("error in block listener: %v", err)
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
				// DONOT: return error here but continue the service
				log.Printf("error handling new block header: %v", err)
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
	if err := rdb.Publish(ctx, cfg.REDIS_PUBSUB_CH, blockDataInBytes).Err(); err != nil {
		return err
	}

	log.Printf("Published new Block %d to Redis channel: %s\n", blockNumber, cfg.REDIS_PUBSUB_CH)
	return nil
}

func handleGracefulShutdown(ethClient *ethclient.Client, rdb *redis.Client, cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal
	sig := <-sigCh
	log.Printf("Received signal %v. Initating graceful shut down...", sig)

	// Close all client connections and trigger cancellation of context
	ethClient.Close()
	rdb.Close()
	cancel()
}
