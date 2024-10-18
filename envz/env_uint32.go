package envz

import (
	"fmt"
	"os"
	"strconv"
)

func Uint32(key string) (uint32, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseUint(env, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return uint32(value), nil
}

func Uint32OrDefault(key string, defaultValue uint32) uint32 {
	value, err := Uint32(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustUint32(key string) uint32 {
	env, err := Uint32(key)
	if err != nil {
		panic(err)
	}

	return env
}
