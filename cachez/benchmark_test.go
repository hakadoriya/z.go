package cachez_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hakadoriya/z.go/cachez"
)

func BenchmarkCache_LRU_Set(b *testing.B) {
	cache := cachez.New[string, int](
		cachez.WithMaxCapacity(10000),
		cachez.WithEvictionPolicy(cachez.LRU),
	)

	b.ResetTimer()

	for i := range b.N {
		key := fmt.Sprintf("key%d", i%10000)
		cache.Set(key, i)
	}
}

func BenchmarkCache_LRU_Get(b *testing.B) {
	cache := cachez.New[string, int](
		cachez.WithMaxCapacity(10000),
		cachez.WithEvictionPolicy(cachez.LRU),
	)

	// 事前にデータを投入
	for i := range 10000 {
		key := fmt.Sprintf("key%d", i)
		cache.Set(key, i)
	}

	b.ResetTimer()

	for i := range b.N {
		key := fmt.Sprintf("key%d", i%10000)
		cache.Get(key)
	}
}

func BenchmarkCache_LFU_Set(b *testing.B) {
	cache := cachez.New[string, int](
		cachez.WithMaxCapacity(10000),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

	b.ResetTimer()

	for i := range b.N {
		key := fmt.Sprintf("key%d", i%10000)
		cache.Set(key, i)
	}
}

func BenchmarkCache_LFU_Get(b *testing.B) {
	cache := cachez.New[string, int](
		cachez.WithMaxCapacity(10000),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

	// 事前にデータを投入
	for i := range 10000 {
		key := fmt.Sprintf("key%d", i)
		cache.Set(key, i)
	}

	b.ResetTimer()

	for i := range b.N {
		key := fmt.Sprintf("key%d", i%10000)
		cache.Get(key)
	}
}

func BenchmarkCache_SetWithTTL(b *testing.B) {
	cache := cachez.New[string, int](
		cachez.WithMaxCapacity(10000),
	)

	b.ResetTimer()

	for i := range b.N {
		key := fmt.Sprintf("key%d", i%10000)
		cache.SetWithTTL(key, i, 5*time.Minute)
	}
}

func BenchmarkCache_ConcurrentAccess(b *testing.B) {
	cache := cachez.New[string, int](
		cachez.WithMaxCapacity(10000),
	)

	// 事前にデータを投入
	for i := range 10000 {
		key := fmt.Sprintf("key%d", i)
		cache.Set(key, i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%10000)
			if i%2 == 0 {
				cache.Get(key)
			} else {
				cache.Set(key, i)
			}

			i++
		}
	})
}
