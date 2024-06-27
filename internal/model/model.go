package model

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Block represents an Ethereum block, containing both header and body information.
type Block struct {
	Header *types.Header `json:"header"` // Header contains metadata about the block.
	Body   *types.Body   `json:"body"`   // Body contains the transactions and uncles of the block.
}

// Data represents storage data in Redis DB for Ethereum-related information.
// It includes Ethereum block data, transaction hashes mapped to transactions,
// and event data mapped to specific event addresses. Block data byitself
// contains Tx hashes but as per the challenge description, we store it explicitly here
type Data struct {
	Block             Block                         `json:"block"`              // Block holds Ethereum block information.
	TransactionHashes map[string]*types.Transaction `json:"transaction_hashes"` // TransactionHashes maps transaction hashes to their corresponding transactions.
	Events            map[string][]*types.Log       `json:"events"`             // Events maps from each transaction hashes to lists of event logs.
}

// FormatBlockData: Extracts the data from Ethereum Block and format the data
// as per model.Data and then marshalls the result into bytes.
func FormatBlockData(client *ethclient.Client, block *types.Block) ([]byte, error) {

	blockData := Data{
		Block:             Block{Header: block.Header(), Body: block.Body()},
		TransactionHashes: make(map[string]*types.Transaction),
		Events:            make(map[string][]*types.Log),
	}

	// Get transaction hashes and all events related to each transaction in each block
	for _, tx := range block.Transactions() {
		blockData.TransactionHashes[tx.Hash().Hex()] = tx

		// Fetch events for each transaction
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return nil, err
		}

		// Store all events related to a particular transaction
		blockData.Events[receipt.TxHash.Hex()] = receipt.Logs

	}

	// Marshal BlockData to bytes
	blockDataJSON, err := json.Marshal(blockData)
	if err != nil {
		return nil, err
	}

	return blockDataJSON, nil
}
