package cachez_test

import (
	"testing"
	"time"

	"github.com/hakadoriya/z.go/cachez"
)

func TestCache_LRU_Basic(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(3),
		cachez.WithEvictionPolicy(cachez.LRU),
	)

	// 基本的な Set/Get のテスト
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	val, ok := cache.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("❌: expected value1, got %v", val)
	}

	// 容量超過時のエビクションテスト
	cache.Set("key4", "value4") // key2 がエビクトされるはず（key1 は Get でアクセスされたため）

	if _, ok := cache.Get("key2"); ok {
		t.Error("key2 should have been evicted")
	}

	if _, ok := cache.Get("key1"); !ok {
		t.Error("key1 should still exist")
	}
}

func TestCache_LFU_Basic(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithMaxCapacity(3),
		cachez.WithEvictionPolicy(cachez.LFU),
	)

	// 基本的な Set/Get のテスト
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// key1 を複数回アクセス
	cache.Get("key1")
	cache.Get("key1")
	cache.Get("key2")

	// 容量超過時のエビクションテスト
	cache.Set("key4", "value4") // key3 がエビクトされるはず（最も使用頻度が低い）

	if _, ok := cache.Get("key3"); ok {
		t.Error("key3 should have been evicted")
	}

	if _, ok := cache.Get("key1"); !ok {
		t.Error("key1 should still exist")
	}
}

func TestCache_TTL(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string]() // デフォルトオプション

	// TTL 付きで値を設定
	cache.SetWithTTL("key1", "value1", 100*time.Millisecond)

	// 即座に取得
	if val, ok := cache.Get("key1"); !ok || val != "value1" {
		t.Errorf("❌: expected value1, got %v", val)
	}

	// TTL が経過するまで待機
	time.Sleep(150 * time.Millisecond)

	// TTL 経過後は取得できないはず
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should have expired")
	}
}

func TestCache_Delete(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string]()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.Delete("key1")

	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should have been deleted")
	}

	if val, ok := cache.Get("key2"); !ok || val != "value2" {
		t.Errorf("❌: key2 should still exist with value2, got %v", val)
	}
}

func TestCache_Clear(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string]()

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	if size := cache.Size(); size != 3 {
		t.Errorf("❌: expected size 3, got %d", size)
	}

	cache.Clear()

	if size := cache.Size(); size != 0 {
		t.Errorf("❌: expected size 0 after clear, got %d", size)
	}

	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should not exist after clear")
	}
}

func TestCache_DefaultTTL(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string](
		cachez.WithDefaultTTL(100 * time.Millisecond),
	)

	// デフォルト TTL で値を設定
	cache.Set("key1", "value1")

	// 即座に取得
	if val, ok := cache.Get("key1"); !ok || val != "value1" {
		t.Errorf("❌: expected value1, got %v", val)
	}

	// デフォルト TTL が経過するまで待機
	time.Sleep(150 * time.Millisecond)

	// TTL 経過後は取得できないはず
	if _, ok := cache.Get("key1"); ok {
		t.Error("key1 should have expired")
	}
}

func TestCache_UpdateExisting(t *testing.T) {
	t.Parallel()

	cache := cachez.New[string, string]()

	cache.Set("key1", "value1")
	cache.Set("key1", "value2") // 既存のキーを更新

	if val, ok := cache.Get("key1"); !ok || val != "value2" {
		t.Errorf("❌: expected updated value2, got %v", val)
	}

	// サイズは変わらないはず
	if size := cache.Size(); size != 1 {
		t.Errorf("❌: expected size 1, got %d", size)
	}
}
