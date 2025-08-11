package cachez

import (
	"testing"
	"time"
)

func Test_lruCacheG_removeElement(t *testing.T) {
	cache := New[string, string](
		WithMaxCapacity(5),
		WithEvictionPolicy(LRU),
	)

	before := _testAssert_lruCacheG_removeElement
	_testAssert_lruCacheG_removeElement = true
	t.Cleanup(func() { _testAssert_lruCacheG_removeElement = before })

	// Set entries to create list elements
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")
	cache.Set("key5", "value5")
	
	// This should trigger eviction and call removeElement
	cache.Set("key6", "value6")
}

func Test_lruCacheG_Get_TypeAssertion(t *testing.T) {
	cache := New[string, string](
		WithMaxCapacity(5),
		WithEvictionPolicy(LRU),
	)

	before := _testAssert_lruCacheG_Get
	_testAssert_lruCacheG_Get = true
	t.Cleanup(func() { _testAssert_lruCacheG_Get = before })

	// Set an entry
	cache.Set("key1", "value1")
	
	// Try to get it - should fail due to type assertion flag
	val, ok := cache.Get("key1")
	if ok {
		t.Errorf("expected Get to fail with type assertion, got value: %v", val)
	}
}

func Test_lruCacheG_SetWithTTL_TypeAssertion(t *testing.T) {
	cache := New[string, string](
		WithMaxCapacity(5),
		WithEvictionPolicy(LRU),
	)

	before := _testAssert_lruCacheG_SetWithTTL
	_testAssert_lruCacheG_SetWithTTL = true
	t.Cleanup(func() { _testAssert_lruCacheG_SetWithTTL = before })

	// Set an entry first
	cache.Set("key1", "value1")
	
	// Try to update it with TTL - should fail silently due to type assertion flag
	cache.SetWithTTL("key1", "updated", 1*time.Hour)
	
	// Reset the flag to check if value was not updated
	_testAssert_lruCacheG_SetWithTTL = false
	
	// Value should still be the original since update failed
	val, ok := cache.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("expected original value 'value1', got: %v", val)
	}
}
