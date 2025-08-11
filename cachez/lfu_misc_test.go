package cachez

import (
	"errors"
	"testing"
)

func Test_lfuCacheG_evictLeastFrequent(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic")
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("expected error, got %v", r)
		}
		if !errors.Is(err, ErrAssertionFailed) {
			t.Errorf("expected ErrAssertionFailed, got %v", err)
		}
	}()

	cache := New[string, string](
		WithMaxCapacity(5),
		WithEvictionPolicy(LFU),
	)

	before := _testAssert_lfuCacheG_evictLeastFrequent
	_testAssert_lfuCacheG_evictLeastFrequent = true
	t.Cleanup(func() { _testAssert_lfuCacheG_evictLeastFrequent = before })

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")
	cache.Set("key5", "value5")
	cache.Set("key6", "value6")
}

func Test_minHeapG_Push(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("expected panic")
		}
		err, ok := r.(error)
		if !ok {
			t.Errorf("expected error, got %v", r)
		}
		if !errors.Is(err, ErrAssertionFailed) {
			t.Errorf("expected ErrAssertionFailed, got %v", err)
		}
	}()

	cache := New[string, string](
		WithMaxCapacity(5),
		WithEvictionPolicy(LFU),
	)

	before := _testAssert_minHeapG_Push
	_testAssert_minHeapG_Push = true
	t.Cleanup(func() { _testAssert_minHeapG_Push = before })

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")
	cache.Set("key5", "value5")
	cache.Set("key6", "value6")
}
