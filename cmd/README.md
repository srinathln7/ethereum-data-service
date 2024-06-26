# VC-ETH-DATA-Service commands

The `cmd` package is responsible for setting up and managing commands through Cobra, a popular CLI library for Go. It provides commands to start different services related to Ethereum data handling and an HTTP API server.

### `Init()` Function

- **Initialization**: Loads configuration (`config.LoadConfig()`) and initializes clients (`client.InitClient()`). Errors during initialization cause the CLI to exit with a fatal log message.

### Commands Defined in `RootCmd`

- **Root Command (`ethereum_api_service`)**: 
  - Prints welcome message and usage instructions for sub-commands (`bootstrap`, `api-server`, `pub`, `sub`).

### `bootstrapCmd`

- **`bootstrap` Command**: Starts the BlockBootstrapper service.
  - **Functionality**: Spawns a goroutine (`bootstrapper.RunBootstrapSvc()`) to run the bootstrapper service using configured client and settings.
  - **Shutdown**: Uses `handleShutdown()` to handle graceful shutdown of the bootstrapper service.

### `pubCmd`

- **`pub` Command**: Starts the BlockNotification service.
  - **Functionality**: Spawns a goroutine (`pub.RunBlockNotifierSvc()`) to run the notification service using configured client and settings.
  - **Shutdown**: Uses `handleShutdown()` to handle graceful shutdown of the notification service.

### `subCmd`

- **`sub` Command**: Starts the BlockSubscriber service.
  - **Functionality**: Spawns a goroutine (`sub.RunBlockSubscriberSvc()`) to run the subscriber service using the Redis client and configured settings.
  - **Shutdown**: Uses `handleShutdown()` to handle graceful shutdown of the subscriber service.

### `apiServerCmd`

- **`api-server` Command**: Starts the HTTP-API server.
  - **Functionality**: Spawns a goroutine (`v1.RunAPIServer()`) to run the HTTP API server using the Redis client and configured settings.
  - **Shutdown**: Uses `handleShutdown()` to handle graceful shutdown of the API server.

### `handleShutdown()` Function

- **Graceful Shutdown Handling**:
  - **Signal Handling**: Captures interrupt (`os.Interrupt`) and termination signals (`syscall.SIGTERM`) using `os/signal`.
  - **Shutdown Initiation**: Closes the `shutdown` channel to signal all services to start shutting down.
  - **Wait for Goroutines**: Waits for all spawned goroutines (services) to finish using `sync.WaitGroup`.
  - **Timeout Handling**: Sets a timeout (`cfg.DEFAULT_TIMEOUT`) for shutdown operations and logs if shutdown exceeds this timeout.

- The `cmd` package effectively manages starting and stopping of Ethereum-related services and an HTTP API server using Cobra for command-line interface management.
- Each command (`bootstrap`, `pub`, `sub`, `api-server`) spawns a goroutine for the respective service and handles graceful shutdown upon receiving termination signals.
- Error handling ensures that initialization failures are logged and cause immediate termination of the CLI.

This setup provides a robust mechanism to start and manage Ethereum data services and an API server through a CLI interface, ensuring reliability and graceful shutdown during operational tasks.