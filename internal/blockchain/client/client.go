package blockchain

import (
	"ethereum-data-service/internal/config"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumClient struct {
	client *ethclient.Client
}

func NewEthereumClient() (*EthereumClient, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	// Create an Ethereum client
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	return &EthereumClient{client: client}, nil
}
