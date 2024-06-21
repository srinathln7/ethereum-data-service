package service

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetLatest50BlockNumbers returns an array containing the most recent 50 block numbers
func GetLatest50BlockNumbers(client *ethclient.Client) [50]uint64 {

	latestBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to get the latest block number: %v", err)
	}

	var blockNumbers [50]uint64
	start := uint64(0)
	if latestBlockNumber >= 50 {
		start = latestBlockNumber - 49
	}
	for i := 0; i < 50; i++ {
		if start+uint64(i) > latestBlockNumber {
			break
		}
		blockNumbers[i] = start + uint64(i)
	}
	return blockNumbers
}
