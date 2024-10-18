package envz

import (
	"fmt"
	"os"
	"strconv"
)

func Int(key string) (int, error) {
	env, found := os.LookupEnv(key)
	if !found {
		return 0, fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	value, err := strconv.Atoi(env)
	if err != nil {
		return 0, fmt.Errorf("strconv.Atoi: %w", err)
	}

	return value, nil
}

func IntOrDefault(key string, defaultValue int) int {
	value, err := Int(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustInt(key string) int {
	env, err := Int(key)
	if err != nil {
		panic(err)
	}

	return env
}
