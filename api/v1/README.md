# API-server

## Overview
The `server.go` file initializes and manages the API server for handling HTTP requests, using the Gin web framework. It also ensures graceful shutdown of the server when required.

### Details

1. **RunAPIServer Function**
   - **Initialization**: 
     - Sets Gin framework mode to release to optimize performance.
     - Creates a new Gin router (`gin.New()`) without default middleware.
     - Adds essential middleware (`gin.Logger()` for request logging and `gin.Recovery()` for recovering from panics).
   
   - **Route Setup**:
     - Registers endpoint handlers using `setupHandlers(router, rdb)`, where `router` is the Gin router instance and `rdb` is the Redis client.
   
   - **Server Configuration**:
     - Configures an HTTP server (`http.Server`) to listen on the specified port (`cfg.API_PORT`) with the configured Gin router as its handler.
   
   - **Server Startup**:
     - Starts the HTTP server in a separate goroutine (`go func()`), allowing it to handle incoming requests concurrently.
   
   - **Graceful Shutdown**:
     - Waits for a shutdown signal (`<-shutdown`) to initiate server shutdown.
     - Sets a deadline using `context.WithTimeout` to wait for ongoing requests to finish within the specified default timeout (`cfg.DEFAULT_TIMEOUT`).
     - Calls `srv.Shutdown(ctx)` to gracefully stop the server, allowing ongoing requests to complete before shutting down.
   
   - **Logging**:
     - Logs server startup (`log.Printf`) and shutdown (`log.Println`) events for monitoring and debugging purposes.

### `handlers.go`

### Overview
The `handlers.go` file defines endpoint handlers for various API routes, managing HTTP request handling and response generation.

#### Details

1. **Endpoint Handlers**:
   - **Route Definition**: 
     - Each handler corresponds to a specific API endpoint (`/`, `/v1/blocks`, `/v1/events`, `/v1/block`, `/v1/tx`, `/favicon.ico`).
   
   - **Functionality**:
     - Uses Gin framework's `router.GET()` to define HTTP GET endpoints and associate them with handler functions.
     - Each handler function takes a Gin `Context` (`gin.Context`) as a parameter, which provides access to HTTP request parameters, headers, and response writer.
   
   - **Error Handling**:
     - Checks for required query parameters (`address`, `block_number`, `tx_hash`) in request queries and responds with appropriate HTTP status codes and error messages if parameters are missing.
     - Logs internal server errors (`http.StatusInternalServerError`) along with detailed error messages when fetching data from Redis fails (`storage` package functions like `GetEventsByAddress`, `GetAllBlockNumbers`, `GetBlockByNumber`, `GetTransactionByHash`).

2. **Utility Handler**:
   - **`handleFavicon` Function**:
     - Handles `favicon.ico` requests specifically by aborting further processing (`c.Abort()`), preventing logging of unnecessary requests.

#### Functionality Summary
- **Initialization**: Sets up Gin framework for handling HTTP requests with essential middleware.
- **Route Setup**: Registers endpoint handlers for API routes.
- **Server Startup**: Starts the HTTP server to listen on the configured port.
- **Graceful Shutdown**: Handles server shutdown gracefully, ensuring ongoing requests complete before shutting down.
- **Endpoint Handlers**: Implements specific logic for each API endpoint, validating inputs, fetching data from Redis, and returning appropriate HTTP responses.
- **Logging**: Provides logging for server startup, shutdown, and internal errors to aid in monitoring and troubleshooting.