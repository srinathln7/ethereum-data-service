package client

import (
	"ethereum-data-service/internal/config"
	"ethereum-data-service/pkg/enum"

	eth_err "ethereum-data-service/pkg/err"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	ETH_HTTPS *ethclient.Client
	ETH_WSS   *ethclient.Client
	REDIS     *redis.Client
}

// InitClient initializes and returns all clients
func InitClient() (*Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	httpsETHClient, err := newETHClient(cfg, enum.HTTPS)
	if err != nil {
		return nil, err
	}

	wssETHClient, err := newETHClient(cfg, enum.WSS)
	if err != nil {
		return nil, err
	}

	rdb, err := newRedisClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		ETH_HTTPS: httpsETHClient,
		ETH_WSS:   wssETHClient,
		REDIS:     rdb,
	}, nil
}

// newETHClient initializes and returns a new ETH client as per the specified protocol
func newETHClient(cfg *config.Config, protocol enum.Protocol) (*ethclient.Client, error) {
	var rpcURL string
	switch protocol {
	case enum.HTTPS:
		rpcURL = cfg.ETH_HTTPS_URL
	case enum.WSS:
		rpcURL = cfg.ETH_WSS_URL
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

// newRedisClient initializes and returns a new Redis client
func newRedisClient(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.REDIS_ADDR, // Redis server address
		DB:   cfg.REDIS_DB,   // Use default DB
	})

	return rdb, nil
}
