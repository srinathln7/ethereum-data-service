# Ethereum Data Service

## Pre-requesties
* [Go]((https://go.dev/doc/install)
* [Redis](https://redis.io/) 
* Docker (v20.10.21 or higher)
* Docker Compose (v20.20.2 or higher)

## Proposed Architecture

TBD 

##  Project Structure

Refer [here](https://github.com/srinathln7/ethereum-data-service/blob/main/CHALLENGE.md) for the challenge description and the project structure.


## Get Started

To build and start all the required services run

```
make buildup
```

To view logs real-time for every service 

```
docker logs -f vc-api-server

docker logs -f vc-bootstrapper

docker logs -f vc-blocknotifier

docker logs -f vc-blocksubscriber

```


## Test

```
# VC-00: - List all routes 
curl -X GET http://localhost:8080/ | jq 

# VC-01: List all block numbers currently available in local data store
curl -X GET http://localhost:8080/blocks | jq

# VC-02: Get block info associated with a given block number 
curl -X GET http://localhost:8080/block?block_number=20162001 | jq

# VC-03: Get Tx info associated with a given tc hash
curl -X GET http://localhost:8080/tx?tx_hash=0xec9057951284893e709fd5e4d57f76f1013145f12d6e366664c060b6f2baf559 | jq

# VC-04: Get all events associated with a particular address
curl -X GET "http://localhost:8080/events?address=0x38AfDe1E098B0cA181E144A639339198c8CF3Efa" | jq

```


## References

* [Ethereum JSON-RPC Specification](https://ethereum.github.io/execution-apis/api-documentation/)

* [Redis-Docs](https://redis.io/docs/latest/)
