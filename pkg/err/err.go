package err

import (
	"errors"
	"fmt"
)

var (
	ErrEnvFileMissing  = errors.New("environment config variable missing")
	ErrInvalidProtocol = errors.New("invalid protocol specified")
)

func ConfigKeyMissingError(key string) error {
	return fmt.Errorf("specified key:%s missing from config", key)
}
