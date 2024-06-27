package storage

import (
	"context"
	"encoding/json"
	"ethereum-data-service/internal/model"
	"log"
	"time"

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

	if err := IdxBlockAndStore(ctx, rdb, &blockData, expiryTime); err != nil {
		return err
	}

	if err := IdxTxAndStore(ctx, rdb, &blockData, expiryTime); err != nil {
		return err
	}

	if err := IdxEventsAndStore(ctx, rdb, &blockData, expiryTime); err != nil {
		return err
	}

	log.Printf("Stored block %d in Redis\n", blockData.Block.Header.Number)
	return nil
}
