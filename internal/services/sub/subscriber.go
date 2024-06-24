package sub

import (
	"context"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/storage"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
)

func RunBlockSubsriber(rdb *redis.Client, cfg *config.Config) {
	// Create a common context instance
	ctx, cancel := context.WithCancel(context.Background())

	// Subscribe to the Redis channel
	subscriber := rdb.Subscribe(ctx, cfg.REDIS_PUBSUB_CH)

	// Channel to receive messages
	ch := subscriber.Channel()

	go handleGracefulShutdown(subscriber, cancel)

	for msg := range ch {
		err := storage.AddBlockDataToDB(ctx, rdb, []byte(msg.Payload), cfg.REDIS_KEY_EXPIRY_TIME)
		if err != nil {
			log.Printf("error adding block data to storage: %v\n", err)
		}
	}

}

func handleGracefulShutdown(sub *redis.PubSub, cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a termination signal
	sig := <-sigCh
	log.Printf("Received signal %v. Initating graceful shut down...", sig)

	// Close all client connections and trigger cancellation of context
	sub.Close()
	cancel()
}
