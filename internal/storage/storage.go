package storage

import (
	"context"
	"encoding/json"
	"ethereum-data-service/internal/model"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/redis/go-redis/v9"
)

var (
	BLOCK_PREFIX string = "block:"
	TX_PREFIX    string = "tx:"
	EVENT_PREFIX string = "event:"
)

func AddBlockDataToDB(ctx context.Context, rdb *redis.Client, payload []byte, expiryTime time.Duration) error {
	var blockData model.Data

	// Deserialize the block data
	err := json.Unmarshal(payload, &blockData)
	if err != nil {
		return err
	}

	// Index the block by its number
	if err := IndexBlock(ctx, rdb, &blockData, expiryTime); err != nil {
		return err
	}

	// Index transactions
	if err := IndexTransactions(ctx, rdb, &blockData, expiryTime); err != nil {
		return err
	}

	// Index events by address and transaction index
	if err := IndexEvents(ctx, rdb, &blockData, expiryTime); err != nil {
		return err
	}

	log.Printf("Stored block %d in Redis\n", blockData.Block.Header.Number)
	return nil
}

// IndexBlock: indexes the block data by its block number in Redis.
func IndexBlock(ctx context.Context, rdb *redis.Client, blockData *model.Data, expiryTime time.Duration) error {
	blockKey := fmt.Sprint(BLOCK_PREFIX, blockData.Block.Header.Number)
	blockDataJSON, err := json.Marshal(blockData.Block)
	if err != nil {
		return fmt.Errorf("error marshalling block data: %v", err)
	}
	if err := rdb.Set(ctx, blockKey, blockDataJSON, expiryTime).Err(); err != nil {
		return fmt.Errorf("error storing block data in Redis: %v", err)
	}
	return nil
}

// IndexTransactions: indexes each transaction data against its hash in Redis.
func IndexTransactions(ctx context.Context, rdb *redis.Client, blockData *model.Data, expiryTime time.Duration) error {
	for txHash, tx := range blockData.TransactionHashes {
		txKey := fmt.Sprint(TX_PREFIX, txHash)
		txJSON, err := json.Marshal(tx)
		if err != nil {
			return fmt.Errorf("error marshalling transaction %s: %v", txHash, err)
		}
		if err := rdb.Set(ctx, txKey, txJSON, expiryTime).Err(); err != nil {
			return fmt.Errorf("error storing transaction %s in Redis: %v", txHash, err)
		}
	}
	return nil
}

// IndexEvents: indexes each event by its address in Redis.
func IndexEvents(ctx context.Context, rdb *redis.Client, blockData *model.Data, expiryTime time.Duration) error {
	for _, events := range blockData.Events {
		for _, event := range events {
			eventJSON, err := json.Marshal(event)
			if err != nil {
				return fmt.Errorf("error marshalling event %+v: %v", event, err)
			}
			// Example key: `event:0x6000da47483062A0D734Ba3dc7576Ce6A0B645C4_0`
			addressKey := fmt.Sprint(EVENT_PREFIX, event.Address, "_", event.TxIndex)
			if err := rdb.Set(ctx, addressKey, eventJSON, expiryTime).Err(); err != nil {
				return fmt.Errorf("error storing event %s_%d in Redis: %v", event.Address, event.TxIndex, err)
			}
		}
	}
	return nil
}

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
