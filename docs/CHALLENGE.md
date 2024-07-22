# Ethereum Data Service

## Objective

To write a service which uses the Ethereum JSON-RPC API to store the following information in a local datastore for the most recent 50 blocks and provide a basic API:

1. Get and store the block and all transaction hashes in the block.
2. Get and store all events related to each transaction in each block.
3. Use the stored data to serve three endpoints:
   - Query all events related to a particular address.
   - Get block by number.
   - Get transaction by hash.

### Use the following endpoints for Ethereum Mainnet RPC access:
- HTTP: `https://mainnet.ethereum.validationcloud.io/v1/JFt58zlN7gcQlLYnMZcOD75LpethJgD6Eq5nKOxC9F0`
- WebSocket: `wss://mainnet.ethereum.validationcloud.io/v1/wss/JFt58zlN7gcQlLYnMZcOD75LpethJgD6Eq5nKOxC9F0`

## Project structure

```
├── api
│   └── v1
│       ├── handlers.go
│       ├── README.md
│       └── server.go
├── cmd
│   ├── cmd.go
│   └── README.md
├── DESIGN.md
├── docker-compose.yml
├── Dockerfile
├── docs
│   ├── CHALLENGE.md
│   ├── html
│   │   ├── bootstrap.html
│   │   ├── client.html
│   │   ├── config.html
│   │   ├── index.html
│   │   ├── model.html
│   │   ├── pub.html
│   │   ├── storage.html
│   │   └── sub.html
│   └── REDIS.md
├── go.mod
├── go.sum
├── internal
│   ├── client
│   │   ├── client.go
│   │   └── README.md
│   ├── config
│   │   ├── config.go
│   │   └── README.md
│   ├── model
│   │   ├── model.go
│   │   └── README.md
│   ├── services
│   │   ├── bootstrapper
│   │   │   ├── bootstrap.go
│   │   │   └── README.md
│   │   ├── pub
│   │   │   ├── publisher.go
│   │   │   └── README.md
│   │   └── sub
│   │       ├── README.md
│   │       └── subscriber.go
│   └── storage
│       ├── index.go
│       ├── queries.go
│       ├── README.md
│       └── storage.go
├── LICENSE
├── main.go
├── Makefile
├── pkg
│   ├── enum
│   │   └── enum.go
│   ├── err
│   │   └── err.go
│   └── util
│       └── util.go
├── README.md
└── sample.env

18 directories, 44 files

```
