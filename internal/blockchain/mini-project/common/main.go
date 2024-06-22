package common

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-redis/redis/v8"
)

// NewRedisClient initializes and returns a new Redis client
func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Use default DB
	})
	return rdb
}

// Block: Header and Body
type Block struct {
	Header *types.Header
	Body   *types.Body
}

// BlockData represents the structure of the block for serialization
type BlockData struct {
	Block             Block
	TransactionHashes map[string]*types.Transaction
	Events            map[string][]*types.Log
}
