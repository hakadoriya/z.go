package envz

import (
	"fmt"
	"os"
	"strconv"
)

func Int32(key string) (int32, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseInt(env, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return int32(value), nil
}

func Int32OrDefault(key string, defaultValue int32) int32 {
	value, err := Int32(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustInt32(key string) int32 {
	env, err := Int32(key)
	if err != nil {
		panic(err)
	}

	return env
}
