# VC-Ethereum Data Service Architecture

```mermaid
graph TD

    classDef ethereumNodeStyle fill:#7FB3D5,stroke:#333,stroke-width:1px,color:white;
    classDef blockNotificationStyle fill:#96C7C0,stroke:#333,stroke-width:1px,color:white;
    classDef bootstrapperStyle fill:#FAC15E,stroke:#333,stroke-width:1px,color:white;
    classDef redisChannelStyle fill:#F08A5D,stroke:#333,stroke-width:1px,color:white;
    classDef blockSubscriberStyle fill:#7AA6A6,stroke:#333,stroke-width:1px,color:white;
    classDef dataFormatterStyle fill:#9b59b6,stroke:#333,stroke-width:1px,color:white;
    classDef redisDBStyle fill:#F6CF71,stroke:#333,stroke-width:1px,color:black;
    classDef apiServiceStyle fill:#B6D7A8,stroke:#333,stroke-width:1px,color:black;
    classDef clientsStyle fill:#B2B2B2,stroke:#333,stroke-width:1px,color:black;
    
    subgraph VC-ETH-Data-Service
        direction TB
        
        ethereumNode["Ethereum Node<br/>(WebSocket & HTTPS RPC)"]:::ethereumNodeStyle
        blockNotification["Block<br/>Notification"]:::blockNotificationStyle
        bootstrapper["Bootstrapper"]:::bootstrapperStyle
        redisChannel(["Redis Channel"]):::redisChannelStyle
        blockSubscriber["Block<br/>Subscriber"]:::blockSubscriberStyle
        dataFormatter["Data<br/>Formatter"]:::dataFormatterStyle
        redisDB[("Redis DB")]:::redisDBStyle
        apiService["API Service"]:::apiServiceStyle
        clients["Clients"]:::clientsStyle
        
        ethereumNode <--> |WebSocket| blockNotification
        bootstrapper -.-> |HTTPS<br/>Request| ethereumNode
        ethereumNode -.-> |HTTPS<br/>Response| bootstrapper
        blockNotification --> |New Block<br/>Info| redisChannel
        bootstrapper --> |Latest 50<br/>Block Info| dataFormatter
        redisChannel --> blockSubscriber
        blockSubscriber --> |Processed<br/>Block Data| dataFormatter
        dataFormatter --> |Formatted<br/>Block Data| redisDB
        redisDB --> |Data<br/>Retrieval| apiService
        clients --> |HTTP<br/>Request| apiService
        apiService --> |HTTP<br/>Response| clients
    end

```

## Overall Flow of VC-Ethereum Data Service Architecture

The VC-Ethereum Data Service architecture facilitates the retrieval, processing, and distribution of Ethereum blockchain data through various components as follows

**Ethereum Node**
- **Role:** Acts as the primary interface to the Ethereum blockchain, providing both real-time updates via WebSocket and historical data (most recent 50 blocks) retrieval via HTTPS RPC.

**Bootstrapper**
- **Role:** Fetches historical block data from the Ethereum Node using HTTPS RPC.
- **Flow:** Sends an HTTPS request to the Ethereum Node to retrieve the latest 50 blocks for initial data synchronization.

**Block Notification**
- **Role:** Listens for new blocks in real-time using WebSocket subscription from the Ethereum Node.
- **Flow:** Establishes a bi-directional WebSocket connection with the Ethereum Node to receive immediate updates on new blocks.

**Redis Channel**
- **Role:** Acts as an event-driven message broker for asynchronous communication.
- **Flow:** Receives new block information from the Block Notification service and forwards it to downstream components.

**Block Subscriber**
- **Role:** Subscribes to the Redis Channel to process incoming block data.
- **Flow:** Listens to the Redis Channel for new block updates, processes the data, and prepares it for further handling.

**Data Formatter**
- **Role:** A module that formats raw block data into a structured format according to predefined data models.
- **Integration:** Importable into other services such as Bootstrapper and Block Subscriber to ensure consistency in data formatting before storage.

**Redis DB**
- **Role:** Central storage for processed Ethereum blockchain data.
- **Flow:** Stores formatted block data received from the Data Formatter and Block Subscriber for efficient data retrieval.

**API Service**
- **Role:** Provides an interface for external clients to query Ethereum blockchain data.
- **Flow:** Receives HTTP requests from clients, retrieves requested data from Redis DB, and sends back HTTP responses containing the queried blockchain data.

**Clients**
- **Role:** External consumers of Ethereum blockchain data.
- **Flow:** Make HTTP requests to the API Service to fetch specific blockchain data of interest.


The `Data Formatter` component in this architecture is not a standalone service but rather a modular component that integrates into other services like the Bootstrapper and Block Subscriber. Its role is to ensure that all incoming block data, whether historical (from Bootstrapper) or real-time (from Block Subscriber), adheres to a consistent data format as specified in `model.Data` struct before being stored in Redis DB. This modular approach promotes code reusability and maintains data integrity across different parts of the Ethereum Data Service architecture. This module allows for efficient data processing, storage, and retrieval, ensuring that formatted blockchain data is readily available for consumption by API clients while maintaining consistency and reliability throughout the system.

Overall, this design strives to efficiently combine real-time and historical Ethereum data processing, leveraging Redis for both message brokering and data storage. The API Service ensures that clients have quick and reliable access to the latest blockchain data. 
