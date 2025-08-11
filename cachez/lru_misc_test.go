package cachez

import "testing"

func Test_lruCacheG_removeElement(t *testing.T) {
	cache := New[string, string](
		WithMaxCapacity(5),
		WithEvictionPolicy(LRU),
	)

	before := _testAssert_lruCacheG_removeElement
	_testAssert_lruCacheG_removeElement = true
	t.Cleanup(func() { _testAssert_lruCacheG_removeElement = before })

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")
	cache.Set("key5", "value5")
}
