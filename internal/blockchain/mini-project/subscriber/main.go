package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"project/common"

	"github.com/go-redis/redis/v8"
)

// Function to add a new block and trim the list
func addBlock(rdb *redis.Client, blockData common.BlockData) error {

	expiryTime := time.Duration(70 * time.Second)
	ctx := context.Background()

	// Index the block by its number for quick retrieval
	blockKey := fmt.Sprintf("block:%d", blockData.Block.Header.Number)
	blockDataJSON, err := json.Marshal(blockData.Block)
	if err != nil {
		return err
	}

	if err := rdb.Set(ctx, blockKey, blockDataJSON, expiryTime).Err(); err != nil {
		return err
	}

	// Index transactions and events with block's expiry
	for txHash, tx := range blockData.TransactionHashes {
		txKey := fmt.Sprintf("tx:%s", txHash)
		txJSON, err := json.Marshal(tx)
		if err != nil {
			return err
		}
		if err := rdb.Set(ctx, txKey, txJSON, expiryTime).Err(); err != nil {
			return err
		}
	}

	// Index address to retrieve all events associated with the address
	for _, events := range blockData.Events {
		for _, event := range events {
			eventJSON, err := json.Marshal(event)
			if err != nil {
				return err
			}
			addressKey := fmt.Sprintf("event:%s", event.Address)
			if err := rdb.Set(ctx, addressKey, eventJSON, expiryTime).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	// Initialize Redis client
	rdb := common.NewRedisClient()

	// Subscribe to the Redis channel
	subscriber := rdb.Subscribe(context.Background(), "new_blocks")
	defer subscriber.Close()

	// Channel to receive messages
	ch := subscriber.Channel()

	for msg := range ch {
		// Initialize the blockData variable
		var blockData common.BlockData

		// Deserialize the block data
		err := json.Unmarshal([]byte(msg.Payload), &blockData)
		if err != nil {
			log.Printf("Failed to unmarshal block data: %v", err)
			continue
		}

		// Add the block data to Redis
		err = addBlock(rdb, blockData)
		if err != nil {
			log.Printf("Failed to store block data in Redis: %v", err)
		} else {
			fmt.Printf("Stored block %d in Redis\n", blockData.Block.Header.Number)
		}
	}
}
