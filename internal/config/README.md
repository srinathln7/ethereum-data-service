# README.md

## Overview

The `config` package is responsible for loading and managing configuration settings for the application. These settings include various timeouts, URLs, and other parameters needed for connecting to services such as Ethereum nodes and Redis.

## Structs

### Config

Holds configuration settings for the application.

- **Fields**:
  - `DEFAULT_TIMEOUT time.Duration`: Default timeout for network requests.
  - `API_PORT string`: Port for the API server.
  - `ETH_HTTPS_URL string`: HTTPS URL for accessing the Ethereum network.
  - `ETH_WSS_URL string`: WebSocket URL for accessing the Ethereum network.
  - `REDIS_DB int`: Redis database number to use.
  - `REDIS_ADDR string`: Address of the Redis server.
  - `REDIS_PUBSUB_CH string`: Redis Pub/Sub channel name for messaging.
  - `REDIS_KEY_EXPIRY_TIME time.Duration`: Expiration time for keys stored in Redis and is calculated based on avg. ETH block time.
  - `NUM_BLOCKS_TO_SYNC int`: Number of recent blocks to sync during initialization.
  - `BOOTSTRAP_TIMEOUT time.Duration`: Time after which the bootstrap service exits itself gracefully.

## Functions

### LoadConfig

Loads configuration settings from environment variables.

- **Returns**:
  - `*Config`: A struct containing the loaded configuration settings.
  - `error`: An error if any configuration setting is missing or invalid.

- **Behavior**:
  1. Defines the required environment variable keys.
  2. Retrieves environment variable values using the `GetEnvMap` utility function.
  3. Converts string values to appropriate types (e.g., integers, durations).
  4. Initializes and returns a `Config` struct with the converted values.
  5. Returns an error if any required environment variable is missing or cannot be converted.

