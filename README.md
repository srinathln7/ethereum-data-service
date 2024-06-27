# VC-Ethereum Data-API-Service

VC-Ethereum Data-API-Service stores block info, transaction hashes, and events for the latest 50 blocks and provides API endpoints for querying events by address, blocks by number, and transactions by hash. Please refer to the [project objective](https://github.com/srinathln7/ethereum-data-service/blob/main/docs/CHALLENGE.md) here which outlines the goal of the project and how the codebase is organized.

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
- [Redis Storage](https://github.com/srinathln7/ethereum-data-service/tree/main/internal/storage)
- [Data Formatter](https://github.com/srinathln7/ethereum-data-service/tree/main/internal/model)


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

```

docker logs -f vc-api-server

docker logs -f vc-bootstrapper

docker logs -f vc-blocknotifier

docker logs -f vc-blocksubscriber
```

## API Endpoints

|   ID     | Route                                      | Description                                                    | Avg. Resp Time   |
|----------|--------------------------------------------|----------------------------------------------------------------|------------------|
|  VC-00   | GET `/`                                    | List all routes                                                |     40.44 µs     |
|  VC-01   | GET `/v1/blocks`                           | List all block numbers currently available in local data store |     9.13 ms      |
|  VC-02   | GET `/v1/block?block_number=<block_number>`| Get block info associated with a given block number            |     48.17 ms     |
|  VC-03   | GET `/v1/tx?tx_hash=<tx_hash>`             | Get transaction info associated with a given transaction hash  |     631.97 µs    |
|  VC-04   | GET `/v1/events?address=<address>`         | Get all events associated with a particular address            |     187.51 ms    |

`VC-02`, `VC-03`, `VC-04` all get their info from the local data store. 

Please note: When querying `VC-04` with a widely used contract address such as `0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48` for Circle USDC Token, which can potentially involve fetching thousands of events, there may be a slight delay in response time, approaching close to a second. However, despite occasional delays, the average response time for `VC-04` remains around 200ms.

## Test

You can test the service using the following curl commands:

```sh
# VC-00: List all routes 
curl -X GET http://localhost:8080/ | jq 

# VC-01: List all block numbers currently available in local data store
curl -X GET http://localhost:8080/v1/blocks | jq

# VC-02: Get block info associated with a given block number 
curl -X GET http://localhost:8080/v1/block?block_number=<$BLOCK_NUMBER> | jq

# VC-03: Get transaction info associated with a given transaction hash
curl -X GET http://localhost:8080/v1/tx?tx_hash=<$TX_HASH> | jq

# VC-04: Get all events associated with a particular address
curl -X GET "http://localhost:8080/v1/events?address=<$ADDR>" | jq
```

Alternatively, you can test the service in your browser. 

## Observation and Analysis

To view real-time analytics of our local data store, navigate to http://localhost:5540 in your browser to access the Redis Insight GUI. Then, manually connect the database by clicking `Add connection details manually` and using the following details 

```
host: redis
port: 6379
```

After continuously running the application for several hours, we monitored Redis using the Redis-Insight tool and observed that Redis peak memory consumption reached approximately 50MB. The number of indexed keys maintained averaged between 30,000 to 40,000, with Redis utilizing less than 1% of CPU resources. This performance is well within Redis's capabilities, as outlined in the official [Redis FAQ](https://redis.io/docs/latest/develop/get-started/faq/#:~:text=Redis%20can%20handle%20up%20to), which states Redis can manage up to 2^32 keys and has been tested to handle at least 250 million keys per instance.

This validates that Redis is an optimal choice for our current project requirements.

## Limitations and Future Improvements

In our current architecture, each service operates with a single instance, making each service vulnerable to being a single point of failure. This becomes particularly critical because if Redis experiences downtime, it impacts the entire application. To address this, we can implement strategies such as deploying Redis in a high-availability configuration using Redis Cluster. Additionally, we can deploy these service on Kubernetes and enable automatic scaling and resilience by managing multiple instances of each service. Implementing load balancing across these instances can further enhance availability and fault tolerance and incorporating monitoring and alerting mechanisms helps in promptly identifying and mitigating issues before they impact the entire system. These approaches collectively aim to enhance the reliability and availability of our application architecture.

### Easy-to-Query APIs

To expose the stored data to customers in an easy-to-query API, I would consider implementing a GraphQL API on top of the existing API service. GraphQL provides a flexible and efficient way to query and retrieve data, allowing clients to request only the data they need.

### API Security

To secure the API, I would implement the following measures:

- **Authentication and Authorization:** Implement token-based authentication (e.g., JWT) and role-based access control (RBAC) to restrict access to the API.
  
- **Rate Limiting:** Implement rate limiting to prevent abuse and protect the service from being overwhelmed by too many requests.
    
- **HTTPS:** Use HTTPS to encrypt communication between clients and our API service.


## References

- [Ethereum JSON-RPC Specification](https://ethereum.github.io/execution-apis/api-documentation/)
- [Redis Documentation](https://redis.io/docs/latest/)

