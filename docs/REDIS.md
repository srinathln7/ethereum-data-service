# Why Redis?

Redis emerges as an ideal choice for our project due to its exceptional performance, real-time capabilities, and scalability, which are essential for effectively managing data related to the most recent 50 Ethereum blocks and supporting a basic API. The factors considered while choosing Redis are outlined below:

## Estimate Storage Requirements 

Based on recent Ethereum block sizes from [Etherscan](https://etherscan.io/chart/blocksize), with an average block size of 150-200KB over the last 6 months, let us perform a simple back of the envelope calculation to appropximate our storage requirements

- **Block Data**: Approximately 200 KB per block.
- **Transaction Data**:
  - Average of 200 transactions per block.
  - Each transaction key (Hash: 32 bytes) and actual tx data (750 bytes) contribute to about 156 KB per block.
- **Event Data**:
  - Assuming an average of 2 events per transaction, contributing approximately 312 KB per block.

**Total Estimated Storage per Block**: 200 KB (block) + 156 KB (transactions) + 312 KB (events) ≈ 668 KB 
**Total for 50 Blocks**: 668 KB * 50 ≈ 33.4 MB

With additional overhead for indexing and metadata, our total storage requirement approximates to about 40MB, comfortably manageable within Redis's in-memory model.

### Performance and Reliability

- **In-Memory Operations**: Redis excels in rapid read and write operations, critical for handling real-time updates and queries on Ethereum block data.
- **Low Latency**: By leveraging Redis's in-memory storage, we ensure minimal latency in data access and updates, vital for maintaining responsive API performance.

### Real-Time Capabilities

- **Pub/Sub Model**: Redis's publish/subscribe feature facilitates seamless integration with Ethereum's WebSocket feeds, enabling real-time data updates and notifications.

## Flexibility and Scalability

- **TTL Support**: Redis's Time-To-Live (TTL) feature automates the management of data expiration, ensuring efficient use of memory resources as new blocks continuously arrive.
- **Clustering and Sharding**: Redis offers robust clustering and sharding capabilities, providing scalability to handle future increases in data volume without compromising performance.

### Concurrent Requests

- **Concurrency Support**: Redis efficiently handles multiple concurrent client requests thereby enhancing application responsiveness.

### Performance Monitoring

- **Redis-Insight**: Redis offers an excellent GUI through `redis-insight`, which can be a useful tool to monitor memory consumption in real-time. This provides a valuable alternative to `redis-cli`.

### Developer-Friendly Environment

- **Simplicity and Go Support**: Redis's straightforward data structures and extensive library support in Go streamline development efforts and maintenance tasks, enhancing developer productivity and code maintainability.


