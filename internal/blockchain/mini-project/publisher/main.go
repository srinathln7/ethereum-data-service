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

func publishBlock(client *ethclient.Client, rdb *redis.Client, block *types.Block) error {

	// Convert types.Block to common.BlockData
	blockData := common.BlockData{
		Block:             common.Block{Header: block.Header(), Body: block.Body()},
		TransactionHashes: make(map[string]*types.Transaction),
		Events:            make(map[string][]*types.Log),
	}

	// Get transaction hashes and all events related to each transaction in each block
	for _, tx := range block.Transactions() {
		blockData.TransactionHashes[tx.Hash().Hex()] = tx

		// Fetch events for each transaction
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return err
		}

		blockData.Events[tx.Hash().Hex()] = receipt.Logs
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
			if err := publishBlock(client, rdb, block); err != nil {
				log.Fatalf("Failed to publish block: %v", err)
			}
			log.Printf("Published new block: %d\n", blockNumber)
		}
	}
}
