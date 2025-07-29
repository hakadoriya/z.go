package cachez_test

import (
	"testing"
	"time"

	"github.com/hakadoriya/z.go/cachez"
)

// TestCase represents a test case with ID and Name
type TestCase struct {
	ID   int
	Name string
}

func TestCache_Generic_IntKey(t *testing.T) {
	t.Parallel()

	cache := cachez.New[int, string](
		cachez.WithMaxCapacity(3),
	)

	// Test with integer keys
	cache.Set(1, "one")
	cache.Set(2, "two")
	cache.Set(3, "three")

	val, ok := cache.Get(1)
	if !ok || val != "one" {
		t.Errorf("❌: expected 'one', got %v", val)
	}

	// Test eviction
	cache.Set(4, "four")

	if _, ok := cache.Get(2); ok {
		t.Error("key 2 should have been evicted")
	}
}

func TestCache_Generic_StructKey(t *testing.T) {
	t.Parallel()

	type CacheKey struct {
		UserID     int
		ResourceID string
	}

	cache := cachez.New[CacheKey, TestCase](
		cachez.WithMaxCapacity(10),
	)

	key1 := CacheKey{UserID: 1, ResourceID: "res1"}
	val1 := TestCase{ID: 100, Name: "Test 1"}

	key2 := CacheKey{UserID: 2, ResourceID: "res2"}
	val2 := TestCase{ID: 200, Name: "Test 2"}

	cache.Set(key1, val1)
	cache.Set(key2, val2)

	// Get with struct key
	result, ok := cache.Get(key1)
	if !ok {
		t.Error("key1 should exist")
	}
	if result.ID != 100 || result.Name != "Test 1" {
		t.Errorf("❌: unexpected value: %+v", result)
	}

	// Test deletion
	cache.Delete(key1)
	if _, ok := cache.Get(key1); ok {
		t.Error("key1 should have been deleted")
	}
}

func TestCache_Generic_WithTTL(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, int](
		cachez.WithDefaultTTL(50 * time.Millisecond),
	)

	cache.Set("count", 42)

	// Immediate retrieval
	val, ok := cache.Get("count")
	if !ok || val != 42 {
		t.Errorf("❌: expected 42, got %v", val)
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	if _, ok := cache.Get("count"); ok {
		t.Error("value should have expired")
	}
}

func TestCache_Generic_TypeSafety(t *testing.T) {
	t.Parallel()

	// This test demonstrates compile-time type safety
	intCache := cachez.New[string, int]()
	intCache.Set("number", 123)

	// This would cause a compile error if uncommented:
	// intCache.Set("number", "not a number")

	val, ok := intCache.Get("number")
	if !ok || val != 123 {
		t.Errorf("❌: expected 123, got %v", val)
	}

	// The returned value is already the correct type (int)
	// No type assertion needed
	doubled := val * 2
	if doubled != 246 {
		t.Errorf("❌: expected 246, got %d", doubled)
	}
}
