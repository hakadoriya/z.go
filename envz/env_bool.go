package envz

import (
	"fmt"
	"os"
	"strconv"
)

func Bool(key string) (bool, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return false, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.ParseBool(env)
	if err != nil {
		return false, fmt.Errorf("strconv.ParseBool: %w", err)
	}

	return value, nil
}

func BoolOrDefault(key string, defaultValue bool) bool {
	value, err := Bool(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustBool(key string) bool {
	env, err := Bool(key)
	if err != nil {
		panic(err)
	}

	return env
}
