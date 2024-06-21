package err

import "errors"

var (
	ErrEnvVarMissing = errors.New("environment config variable missing")
)
