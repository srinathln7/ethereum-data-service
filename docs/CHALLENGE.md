# Ethereum Data Service

## Challenge Overview

Please write a service which uses the Ethereum JSON-RPC API to store the following information in a local datastore for the most recent 50 blocks and provide a basic API:

1. Get and store the block and all transaction hashes in the block.
2. Get and store all events related to each transaction in each block.
3. Use the stored data to serve three endpoints:
   - Query all events related to a particular address.
   - Get block by number.
   - Get transaction by hash.

### Use the following Validation Cloud endpoints for Ethereum Mainnet RPC access:

- HTTP: `https://mainnet.ethereum.validationcloud.io/v1/JFt58zlN7gcQlLYnMZcOD75LpethJgD6Eq5nKOxC9F0`
- WebSocket: `wss://mainnet.ethereum.validationcloud.io/v1/wss/JFt58zlN7gcQlLYnMZcOD75LpethJgD6Eq5nKOxC9F0`

### Additional Notes:

- Code should be written in Go.
- RPC Methods which might be useful: `eth_blockNumber`, `eth_getBlockByNumber`, `eth_getTransactionByHash`, `eth_getLogs`.
- Feel free to ask us questions as you go.
- Include a Dockerfile to compile and build your application.
- Containerize any dependencies (e.g. databases).
- Include a README file with instructions on how to build the container(s) and run the application.

### Some inspiration for add-ons and discussion topics in our debrief:

- Thought process and code design.
- Technology choices.
- How would you keep the data set up to date?
- How would you expose the stored data to customers in an easy-to-query API?
- How would you handle security of the API?
- How would you improve the performance of your approach?
- How would you adapt your design to store the same data for the entire history of Ethereum Mainnet?
- What would it take to deploy and monitor a service like this in production?


## Project structure

```
.
├── api
│   └── v1
│       ├── handlers.go
│       ├── middleware.go
│       └── server.go
├── CHALLENGE.md
├── cmd
│   └── cmd.go
├── docker-compose.yml
├── Dockerfile
├── docs
│   ├── bootstrap.html
│   ├── client.html
│   ├── config.html
│   ├── index.html
│   ├── model.html
│   ├── pub.html
│   ├── storage.html
│   └── sub.html
├── go.mod
├── go.sum
├── internal
│   ├── client
│   │   └── client.go
│   ├── config
│   │   └── config.go
│   ├── model
│   │   └── model.go
│   ├── services
│   │   ├── bootstrapper
│   │   │   └── bootstrap.go
│   │   ├── pub
│   │   │   └── publisher.go
│   │   └── sub
│   │       └── subscriber.go
│   └── storage
│       ├── index.go
│       ├── queries.go
│       └── storage.go
├── LICENSE
├── main.go
├── Makefile
├── pkg
│   ├── enum
│   │   └── enum.go
│   ├── err
│   │   └── err.go
│   ├── log
│   │   └── logger.go
│   └── util
│       └── util.go
└── README.md

18 directories, 34 files

```
