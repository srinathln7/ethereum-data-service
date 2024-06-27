package storage

import (
	"context"
	"encoding/json"
	"ethereum-data-service/internal/model"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// IdxBlockAndStore: Indexes the block data by its block number and stores in Redis.
func IdxBlockAndStore(ctx context.Context, rdb *redis.Client, blockData *model.Data, expiryTime time.Duration) error {
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

// IdxTxAndStore: Indexes each transaction data against its hash and stores in Redis.
func IdxTxAndStore(ctx context.Context, rdb *redis.Client, blockData *model.Data, expiryTime time.Duration) error {
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

// IdxEventsAndStore: Indexes each event by its address, blocknumber, tx_hash, and tx_idx stores in Redis.
func IdxEventsAndStore(ctx context.Context, rdb *redis.Client, blockData *model.Data, expiryTime time.Duration) error {
	for txHash, events := range blockData.Events {
		for _, event := range events {
			eventJSON, err := json.Marshal(event)
			if err != nil {
				return fmt.Errorf("error marshalling event %+v: %v", event, err)
			}
			// For ex key:`event:0x6000da47483062A0D734Ba3dc7576Ce6A0B645C4_20167294_0xd59016cbaf7c580e83544ac5bd98584f7ec65b6984ddbd6a7647d6873c16f63a_234`
			// This is done to keep the key associated with every event unique
			addressKey := fmt.Sprint(EVENT_PREFIX, event.Address.Hex(), "_", event.BlockNumber, "_", txHash, "_", event.Index)
			if err := rdb.Set(ctx, addressKey, eventJSON, expiryTime).Err(); err != nil {
				return fmt.Errorf("error storing event %s_%d in Redis: %v", event.Address, event.TxIndex, err)
			}
		}
	}
	return nil
}
