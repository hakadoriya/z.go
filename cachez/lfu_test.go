package cachez_test

import (
	"strconv"
	"testing"

	"github.com/hakadoriya/z.go/cachez"
)

func TestLFU_EvictionOrder(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](&cachez.Options{
		MaxCapacity:    5,
		EvictionPolicy: cachez.LFU,
	})

	// 5つの要素を追加
	for i := 1; i <= 5; i++ {
		cache.Set("key"+strconv.Itoa(i), "value"+strconv.Itoa(i))
	}

	// 異なる頻度でアクセス
	// key1: 4回
	for range 3 {
		cache.Get("key1")
	}
	// key2: 3回
	for range 2 {
		cache.Get("key2")
	}
	// key3: 2回
	cache.Get("key3")
	// key4: 1回（Set のみ）
	// key5: 1回（Set のみ）

	// 新しい要素を追加（key4 または key5 がエビクトされるはず）
	cache.Set("key6", "value6")

	// key4 または key5 のいずれかが削除されているはず
	key4Exists := false
	key5Exists := false

	if _, ok := cache.Get("key4"); ok {
		key4Exists = true
	}

	if _, ok := cache.Get("key5"); ok {
		key5Exists = true
	}

	if key4Exists && key5Exists {
		t.Error("either key4 or key5 should have been evicted")
	}

	// 頻度の高いキーは存在するはず
	for _, key := range []string{"key1", "key2", "key3"} {
		if _, ok := cache.Get(key); !ok {
			t.Errorf("%s should still exist", key)
		}
	}
}

func TestLFU_UpdateIncreasesFrequency(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](&cachez.Options{
		MaxCapacity:    3,
		EvictionPolicy: cachez.LFU,
	})

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// key1 と key2 の頻度を増加
	cache.Get("key1")
	cache.Get("key2")
	cache.Set("key1", "updated") // Set も頻度を増加させる

	// 新しい要素を追加（key3 がエビクトされるはず）
	cache.Set("key4", "value4")

	if _, ok := cache.Get("key3"); ok {
		t.Error("key3 should have been evicted")
	}

	if val, ok := cache.Get("key1"); !ok || val != "updated" {
		t.Errorf("key1 should exist with updated value, got %v", val)
	}
}
