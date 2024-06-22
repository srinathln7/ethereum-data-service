package common

import (
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

// BlockData represents the structure of the block for serialization
type BlockData struct {
	Number            uint64   `json:"number"`
	Hash              string   `json:"hash"`
	TransactionHashes []string `json:"transaction_hashes"`
}
