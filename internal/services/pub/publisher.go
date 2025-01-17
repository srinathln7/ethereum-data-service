package pub

import (
	"context"
	"ethereum-data-service/internal/client"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/model"
	"ethereum-data-service/pkg/util"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

// RunBlockNotifierSvc: Listens for new incoming blocks real-time from the Ethereum blockchain,
// extracts and formats the block as per the required format, and then publishes it to the Redis channel.
func RunBlockNotifierSvc(client *client.Client, cfg *config.Config, shutdown chan struct{}) {

	ethClient, rdb := client.ETH_WSS, client.REDIS

	// Create a common context instance
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals for graceful shutdown
	go util.HandleGracefulShutdown(cancel, shutdown)

	log.Println("Listening for new blocks from the Ethereum Blockchain...")
	err := listenForBlocks(ctx, ethClient, rdb, cfg, shutdown)
	if err != nil {
		// Prefered not to throw Fatalf to keep retrying in case of connection timeout
		log.Printf("error in block listener: %v", err)
	}
}

func listenForBlocks(ctx context.Context, ethClient *ethclient.Client, rdb *redis.Client, cfg *config.Config, shutdown chan struct{}) error {
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
			select {
			case <-shutdown:
				log.Println("Shutting down BlockNotifier service...")
				return nil
			default:
				err := handleNewHeader(ctx, ethClient, rdb, cfg, header)
				if err != nil {
					log.Printf("error handling new block header: %v", err)
				}
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

	// Publish the formatted block data to the Redis channel
	if err := rdb.Publish(ctx, cfg.REDIS_PUBSUB_CH, blockDataInBytes).Err(); err != nil {
		return err
	}

	log.Printf("Published new Block %d to Redis channel: %s\n", blockNumber, cfg.REDIS_PUBSUB_CH)
	return nil
}
