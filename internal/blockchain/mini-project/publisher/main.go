package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"project/common"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis/v8"
)

func publishBlock(rdb *redis.Client, block *types.Block) error {
	// Convert types.Block to common.BlockData
	transactionHashes := make([]string, block.Transactions().Len())
	for idx, tx := range block.Transactions() {
		transactionHashes[idx] = tx.Hash().Hex()
	}

	blockData := common.BlockData{
		Number:            block.NumberU64(),
		Hash:              block.Hash().Hex(),
		TransactionHashes: transactionHashes,
	}

	// Marshal BlockData to JSON
	blockDataJSON, err := json.Marshal(blockData)
	if err != nil {
		return fmt.Errorf("failed to marshal block data: %w", err)
	}

	// Publish the block data to Redis channel
	if err := rdb.Publish(context.Background(), "new_blocks", blockDataJSON).Err(); err != nil {
		return err
	}
	return nil
}

func main() {
	// Connect to the Ethereum client using websocket
	client, err := ethclient.Dial("wss://mainnet.ethereum.validationcloud.io/v1/wss/JFt58zlN7gcQlLYnMZcOD75LpethJgD6Eq5nKOxC9F0")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Initialize Redis client
	rdb := common.NewRedisClient()

	// Create a channel to receive new headers
	headers := make(chan *types.Header)
	// Subscribe to new headers
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Failed to subscribe to new headers: %v", err)
	}

	// Listen for new headers
	fmt.Println("Listening for new blocks...")
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case header := <-headers:
			blockNumber := header.Number
			block, err := client.BlockByNumber(context.Background(), blockNumber)
			if err != nil {
				log.Fatalf("Failed to get block: %v", err)
			}
			if err := publishBlock(rdb, block); err != nil {
				log.Fatalf("Failed to publish block: %v", err)
			}
			log.Printf("Published new block: %d\n", block.NumberU64())
		}
	}
}
