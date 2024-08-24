package envz

import (
	"fmt"
	"time"
)

func Second(key string) (time.Duration, error) {
	env, err := Int64(key)
	if err != nil {
		return 0, fmt.Errorf("Int64: %w", err)
	}

	return time.Duration(env) * time.Second, nil
}

func SecondOrDefault(key string, defaultValue time.Duration) time.Duration {
	value, err := Second(key)
	if err != nil {
		return defaultValue
	}

	return value
}

func MustSecond(key string) time.Duration {
	env, err := Second(key)
	if err != nil {
		panic(err)
	}

	return env
}
