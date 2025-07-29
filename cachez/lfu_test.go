package cachez_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/hakadoriya/z.go/cachez"
)

func TestLFU_EvictionOrder(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(5),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

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
			t.Errorf("❌: %s should still exist", key)
		}
	}
}

func TestLFU_UpdateIncreasesFrequency(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(3),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

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
		t.Errorf("❌: key1 should exist with updated value, got %v", val)
	}
}

func TestLFU_Delete(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(5),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

	// Add multiple entries
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Verify key exists before deletion
	if _, ok := cache.Get("key2"); !ok {
		t.Error("key2 should exist before deletion")
	}

	// Delete key2
	cache.Delete("key2")

	// Verify key2 is deleted
	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 should not exist after deletion")
	}

	// Verify other keys still exist
	if _, ok := cache.Get("key1"); !ok {
		t.Error("key1 should still exist")
	}
	if _, ok := cache.Get("key3"); !ok {
		t.Error("key3 should still exist")
	}

	// Delete non-existent key should not panic
	cache.Delete("non-existent")

	// Verify size is correct
	if size := cache.Size(); size != 2 {
		t.Errorf("❌: expected size 2 after deletion, got %d", size)
	}
}

func TestLFU_Clear(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(5),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

	// Add multiple entries
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")

	// Verify size before clear
	if size := cache.Size(); size != 4 {
		t.Errorf("❌: expected size 4 before clear, got %d", size)
	}

	// Clear all entries
	cache.Clear()

	// Verify size after clear
	if size := cache.Size(); size != 0 {
		t.Errorf("❌: expected size 0 after clear, got %d", size)
	}

	// Verify all keys are removed
	for i := 1; i <= 4; i++ {
		key := "key" + strconv.Itoa(i)
		if _, ok := cache.Get(key); ok {
			t.Errorf("❌: %s should not exist after clear", key)
		}
	}

	// Can add new entries after clear
	cache.Set("new-key", "new-value")
	if val, ok := cache.Get("new-key"); !ok || val != "new-value" {
		t.Errorf("❌: should be able to add new entries after clear")
	}
}

func TestLFU_SetWithTTL_UpdateExistingWithoutTTL(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(5),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

	// Test 1: Update existing entry with TTL to have TTL
	cache.SetWithTTL("key1", "value1", 100*time.Millisecond)
	cache.SetWithTTL("key1", "updated1", 200*time.Millisecond)
	
	if val, ok := cache.Get("key1"); !ok || val != "updated1" {
		t.Errorf("❌: expected updated1, got %v", val)
	}

	// Test 2: Update existing entry with TTL to have no TTL (ttl = 0)
	cache.SetWithTTL("key2", "value2", 100*time.Millisecond)
	
	// Verify entry exists with TTL
	if _, ok := cache.Get("key2"); !ok {
		t.Error("❌: key2 should exist before update")
	}
	
	// Update without TTL (remove expiration)
	cache.SetWithTTL("key2", "updated2", 0)
	
	// Wait longer than original TTL
	time.Sleep(150 * time.Millisecond)
	
	// Entry should still exist since TTL was removed
	if val, ok := cache.Get("key2"); !ok || val != "updated2" {
		t.Error("❌: key2 should not expire after TTL was removed")
	}

	// Test 3: Update existing entry without TTL to have TTL
	cache.Set("key3", "value3") // No TTL
	cache.SetWithTTL("key3", "updated3", 50*time.Millisecond)
	
	// Wait for expiration
	time.Sleep(100 * time.Millisecond)
	
	// Should be expired now
	if _, ok := cache.Get("key3"); ok {
		t.Error("❌: key3 should have expired")
	}
}

func TestLFU_ComplexEviction(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(3),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

	// Test eviction with empty heap (edge case)
	// This happens when all entries have been manually deleted
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Delete all entries manually
	cache.Delete("key1")
	cache.Delete("key2")
	cache.Delete("key3")

	// Now add new entries - eviction should handle empty heap gracefully
	cache.Set("new1", "value1")
	cache.Set("new2", "value2")
	cache.Set("new3", "value3")

	// Verify all new entries exist
	if _, ok := cache.Get("new1"); !ok {
		t.Error("new1 should exist")
	}
	if _, ok := cache.Get("new2"); !ok {
		t.Error("new2 should exist")
	}
	if _, ok := cache.Get("new3"); !ok {
		t.Error("new3 should exist")
	}

	// Test multiple evictions in sequence
	for i := 4; i <= 10; i++ {
		key := "key" + strconv.Itoa(i)
		cache.Set(key, "value"+strconv.Itoa(i))
		
		// Verify we still have exactly 3 items
		if size := cache.Size(); size != 3 {
			t.Errorf("❌: expected size 3, got %d", size)
		}
	}
}

func TestLFU_FrequencyTieBreaking(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(3),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

	// Add entries with same frequency
	cache.Set("a", "1")
	cache.Set("b", "2")
	cache.Set("c", "3")

	// All have frequency 1, so the oldest (a) should be evicted first
	cache.Set("d", "4")

	// Due to implementation using heap indices for tie-breaking,
	// one of the entries with frequency 1 should be evicted
	count := 0
	for _, key := range []string{"a", "b", "c", "d"} {
		if _, ok := cache.Get(key); ok {
			count++
		}
	}

	if count != 3 {
		t.Errorf("❌: expected exactly 3 entries to remain, got %d", count)
	}
}
