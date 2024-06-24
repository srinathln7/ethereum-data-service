package sub

import (
	"context"
	"ethereum-data-service/internal/config"
	"ethereum-data-service/internal/storage"
	"log"

	"github.com/redis/go-redis/v9"
)

// RunBlockSubscriber: Subscribes to a Redis channel and stores incoming block data to storage.
func RunBlockSubscriber(rdb *redis.Client, cfg *config.Config, shutdown chan struct{}) {
	// Create a common context instance
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals for graceful shutdown
	go handleGracefulShutdown(cancel, shutdown)

	log.Printf("Subscribed to Redis channel: %s\n", cfg.REDIS_PUBSUB_CH)

	// Subscribe to the Redis channel
	subscriber := rdb.Subscribe(ctx, cfg.REDIS_PUBSUB_CH)

	// Channel to receive messages
	ch := subscriber.Channel()

	for {
		select {
		case <-shutdown:
			log.Println("Shutting down BlockSubscriber service...")
			subscriber.Close()
			cancel()
			return
		case msg := <-ch:
			err := storage.AddBlockDataToDB(ctx, rdb, []byte(msg.Payload), cfg.REDIS_KEY_EXPIRY_TIME)
			if err != nil {
				log.Printf("error adding block data to storage: %v\n", err)
			}
		}
	}
}

func handleGracefulShutdown(cancel context.CancelFunc, shutdown chan struct{}) {
	<-shutdown
	cancel()
}
