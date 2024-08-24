package envz

import (
	"fmt"
	"os"
	"strconv"
)

func Uint64(key string) (uint64, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseUint(env, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return value, nil
}

func Uint64OrDefault(key string, defaultValue uint64) uint64 {
	value, err := Uint64(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustUint64(key string) uint64 {
	env, err := Uint64(key)
	if err != nil {
		panic(err)
	}

	return env
}
