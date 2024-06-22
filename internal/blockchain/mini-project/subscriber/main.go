package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"project/common"

	"github.com/go-redis/redis/v8"
)

// Function to add a new block and trim the list
func addBlock(rdb *redis.Client, blockData common.BlockData) error {
	// Marshal BlockData to JSON
	blockDataJSON, err := json.Marshal(blockData)
	if err != nil {
		return fmt.Errorf("failed to marshal block data: %w", err)
	}

	// Add the new block to the head of the list
	if err := rdb.LPush(context.Background(), "latest_blocks", blockDataJSON).Err(); err != nil {
		return err
	}
	// Trim the list to keep only the latest 50 blocks
	if err := rdb.LTrim(context.Background(), "latest_blocks", 0, 49).Err(); err != nil {
		return err
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
			fmt.Printf("Stored block %d in Redis\n", blockData.Number)
		}
	}
}
