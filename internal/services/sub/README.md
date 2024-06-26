# BlockSubscriber service

## Overview

The `sub` package provides a service to subscribe to a Redis channel, receive incoming block data, and store it in a database.

## Functionality

### RunBlockSubscriberSvc

This function initializes and runs the block subscriber service.

- **Parameters**:
  - `rdb *redis.Client`: The Redis client.
  - `cfg *config.Config`: Configuration settings for the block subscriber service.
  - `shutdown chan struct{}`: Channel to handle shutdown signals.

- **Behavior**:
  1. Creates a context for managing cancellation.
  2. Starts a goroutine to handle graceful shutdown.
  3. Logs the subscription to the Redis channel.
  4. Subscribes to the specified Redis channel.
  5. Listens for messages on the Redis channel:
     - Handles shutdown signals and closes the subscriber gracefully.
     - Stores the received block data in the database.

## Configuration

### config.Config

- `REDIS_PUBSUB_CH`: The Redis channel to subscribe to for receiving block data.
- `REDIS_KEY_EXPIRY_TIME`: The expiry time for storing block data in Redis.

