package config

import (
	util "ethereum-data-service/pkg/util"
	"strconv"
	"time"
)

type Config struct {
	// DEFAULT_TIMEOUT is the default timeout (seconds) duration for network requests.
	DEFAULT_TIMEOUT time.Duration

	// API_PORT is the port on which the API server will listen.
	API_PORT string

	// ETH_HTTPS_URL is the HTTPS URL for accessing the Ethereum network.
	ETH_HTTPS_URL string
	// ETH_WSS_URL is the WebSocket URL for accessing the Ethereum network.
	ETH_WSS_URL string

	// REDIS_DB is the Redis database number to use.
	REDIS_DB int
	// REDIS_ADDR is the address of the Redis server.
	REDIS_ADDR string
	// REDIS_PUBSUB_CH is the Redis Pub/Sub channel name for messaging.
	REDIS_PUBSUB_CH string
	// REDIS_KEY_EXPIRY_TIME is the default expiration time (seconds) for keys stored in Redis.
	// It is calculated based on Ethereum avg block time (~13s). Currently set to 50*13=650s since
	// we need to store info only for about 50 blocks
	REDIS_KEY_EXPIRY_TIME time.Duration

	// NUM_BLOCKS_TO_SYNC is the number of recent blocks to sync during initialization.
	NUM_BLOCKS_TO_SYNC int
	// BOOTSTRAP_TIME_OUT is the time (minute) after the bootstrap service exits itself gracefully.
	// On an avg, to fetch the most recent 50 blocks and load it to Redis, it takes approximately 6 mins.
	// We set it to 10 mins for safety marigin.
	BOOTSTRAP_TIMEOUT time.Duration
}

func LoadConfig() (*Config, error) {
	requiredKeys := []string{
		"DEFAULT_TIMEOUT",
		"API_PORT",
		"ETH_HTTPS_URL", "ETH_WSS_URL",
		"REDIS_ADDR", "REDIS_DB", "REDIS_PUBSUB_CH", "REDIS_KEY_EXPIRY_TIME",
		"NUM_BLOCKS_TO_SYNC", "BOOTSTRAP_TIMEOUT",
	}

	envMap, err := util.GetEnvMap(requiredKeys)
	if err != nil {
		return nil, err
	}

	defaultTimeout, err := strconv.Atoi(envMap["DEFAULT_TIMEOUT"])
	if err != nil {
		return nil, err
	}

	rdb, err := strconv.Atoi(envMap["REDIS_DB"])
	if err != nil {
		return nil, err
	}

	expiryTime, err := strconv.Atoi(envMap["REDIS_KEY_EXPIRY_TIME"])
	if err != nil {
		return nil, err
	}

	syncNum, err := strconv.Atoi(envMap["NUM_BLOCKS_TO_SYNC"])
	if err != nil {
		return nil, err
	}

	bootstrapTimeout, err := strconv.Atoi(envMap["BOOTSTRAP_TIMEOUT"])
	if err != nil {
		return nil, err
	}

	return &Config{
		DEFAULT_TIMEOUT: time.Duration(defaultTimeout) * time.Second,

		API_PORT: envMap["API_PORT"],

		ETH_HTTPS_URL: envMap["ETH_HTTPS_URL"],
		ETH_WSS_URL:   envMap["ETH_WSS_URL"],

		REDIS_DB:              rdb,
		REDIS_KEY_EXPIRY_TIME: time.Duration(expiryTime) * time.Second,
		REDIS_ADDR:            envMap["REDIS_ADDR"],
		REDIS_PUBSUB_CH:       envMap["REDIS_PUBSUB_CH"],

		NUM_BLOCKS_TO_SYNC: syncNum,
		BOOTSTRAP_TIMEOUT:  time.Duration(bootstrapTimeout) * time.Minute,
	}, nil
}
