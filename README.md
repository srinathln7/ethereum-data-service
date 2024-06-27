# VC-Ethereum Data-API-Service

The VC-Ethereum Data-API-Service provides essential APIs to retrieve block, transaction, and event data for the most recent 50 blocks from the Ethereum blockchain

## Objective

Refer to the [project objective](https://github.com/srinathln7/ethereum-data-service/blob/main/docs/CHALLENGE.md) and the project structure. This document outlines the purpose of the project and how the codebase is organized.

## Prerequisites

- [Go](https://go.dev/doc/install)
- [Redis](https://redis.io/)
- [Docker](https://docs.docker.com/engine/install/ubuntu/) (v20.10.21 or higher)
- [Docker Compose](https://docs.docker.com/compose/install/linux/) (v20.20.2 or higher)


## Design and Architecture

Before diving into the codebase, it is highly recommended to review the [design and architecture](https://github.com/srinathln7/ethereum-data-service/blob/main/DESIGN.md) document. This provides a high-level explanation of how various services, components, and modules interact.

For a deeper dive into each service, refer to the README files in the respective directories:
- [API Service](https://github.com/srinathln7/ethereum-data-service/tree/main/api/v1)
- [Bootstrapper Service](https://github.com/srinathln7/ethereum-data-service/tree/main/internal/services/bootstrapper)
- [Block Notification Service](https://github.com/srinathln7/ethereum-data-service/tree/main/internal/services/pub)
- [Block Subscriber Service](https://github.com/srinathln7/ethereum-data-service/tree/main/internal/services/sub)
- [Data Formatter](https://github.com/srinathln7/ethereum-data-service/tree/main/internal/model)
- [Redis Storage](https://github.com/srinathln7/ethereum-data-service/tree/main/internal/storage)


## Get Started

To build and start all required services, rename `sample.env` to `.env` and run:

```
# If Redis is already running on your localhost, stop it before running Redis in the container
sudo systemctl stop redis

make buildup
```

To stop all running services, run:

```
make builddown
```

To view logs in real-time for each service:

```sh
docker logs -f vc-api-server
docker logs -f vc-bootstrapper
docker logs -f vc-blocknotifier
docker logs -f vc-blocksubscriber
```

## API Endpoints

| ID    | Route                                      | Description                                                    | Avg. Resp Time (ms)  |
|-------|--------------------------------------------|----------------------------------------------------------------|----------------------|
| VC-00 | GET `/`                                    | List all routes                                                |                      |
| VC-01 | GET `/v1/blocks`                           | List all block numbers currently available in local data store |                      |
| VC-02 | GET `/v1/block?block_number=<block_number>`| Get block info associated with a given block number            |                      |
| VC-03 | GET `/v1/tx?tx_hash=<tx_hash>`             | Get transaction info associated with a given transaction hash  |                      |
| VC-04 | GET `/v1/events?address=<address>`         | Get all events associated with a particular address            |                      |

VC-02, VC-03, VC-04 all get their info from the local data store.

## Test

You can test the service using the following curl commands:

```sh
# VC-00: List all routes 
curl -X GET http://localhost:8080/ | jq 

# VC-01: List all block numbers currently available in local data store
curl -X GET http://localhost:8080/v1/blocks | jq

# VC-02: Get block info associated with a given block number 
curl -X GET http://localhost:8080/v1/block?block_number=20162001 | jq

# VC-03: Get transaction info associated with a given transaction hash
curl -X GET http://localhost:8080/v1/tx?tx_hash=0xec9057951284893e709fd5e4d57f76f1013145f12d6e366664c060b6f2baf559 | jq

# VC-04: Get all events associated with a particular address
curl -X GET "http://localhost:8080/v1/events?address=0x38AfDe1E098B0cA181E144A639339198c8CF3Efa" | jq
```

Alternatively, you can test the service in your browser. 

To view real-time analytics of our local data store, click [here](http://localhost:5540) to navigate to the Redis Insight GUI. Then, manually connect the database using the following details:

```
host: redis
port: 6379
```

This will enable real-time monitoring of memory consumption and the number of keys generated in our local data store.

## Observation and Analysis

Currently under review.

## Limitations and Future Improvements

Currently under review.

## References

- [Ethereum JSON-RPC Specification](https://ethereum.github.io/execution-apis/api-documentation/)
- [Redis Documentation](https://redis.io/docs/latest/)
