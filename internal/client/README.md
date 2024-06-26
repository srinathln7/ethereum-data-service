# Client 

## Overview

The `client` package is responsible for initializing and managing connections to various services, including Ethereum nodes (both HTTPS and WSS) and Redis.

## Functionality

### InitClient

This function initializes and returns a `Client` struct containing connections to the Ethereum node (both HTTPS and WSS) and Redis.

- **Returns**:
  - `*Client`: A struct containing the initialized clients.
  - `error`: An error if any initialization step fails.

- **Behavior**:
  1. Loads configuration settings.
  2. Initializes the Ethereum HTTPS client.
  3. Initializes the Ethereum WSS client.
  4. Initializes the Redis client.
  5. Returns the initialized clients or an error.

### newETHClient

This function initializes and returns a new Ethereum client based on the specified protocol.

- **Parameters**:
  - `cfg *config.Config`: Configuration settings for the Ethereum client.
  - `protocol enum.Protocol`: The protocol to use (HTTPS or WSS).

- **Returns**:
  - `*ethclient.Client`: The initialized Ethereum client.
  - `error`: An error if the client initialization fails.

- **Behavior**:
  1. Determines the RPC URL based on the specified protocol.
  2. Creates an Ethereum client using the specified protocol.
  3. Returns the initialized client or an error.

### newRedisClient

This function initializes and returns a new Redis client.

- **Parameters**:
  - `cfg *config.Config`: Configuration settings for the Redis client.

- **Returns**:
  - `*redis.Client`: The initialized Redis client.
  - `error`: An error if the client initialization fails.

- **Behavior**:
  1. Creates a new Redis client using the provided configuration settings.
  2. Returns the initialized client or an error.

## Configuration

### config.Config

- `ETH_HTTPS_URL`: The URL for connecting to the Ethereum HTTPS endpoint.
- `ETH_WSS_URL`: The URL for connecting to the Ethereum WSS endpoint.
- `REDIS_ADDR`: The address of the Redis server.
- `REDIS_DB`: The Redis database to use.

## Dependencies

- `github.com/ethereum/go-ethereum/ethclient`: Ethereum client library.
- `github.com/redis/go-redis/v9`: Redis client library.
- `ethereum-data-service/internal/config`: Configuration loading and management.
- `ethereum-data-service/pkg/enum`: Enumerations for various protocols.
- `ethereum-data-service/pkg/err`: Custom error definitions.

