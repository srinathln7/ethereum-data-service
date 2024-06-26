# BlockNotification Service 

## Overview

The `pub` package provides a service to listen for new incoming blocks from the Ethereum blockchain in real-time, format the block data, and publish it to a Redis channel.

## Functionality

### RunBlockNotifierSvc

This function initializes and runs the block notifier service.

- **Parameters**:
  - `client *client.Client`: Contains the Ethereum and Redis clients.
  - `cfg *config.Config`: Configuration settings for the block notifier service.
  - `shutdown chan struct{}`: Channel to handle shutdown signals.

- **Behavior**:
  1. Initializes the WebSocket Ethereum client and Redis client.
  2. Creates a context for managing cancellation.
  3. Starts a goroutine to handle graceful shutdown.
  4. Logs the start of block listening.
  5. Calls `listenForBlocks` to listen for new blocks and handle them.

### listenForBlocks

This function listens for new Ethereum block headers and processes them.

- **Parameters**:
  - `ctx context.Context`: The context for managing cancellation.
  - `ethClient *ethclient.Client`: The Ethereum client.
  - `rdb *redis.Client`: The Redis client.
  - `cfg *config.Config`: Configuration settings.
  - `shutdown chan struct{}`: Channel to handle shutdown signals.

- **Behavior**:
  1. Subscribes to new block headers from the Ethereum blockchain.
  2. Processes each incoming block header:
     - Handles errors and shutdown signals.
     - Calls `handleNewHeader` to process and publish the block data.
  3. Unsubscribes from the block header subscription upon shutdown or error.

### handleNewHeader

This function processes a new block header and publishes the formatted block data to a Redis channel.

- **Parameters**:
  - `ctx context.Context`: The context for managing cancellation.
  - `ethClient *ethclient.Client`: The Ethereum client.
  - `rdb *redis.Client`: The Redis client.
  - `cfg *config.Config`: Configuration settings.
  - `header *types.Header`: The new block header.

- **Behavior**:
  1. Retrieves the block corresponding to the header.
  2. Formats the block data.
  3. Publishes the formatted block data to the specified Redis channel.
  4. Logs the successful publication of the block data.

## Configuration

### config.Config

- `REDIS_PUBSUB_CH`: The Redis channel for publishing block data.

