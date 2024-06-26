# Ethereum Data Service

## Pre-requesties
* Go
* Redis 


## Install

```

# Docker 
docker pull redis
docker run -d --name vc-redis -p 6379:6379 redis

# In case Redis is running locally
sudo systemctl stop redis

docker start vc-redis

```

## Proposed Architecture

TBD 

##  Project Structure

TBD

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


