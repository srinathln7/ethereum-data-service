# Ethereum Data Service

## Proposed Architecture

TBD 

## Proposed Project Structure

Your proposed project structure looks well-organized and modular, adhering to best practices in Go project organization. It provides a clear separation of concerns, which will make your codebase maintainable and scalable. Here's a brief overview of each component and why it's appropriate:


## Test

```

# Get block info associated with a given block number 
curl -X GET http://localhost:8080/block?block_number=20162001 | jq

# Get Tx info associated with a given tc hash
curl -X GET http://localhost:8080/tx?tx_hash=0xec9057951284893e709fd5e4d57f76f1013145f12d6e366664c060b6f2baf559 | jq

# Get events associated with a particular address
curl -X GET "http://localhost:8080/events?address=0x38AfDe1E098B0cA181E144A639339198c8CF3Efa" | jq

```


# References

* [Ethereum JSON-RPC Specification](https://ethereum.github.io/execution-apis/api-documentation/)

* [Ethereum-Block-Explorer](https://blockexplorer.one/ethereum/mainnet)
