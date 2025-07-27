package cachez

import (
	"errors"
	"testing"
)

// Test_lfuCacheG_evictLeastFrequent tests the type assertion failure case in evictLeastFrequent method.
// This test simulates a scenario where heap.Pop returns an unexpected type,
// which should trigger a panic with ErrAssertionFailed.
//
// ja: Test_lfuCacheG_evictLeastFrequent は evictLeastFrequent メソッドの型アサーション失敗ケースをテストします。
// このテストは heap.Pop が予期しない型を返すシナリオをシミュレートし、
// ErrAssertionFailed を伴うパニックが発生することを確認します。
//
//nolint:paralleltest
func Test_lfuCacheG_evictLeastFrequent(t *testing.T) {
	// Set up panic recovery to verify the expected error
	// ja: 期待されるエラーを検証するためのパニックリカバリを設定
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
		WithMaxLength(5),
		WithEvictionPolicy(LFU),
	)

	// Enable test assertion flag to simulate type assertion failure
	// ja: 型アサーション失敗をシミュレートするためのテストアサーションフラグを有効化
	before := _testAssert_lfuCacheG_evictLeastFrequent
	_testAssert_lfuCacheG_evictLeastFrequent = true
	t.Cleanup(func() { _testAssert_lfuCacheG_evictLeastFrequent = before })

	// Fill the cache to capacity
	// ja: キャッシュを容量まで埋める
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")
	cache.Set("key5", "value5")

	// This should trigger eviction and cause panic due to the test flag
	// ja: これによりエビクションが発生し、テストフラグによってパニックが起きる
	cache.Set("key6", "value6")
}

// Test_minHeapG_Push tests the type assertion failure case in minHeapG.Push method.
// This test simulates a scenario where Push receives an unexpected type,
// which should trigger a panic with ErrAssertionFailed.
//
// ja: Test_minHeapG_Push は minHeapG.Push メソッドの型アサーション失敗ケースをテストします。
// このテストは Push が予期しない型を受け取るシナリオをシミュレートし、
// ErrAssertionFailed を伴うパニックが発生することを確認します。
//
//nolint:paralleltest
func Test_minHeapG_Push(t *testing.T) {
	// Set up panic recovery to verify the expected error
	// ja: 期待されるエラーを検証するためのパニックリカバリを設定
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
		WithMaxLength(5),
		WithEvictionPolicy(LFU),
	)

	// Enable test assertion flag to simulate type assertion failure
	// ja: 型アサーション失敗をシミュレートするためのテストアサーションフラグを有効化
	before := _testAssert_minHeapG_Push
	_testAssert_minHeapG_Push = true
	t.Cleanup(func() { _testAssert_minHeapG_Push = before })

	// Fill the cache to capacity
	// ja: キャッシュを容量まで埋める
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")
	cache.Set("key4", "value4")
	cache.Set("key5", "value5")

	// This should trigger eviction and heap.Push, causing panic due to the test flag
	// ja: これによりエビクションと heap.Push が発生し、テストフラグによってパニックが起きる
	cache.Set("key6", "value6")
}
