package client

import (
	"ethereum-data-service/internal/config"
	"ethereum-data-service/pkg/enum"

	eth_err "ethereum-data-service/pkg/err"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

func NewETHClient(protocol enum.Protocol) (*ethclient.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	var rpcURL string
	switch protocol {
	case enum.HTTPS:
		rpcURL = cfg.HTTPSURL
	case enum.WSS:
		rpcURL = cfg.WSSURL
	default:
		return nil, eth_err.ErrInvalidProtocol
	}

	// Create an Ethereum client based on the specified protocol
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewRedisClient initializes and returns a new Redis client
func NewRedisClient() (*redis.Client, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr, // Redis server address
		DB:   0,             // Use default DB
	})

	return rdb, nil
}
