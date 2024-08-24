package envz

import (
	"fmt"
	"os"
	"strconv"
)

func Float64(key string) (float64, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseFloat(env, 64)
	if err != nil {
		return 0, fmt.Errorf("strconv.ParseFloat: %w", err)
	}

	return value, nil
}

func Float64OrDefault(key string, defaultValue float64) float64 {
	value, err := Float64(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustFloat64(key string) float64 {
	env, err := Float64(key)
	if err != nil {
		panic(err)
	}

	return env
}
