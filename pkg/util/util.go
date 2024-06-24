package util

import (
	eth_err "ethereum-data-service/pkg/err"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvMap(keys []string) (map[string]string, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// No .env file found, relying on system environment variables
		return nil, eth_err.ErrEnvFileMissing
	}

	envMap := make(map[string]string)
	for _, key := range keys {
		if value, exists := os.LookupEnv(key); exists {
			envMap[key] = value
		} else {
			return nil, eth_err.ConfigKeyMissingError(key)
		}
	}

	return envMap, nil
}
