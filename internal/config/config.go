package config

import (
	util "ethereum-data-service/pkg/util"
	"strconv"
	"time"
)

type Config struct {
	API_PORT        string
	API_STATIC_FILE string

	ETH_HTTPS_URL string
	ETH_WSS_URL   string

	REDIS_DB              int
	REDIS_ADDR            string
	REDIS_PUBSUB_CH       string
	REDIS_KEY_EXPIRY_TIME time.Duration

	NUM_BLOCKS_TO_SYNC int
}

func LoadConfig() (*Config, error) {
	requiredKeys := []string{
		"API_PORT", "API_STATIC_FILE",
		"ETHEREUM_HTTPS_URL", "ETHEREUM_WSS_URL",
		"REDIS_ADDR", "REDIS_DB", "REDIS_PUBSUB_CH", "REDIS_KEY_EXPIRY_TIME",
		"NUM_BLOCKS_TO_SYNC",
	}

	envMap, err := util.GetEnvMap(requiredKeys)
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

	return &Config{
		API_PORT:        envMap["API_PORT"],
		API_STATIC_FILE: envMap["API_STATIC_FILE"],

		ETH_HTTPS_URL: envMap["ETHEREUM_HTTPS_URL"],
		ETH_WSS_URL:   envMap["ETHEREUM_WSS_URL"],

		REDIS_DB:              rdb,
		REDIS_KEY_EXPIRY_TIME: time.Duration(expiryTime) * time.Second,
		REDIS_ADDR:            envMap["REDIS_ADDR"],
		REDIS_PUBSUB_CH:       envMap["REDIS_PUBSUB_CH"],

		NUM_BLOCKS_TO_SYNC: syncNum,
	}, nil
}
