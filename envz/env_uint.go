package envz

import (
	"fmt"
	"os"
	"strconv"
)

func Uint(key string) (uint, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseUint(env, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseInt: %w", err)
	}

	return uint(value), nil
}

func UintOrDefault(key string, defaultValue uint) uint {
	value, err := Uint(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustUint(key string) uint {
	env, err := Uint(key)
	if err != nil {
		panic(err)
	}

	return env
}
