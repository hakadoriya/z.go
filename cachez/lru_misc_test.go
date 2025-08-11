package cachez

import (
	"testing"
	"time"
)

// Test_lruCacheG_removeElement tests the type assertion failure case in removeElement method.
// This test simulates a scenario where elem.Value contains an unexpected type,
// which should be handled gracefully without panic.
//
// ja: Test_lruCacheG_removeElement は removeElement メソッドの型アサーション失敗ケースをテストします。
// このテストは elem.Value が予期しない型を含むシナリオをシミュレートし、
// パニックせずに適切に処理されることを確認します。
func Test_lruCacheG_removeElement(t *testing.T) {
	cache := New[string, string](
		WithMaxCapacity(5),
		WithEvictionPolicy(LRU),
	)

	// Enable test assertion flag to simulate type assertion failure
	// ja: 型アサーション失敗をシミュレートするためのテストアサーションフラグを有効化
	before := _testAssert_lruCacheG_removeElement
	_testAssert_lruCacheG_removeElement = true
	t.Cleanup(func() { _testAssert_lruCacheG_removeElement = before })

	// Fill the cache to capacity
	// ja: キャッシュを容量まで埋める
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")
	cache.Set("key5", "value5")
	
	// This should trigger eviction and call removeElement with the test flag
	// ja: これによりエビクションが発生し、テストフラグと共に removeElement が呼ばれる
	cache.Set("key6", "value6")
}

// Test_lruCacheG_Get_TypeAssertion tests the type assertion failure case in Get method.
// This test simulates a scenario where elem.Value contains an unexpected type,
// which should return false without panic.
//
// ja: Test_lruCacheG_Get_TypeAssertion は Get メソッドの型アサーション失敗ケースをテストします。
// このテストは elem.Value が予期しない型を含むシナリオをシミュレートし、
// パニックせずに false を返すことを確認します。
func Test_lruCacheG_Get_TypeAssertion(t *testing.T) {
	cache := New[string, string](
		WithMaxCapacity(5),
		WithEvictionPolicy(LRU),
	)

	// Enable test assertion flag to simulate type assertion failure
	// ja: 型アサーション失敗をシミュレートするためのテストアサーションフラグを有効化
	before := _testAssert_lruCacheG_Get
	_testAssert_lruCacheG_Get = true
	t.Cleanup(func() { _testAssert_lruCacheG_Get = before })

	// Set an entry
	// ja: エントリを設定
	cache.Set("key1", "value1")
	
	// Try to get it - should fail due to type assertion flag
	// ja: 取得を試みる - 型アサーションフラグにより失敗するはず
	val, ok := cache.Get("key1")
	if ok {
		t.Errorf("expected Get to fail with type assertion, got value: %v", val)
	}
}

// Test_lruCacheG_SetWithTTL_TypeAssertion tests the type assertion failure case in SetWithTTL method.
// This test simulates a scenario where elem.Value contains an unexpected type when updating an existing entry,
// which should be handled gracefully without updating the value.
//
// ja: Test_lruCacheG_SetWithTTL_TypeAssertion は SetWithTTL メソッドの型アサーション失敗ケースをテストします。
// このテストは既存エントリ更新時に elem.Value が予期しない型を含むシナリオをシミュレートし、
// 値を更新せずに適切に処理されることを確認します。
func Test_lruCacheG_SetWithTTL_TypeAssertion(t *testing.T) {
	cache := New[string, string](
		WithMaxCapacity(5),
		WithEvictionPolicy(LRU),
	)

	// Enable test assertion flag to simulate type assertion failure
	// ja: 型アサーション失敗をシミュレートするためのテストアサーションフラグを有効化
	before := _testAssert_lruCacheG_SetWithTTL
	_testAssert_lruCacheG_SetWithTTL = true
	t.Cleanup(func() { _testAssert_lruCacheG_SetWithTTL = before })

	// Set an entry first
	// ja: まずエントリを設定
	cache.Set("key1", "value1")
	
	// Try to update it with TTL - should fail silently due to type assertion flag
	// ja: TTL 付きで更新を試みる - 型アサーションフラグにより静かに失敗するはず
	cache.SetWithTTL("key1", "updated", 1*time.Hour)
	
	// Reset the flag to check if value was not updated
	// ja: 値が更新されていないことを確認するためフラグをリセット
	_testAssert_lruCacheG_SetWithTTL = false
	
	// Value should still be the original since update failed
	// ja: 更新が失敗したため、値は元のままであるはず
	val, ok := cache.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("expected original value 'value1', got: %v", val)
	}
}
