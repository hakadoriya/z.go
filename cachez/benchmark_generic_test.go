package cachez_test

import (
	"testing"

	"github.com/hakadoriya/z.go/cachez"
)

func BenchmarkCache_IntKey_Set(b *testing.B) {
	cache := cachez.New[int, string](
		cachez.WithMaxCapacity(10000),
	)

	b.ResetTimer()

	for i := range b.N {
		cache.Set(i%10000, "value")
	}
}

func BenchmarkCache_IntKey_Get(b *testing.B) {
	cache := cachez.New[int, string](
		cachez.WithMaxCapacity(10000),
	)

	// Pre-populate cache
	for i := range 10000 {
		cache.Set(i, "value")
	}

	b.ResetTimer()

	for i := range b.N {
		cache.Get(i % 10000)
	}
}

type benchKey struct {
	ID   int
	Type string
}

func BenchmarkCache_StructKey_Set(b *testing.B) {
	cache := cachez.New[benchKey, string](
		cachez.WithMaxCapacity(10000),
	)

	keys := make([]benchKey, 10000)
	for i := range keys {
		keys[i] = benchKey{ID: i, Type: "bench"}
	}

	b.ResetTimer()

	for i := range b.N {
		key := keys[i%10000]
		cache.Set(key, "value")
	}
}

func BenchmarkCache_StructKey_Get(b *testing.B) {
	cache := cachez.New[benchKey, string](
		cachez.WithMaxCapacity(10000),
	)

	keys := make([]benchKey, 10000)
	for i := range keys {
		keys[i] = benchKey{ID: i, Type: "bench"}
		cache.Set(keys[i], "value")
	}

	b.ResetTimer()

	for i := range b.N {
		key := keys[i%10000]
		cache.Get(key)
	}
}