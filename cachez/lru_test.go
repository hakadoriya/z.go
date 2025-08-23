package cachez_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/hakadoriya/z.go/cachez"
)

func TestLRU_EvictionOrder(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxLength(5),
		cachez.WithEvictionPolicy(cachez.LRU),
	)

	// Add 5 elements
	// ja: 5つの要素を追加
	for i := 1; i <= 5; i++ {
		cache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}

	// Access key1, key3, key5 (change usage order)
	// ja: key1, key3, key5 にアクセス（使用順を変更）
	cache.Get("key1")
	cache.Get("key3")
	cache.Get("key5")

	// Add new element (key2 should be evicted)
	// ja: 新しい要素を追加（key2 がエビクトされるはず）
	cache.Set("key6", "value6")

	// key2 should have been deleted
	// ja: key2 は削除されているはず
	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 should have been evicted")
	}

	// Other keys should still exist
	// ja: 他のキーは存在するはず
	for _, key := range []string{"key1", "key3", "key4", "key5", "key6"} {
		if _, ok := cache.Get(key); !ok {
			t.Errorf("❌: %s should still exist", key)
		}
	}
}

func TestLRU_UpdateMovesToFront(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxLength(3),
		cachez.WithEvictionPolicy(cachez.LRU),
	)

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Update key1 (move to front)
	// ja: key1 を更新（最前面に移動）
	cache.Set("key1", "updated")

	// Add new element (key2 should be evicted)
	// ja: 新しい要素を追加（key2 がエビクトされるはず）
	cache.Set("key4", "value4")

	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 should have been evicted")
	}

	if val, ok := cache.Get("key1"); !ok || val != "updated" {
		t.Errorf("❌: key1 should exist with updated value, got %v", val)
	}
}

func TestLRU_SetWithTTL_EdgeCases(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxLength(5),
		cachez.WithEvictionPolicy(cachez.LRU),
	)

	// Test updating existing entry with TTL = 0
	cache.SetWithTTL("key1", "value1", 100*time.Millisecond)
	cache.SetWithTTL("key1", "updated", 0) // Remove TTL

	// Wait longer than original TTL
	time.Sleep(150 * time.Millisecond)

	// Should still exist since TTL was removed
	if val, ok := cache.Get("key1"); !ok || val != "updated" {
		t.Error("❌: key1 should not expire after TTL was removed")
	}

	// Test updating existing entry without TTL to have TTL
	cache.Set("key2", "value2") // No TTL
	cache.SetWithTTL("key2", "updated2", 50*time.Millisecond)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	if _, ok := cache.Get("key2"); ok {
		t.Error("❌: key2 should have expired")
	}
}
