package cachez_test

import (
	"testing"
	"time"

	"github.com/hakadoriya/z.go/cachez"
)

func TestCache_FOP_DefaultValues(t *testing.T) {
	t.Parallel()

	// Create cache with no options
	cache := cachez.New[string, string]()

	// Default capacity should be 100
	// Fill up to 100 entries
	for i := 0; i < 100; i++ {
		cache.Set(string(rune('a'+i%26))+string(rune('0'+i/26)), "value")
	}

	if size := cache.Size(); size != 100 {
		t.Errorf("❌: expected size 100 with default capacity, got %d", size)
	}

	// Adding one more should trigger eviction
	cache.Set("overflow", "value")
	if size := cache.Size(); size != 100 {
		t.Errorf("❌: expected size to remain 100 after eviction, got %d", size)
	}

	// Default eviction policy should be LRU
	// This is implicitly tested by the LRU test cases

	// Default TTL should be 0 (no expiration)
	cache.Set("no-ttl", "value")
	time.Sleep(10 * time.Millisecond)
	if _, ok := cache.Get("no-ttl"); !ok {
		t.Error("value should not expire with default TTL of 0")
	}
}

func TestCache_FOP_WithOptions(t *testing.T) {
	t.Parallel()

	// Test WithMaxCapacity
	cache1 := cachez.New[string, int](
		cachez.WithMaxCapacity(5),
	)

	for i := 0; i < 10; i++ {
		cache1.Set(string(rune('a'+i)), i)
	}

	if size := cache1.Size(); size != 5 {
		t.Errorf("❌: expected size 5 with custom capacity, got %d", size)
	}

	// Test WithEvictionPolicy
	cache2 := cachez.New[string, string](
		cachez.WithEvictionPolicy(cachez.LFU),
		cachez.WithMaxCapacity(3),
	)

	cache2.Set("a", "1")
	cache2.Set("b", "2")
	cache2.Set("c", "3")

	// Access 'a' multiple times
	cache2.Get("a")
	cache2.Get("a")
	cache2.Get("b")

	// Add new item, 'c' should be evicted (least frequently used)
	cache2.Set("d", "4")

	if _, ok := cache2.Get("c"); ok {
		t.Error("'c' should have been evicted by LFU policy")
	}

	// Test WithDefaultTTL
	cache3 := cachez.New[string, string](
		cachez.WithDefaultTTL(50 * time.Millisecond),
	)

	cache3.Set("ttl-test", "value")

	// Should exist immediately
	if _, ok := cache3.Get("ttl-test"); !ok {
		t.Error("value should exist immediately after set")
	}

	// Should expire after TTL
	time.Sleep(100 * time.Millisecond)
	if _, ok := cache3.Get("ttl-test"); ok {
		t.Error("value should have expired after default TTL")
	}
}

func TestCache_FOP_MultipleOptions(t *testing.T) {
	t.Parallel()

	// Test combining multiple options
	cache := cachez.New[int, string](
		cachez.WithMaxCapacity(10),
		cachez.WithEvictionPolicy(cachez.LFU),
		cachez.WithDefaultTTL(100*time.Millisecond),
	)

	// Test capacity
	for i := 0; i < 15; i++ {
		cache.Set(i, "value")
	}

	if size := cache.Size(); size != 10 {
		t.Errorf("❌: expected size 10, got %d", size)
	}

	// Test TTL
	cache.Set(100, "ttl-value")
	time.Sleep(150 * time.Millisecond)
	if _, ok := cache.Get(100); ok {
		t.Error("value should have expired")
	}
}

func TestCache_FOP_InvalidCapacity(t *testing.T) {
	t.Parallel()

	// Test that invalid capacity is ignored
	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(-10),
	)

	// Should use default capacity
	for i := 0; i < 100; i++ {
		cache.Set(string(rune('a'+i%26))+string(rune('0'+i/26)), "value")
	}

	if size := cache.Size(); size != 100 {
		t.Errorf("❌: expected default capacity 100, got %d", size)
	}
}
