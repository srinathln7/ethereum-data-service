package config

import (
	util "ethereum-data-service/pkg/util"
	"strconv"
)

type Config struct {
	Port          int
	HTTPSURL      string
	WSSURL        string
	RedisAddr     string
	RedisPubSubCh string
}

func LoadConfig() (*Config, error) {
	requiredKeys := []string{"PORT", "ETHEREUM_HTTPS_URL", "ETHEREUM_WSS_URL", "REDIS_ADDR", "REDIS_PUBSUB_CH"}
	envMap, err := util.GetEnvMap(requiredKeys)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(envMap["PORT"])
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:          port,
		HTTPSURL:      envMap["ETHEREUM_HTTPS_URL"],
		WSSURL:        envMap["ETHEREUM_WSS_URL"],
		RedisAddr:     envMap["REDIS_ADDR"],
		RedisPubSubCh: envMap["REDIS_PUBSUB_CH"],
	}, nil
}
