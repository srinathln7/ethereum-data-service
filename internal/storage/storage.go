package storage

import (
	"context"
	"encoding/json"
	"ethereum-data-service/internal/model"
	"fmt"
	"log"
	"time"

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
			addressKey := fmt.Sprintf("event:%s", event.Address)
			if err := rdb.Set(ctx, addressKey, eventJSON, expiryTime).Err(); err != nil {
				return err
			}
		}
	}

	log.Printf("Stored block %d in Redis\n", blockData.Block.Header.Number)
	return nil
}
