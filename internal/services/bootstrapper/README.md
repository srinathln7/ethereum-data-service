# Bootstrapper

## Overview

The `bootstrapper` package provides a service to fetch the most recent blocks from the Ethereum blockchain through HTTPS RPC endpoints and store them in a Redis database. The service handles initialization, execution, and graceful shutdown.

## Functionality

### RunBootstrapSvc

This function initializes and runs the bootstrap service.

- **Parameters**:
  - `client *client.Client`: Contains the Ethereum and Redis clients.
  - `cfg *config.Config`: Configuration settings for the bootstrap service.

- **Behavior**:
  1. Initializes the Ethereum and Redis clients.
  2. Creates a context with a timeout based on the configuration.
  3. Logs the start time.
  4. Calls `loadRecentBlockData` to fetch and store the latest Ethereum blocks.
  5. Logs the total execution time.
  6. Shuts down the service gracefully.

### loadRecentBlockData

This function fetches the most recent Ethereum blocks and stores them in Redis.

- **Parameters**:
  - `ctx context.Context`: The context for managing cancellation and timeout.
  - `ethClient *ethclient.Client`: The Ethereum client.
  - `rdb *redis.Client`: The Redis client.
  - `cfg *config.Config`: Configuration settings.

- **Behavior**:
  1. Retrieves the latest block number from the Ethereum blockchain.
  2. Logs the number of blocks to be synchronized and the latest block height.
  3. Iterates over the specified number of recent blocks, starting from the oldest.
  4. For each block:
     - Retrieves the block data.
     - Formats the block data.
     - Stores the block data in Redis with the pre-defined expiry time.
  5. Logs the successful loading of blocks into Redis.


