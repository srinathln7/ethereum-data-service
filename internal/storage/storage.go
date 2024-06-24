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

func AddBlockDataToDB(ctx context.Context, rdb *redis.Client, payload []byte, expiryTime time.Duration) error {

	var blockData model.Data

	// Deserialize the block data
	err := json.Unmarshal(payload, &blockData)
	if err != nil {
		return err
	}

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
			addressKey := fmt.Sprintf("events:%s", event.Address)
			if err := rdb.Set(ctx, addressKey, eventJSON, expiryTime).Err(); err != nil {
				return err
			}
		}
	}

	log.Printf("Stored block %d in Redis\n", blockData.Block.Header.Number)
	return nil
}

// GetEventsByAddress retrieves all events related to a specific Ethereum address from Redis.
// It takes a Redis client and an address as input, fetches the stored event data, and unmarshals it into a slice of Ethereum log events.
// Returns a slice of logs or an error if any operation fails.
func GetEventsByAddress(rdb *redis.Client, address string) ([]*types.Log, error) {
	data, err := rdb.Get(context.Background(), "events:"+address).Result()
	if err != nil {
		return nil, err
	}

	var events []*types.Log
	if err := json.Unmarshal([]byte(data), &events); err != nil {
		return nil, err
	}

	return events, nil
}

// GetBlockByNumber retrieves a specific Ethereum block by its number from Redis.
// It takes a Redis client and a block number as input, fetches the stored block data, and unmarshals it into a Block struct.
// Returns a pointer to the Block struct or an error if any operation fails.
func GetBlockByNumber(rdb *redis.Client, blockNumber string) (*model.Block, error) {
	data, err := rdb.Get(context.Background(), "block:"+blockNumber).Result()
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
	data, err := rdb.Get(context.Background(), "tx:"+txHash).Result()
	if err != nil {
		return nil, err
	}

	var tx types.Transaction
	if err := json.Unmarshal([]byte(data), &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}
