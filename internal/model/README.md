# README.md

## Overview

The `model` package defines the structure for Ethereum block data and provides functionality to format and serialize this data for storage and processing.

## Structs

### Block

Represents an Ethereum block, containing both header and body information.

- **Fields**:
  - `Header *types.Header`: Contains metadata about the block.
  - `Body *types.Body`: Contains the transactions and uncles of the block.

### Data

Represents storage data in Redis DB for Ethereum-related information. It includes Ethereum block data, transaction hashes mapped to transactions, and event data mapped to specific event addresses.

- **Fields**:
  - `Block Block`: Holds Ethereum block information.
  - `TransactionHashes map[string]*types.Transaction`: Maps transaction hashes to their corresponding transactions.
  - `Events map[string][]*types.Log`: Maps transaction hashes to lists of event logs.

## Functions

### FormatBlockData

Extracts the data from an Ethereum block, formats it as `model.Data`, and then marshals the result into bytes.

- **Parameters**:
  - `client *ethclient.Client`: The Ethereum client used to fetch transaction receipts.
  - `block *types.Block`: The Ethereum block to be formatted.

- **Returns**:
  - `[]byte`: The serialized block data.
  - `error`: An error if any operation fails.

- **Behavior**:
  1. Initializes a `Data` struct with block header and body.
  2. Iterates over transactions in the block to populate `TransactionHashes` and `Events`.
  3. Fetches transaction receipts to extract event logs.
  4. Marshals the `Data` struct into a JSON byte array.
  5. Returns the serialized data or an error if any step fails.

