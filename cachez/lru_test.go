package cachez_test

import (
	"strconv"
	"testing"

	"github.com/hakadoriya/z.go/cachez"
)

func TestLRU_EvictionOrder(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](&cachez.Options{
		MaxCapacity:    5,
		EvictionPolicy: cachez.LRU,
	})

	// 5つの要素を追加
	for i := 1; i <= 5; i++ {
		cache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}

	// key1, key3, key5 にアクセス（使用順を変更）
	cache.Get("key1")
	cache.Get("key3")
	cache.Get("key5")

	// 新しい要素を追加（key2 がエビクトされるはず）
	cache.Set("key6", "value6")

	// key2 は削除されているはず
	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 should have been evicted")
	}

	// 他のキーは存在するはず
	for _, key := range []string{"key1", "key3", "key4", "key5", "key6"} {
		if _, ok := cache.Get(key); !ok {
			t.Errorf("%s should still exist", key)
		}
	}
}

func TestLRU_UpdateMovesToFront(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](&cachez.Options{
		MaxCapacity:    3,
		EvictionPolicy: cachez.LRU,
	})

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// key1 を更新（最前面に移動）
	cache.Set("key1", "updated")

	// 新しい要素を追加（key2 がエビクトされるはず）
	cache.Set("key4", "value4")

	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 should have been evicted")
	}

	if val, ok := cache.Get("key1"); !ok || val != "updated" {
		t.Errorf("key1 should exist with updated value, got %v", val)
	}
}
