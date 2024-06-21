package config

import (
	eth_err "ethereum-data-service/pkg/err"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      int
	RPCURL    string
	WSSURL    string
	RedisAddr string
}

func LoadConfig() (*Config, error) {

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// No .env file found, relying on system environment variables
		return nil, eth_err.ErrEnvVarMissing
	}

	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Port:      port,
		RPCURL:    getEnv("ETHEREUM_RPC_URL", ""),
		WSSURL:    getEnv("ETHEREUM_WSS_URL", ""),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
