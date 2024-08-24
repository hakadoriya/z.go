package envz

import (
	"fmt"
	"os"
	"strconv"
)

func Float32(key string) (float32, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseFloat(env, 32)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseFloat: %w", err)
	}

	return float32(value), nil
}

func Float32OrDefault(key string, defaultValue float32) float32 {
	value, err := Float32(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustFloat32(key string) float32 {
	env, err := Float32(key)
	if err != nil {
		panic(err)
	}

	return env
}
