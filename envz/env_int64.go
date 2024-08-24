package envz

import (
	"fmt"
	"os"
	"strconv"
)

func Int64(key string) (int64, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseInt(env, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return value, nil
}

func Int64OrDefault(key string, defaultValue int64) int64 {
	value, err := Int64(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustInt64(key string) int64 {
	env, err := Int64(key)
	if err != nil {
		panic(err)
	}

	return env
}
