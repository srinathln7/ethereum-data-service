package storage

import (
	"context"
	"encoding/json"
	"ethereum-data-service/internal/model"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/redis/go-redis/v9"
)

// GetEventsByAddress retrieves all events related to a specific Ethereum address from Redis.
// It takes a Redis client and an address as input, fetches the stored event data, and unmarshals it into a slice of Ethereum log events.
// Returns a slice of logs or an error if any operation fails.
func GetEventsByAddress(rdb *redis.Client, address string) ([]*types.Log, error) {

	// Use the Redis wild card pattern
	keyPattern := fmt.Sprint(EVENT_PREFIX, address, "_*")

	ctx := context.Background()

	// Get all keys matching the pattern from Redis
	keys, err := rdb.Keys(ctx, keyPattern).Result()
	if err != nil {
		return nil, fmt.Errorf("error fetching keys from Redis: %v", err)
	}

	events := make([]*types.Log, len(keys))

	// Iterate through keys and retrieve events
	for idx, key := range keys {
		eventJSON, err := rdb.Get(ctx, key).Result()
		if err != nil {
			return nil, fmt.Errorf("error fetching event from Redis: %v", err)
		}

		var event types.Log
		if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
			return nil, fmt.Errorf("error unmarshalling event JSON: %v", err)
		}

		events[idx] = &event
	}

	return events, nil
}

// GetBlockByNumber retrieves a specific Ethereum block by its number from Redis.
// It takes a Redis client and a block number as input, fetches the stored block data, and unmarshals it into a Block struct.
// Returns a pointer to the Block struct or an error if any operation fails.
func GetBlockByNumber(rdb *redis.Client, blockNumber string) (*model.Block, error) {
	data, err := rdb.Get(context.Background(), BLOCK_PREFIX+blockNumber).Result()
	if err != nil {
		return nil, err
	}

	var block model.Block
	if err := json.Unmarshal([]byte(data), &block); err != nil {
		return nil, err
	}

	return &block, nil
}

// GetTransactionByHash retrieves a specific Ethereum transaction by its hash from Redis.
// It takes a Redis client and a transaction hash as input, fetches the stored transaction data, and unmarshals it into a Transaction struct.
// Returns a pointer to the Transaction struct or an error if any operation fails.
func GetTransactionByHash(rdb *redis.Client, txHash string) (*types.Transaction, error) {
	data, err := rdb.Get(context.Background(), TX_PREFIX+txHash).Result()
	if err != nil {
		return nil, err
	}

	var tx types.Transaction
	if err := json.Unmarshal([]byte(data), &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

// GetAllBlockNumbers: Retrieves all block numbers stored in Redis.
// It takes a Redis client as input, constructs a Redis key pattern to fetch all block numbers,
// and retrieves all keys matching the pattern from Redis.
// Returns a slice of strings representing block numbers or an error if any operation fails.
func GetAllBlockNumbers(rdb *redis.Client) ([]string, error) {
	// Use the Redis wild card pattern
	keyPattern := fmt.Sprint(BLOCK_PREFIX, "*")

	// Get all keys matching the pattern from Redis
	keys, err := rdb.Keys(context.Background(), keyPattern).Result()
	if err != nil {
		return nil, fmt.Errorf("error fetching keys from Redis: %v", err)
	}

	return keys, nil
}
