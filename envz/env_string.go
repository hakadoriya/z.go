package envz

import (
	"fmt"
	"os"
)

func String(key string) (string, error) {
	value, found := os.LookupEnv(key)
	if !found {
		return "", fmt.Errorf("%s: %w", key, ErrEnvironmentVariableIsEmpty)
	}

	return value, nil
}

func StringOrDefault(key string, defaultValue string) string {
	value, err := String(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustString(key string) string {
	value, err := String(key)
	if err != nil {
		panic(err)
	}

	return value
}
