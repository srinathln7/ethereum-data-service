# Storage (Redis)

## Overview

The `storage` package provides functionality to store and retrieve Ethereum block data, transactions, and events in Redis. It includes methods to index and query these data types efficiently.

## Functionality

### AddBlockDataToDB

This function stores block data, transactions, and events in Redis.

- **Parameters**:
  - `ctx context.Context`: The context for managing cancellation.
  - `rdb *redis.Client`: The Redis client.
  - `payload []byte`: The block data to store.
  - `expiryTime time.Duration`: The expiry time for the stored data.

- **Behavior**:
  1. Deserializes the block data from JSON.
  2. Indexes the block by its number.
  3. Indexes the transactions by their hashes.
  4. Indexes the events by their addresses.
  5. Logs the successful storage of the block data.

### IdxBlockAndStore

This function indexes block data by its block number in Redis.

- **Parameters**:
  - `ctx context.Context`: The context for managing cancellation.
  - `rdb *redis.Client`: The Redis client.
  - `blockData *model.Data`: The block data to index.
  - `expiryTime time.Duration`: The expiry time for the stored data.

- **Behavior**:
  1. Creates a Redis key using the block number.
  2. Serializes the block data to JSON.
  3. Stores the serialized block data in Redis with the specified expiry time.

### IdxTxAndStore

This function indexes each transaction by its hash in Redis.

- **Parameters**:
  - `ctx context.Context`: The context for managing cancellation.
  - `rdb *redis.Client`: The Redis client.
  - `blockData *model.Data`: The block data containing transactions to index.
  - `expiryTime time.Duration`: The expiry time for the stored data.

- **Behavior**:
  1. Iterates over each transaction in the block data.
  2. Creates a Redis key using the transaction hash.
  3. Serializes the transaction data to JSON.
  4. Stores the serialized transaction data in Redis with the specified expiry time.

### IdxEventsAndStore

This function indexes each event by its address in Redis. Before indexing, all addresses are converted to lower case to eliminate any case sensitivity.

- **Parameters**:
  - `ctx context.Context`: The context for managing cancellation.
  - `rdb *redis.Client`: The Redis client.
  - `blockData *model.Data`: The block data containing events to index.
  - `expiryTime time.Duration`: The expiry time for the stored data.

- **Behavior**:
  1. Iterates over each event in the block data.
  2. Creates a Redis key using the event's address, block number, transaction hash, and index.
  3. Serializes the event data to JSON.
  4. Stores the serialized event data in Redis with the specified expiry time.

### GetEventsByAddress

This function retrieves all events related to a specific Ethereum address from Redis.

- **Parameters**:
  - `rdb *redis.Client`: The Redis client.
  - `address string`: The Ethereum address to retrieve events for.

- **Behavior**:
  1. Constructs a Redis key pattern using the address.
  2. Retrieves all keys matching the pattern from Redis.
  3. Iterates over the keys and retrieves the associated event data.
  4. Deserializes the event data from JSON.
  5. Returns a slice of logs representing the events.

### GetBlockByNumber

This function retrieves a specific Ethereum block by its number from Redis.

- **Parameters**:
  - `rdb *redis.Client`: The Redis client.
  - `blockNumber string`: The block number to retrieve.

- **Behavior**:
  1. Constructs the Redis key using the block number.
  2. Retrieves the stored block data from Redis.
  3. Deserializes the block data from JSON.
  4. Returns a pointer to the `model.Block` struct.

### GetTransactionByHash

This function retrieves a specific Ethereum transaction by its hash from Redis.

- **Parameters**:
  - `rdb *redis.Client`: The Redis client.
  - `txHash string`: The transaction hash to retrieve.

- **Behavior**:
  1. Constructs the Redis key using the transaction hash.
  2. Retrieves the stored transaction data from Redis.
  3. Deserializes the transaction data from JSON.
  4. Returns a pointer to the `types.Transaction` struct.

### GetAllBlockNumbers

This function retrieves all block numbers stored in Redis.

- **Parameters**:
  - `rdb *redis.Client`: The Redis client.

- **Behavior**:
  1. Constructs a Redis key pattern to fetch all block numbers.
  2. Retrieves all keys matching the pattern from Redis.
  3. Returns a slice of strings representing the block numbers.

## Configuration

### config.Config

- `REDIS_PUBSUB_CH`: The Redis channel for publishing block data.
- `REDIS_KEY_EXPIRY_TIME`: The expiry time for storing block data in Redis.


## Usage

To use the storage service, initialize the necessary Redis client and configuration, and call the appropriate functions to index and query block data, transactions, and events.