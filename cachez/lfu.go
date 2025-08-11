package cachez

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrAssertionFailed = errors.New("assertion failed")

// lfuCacheG is a generic cache implementing the LFU eviction policy.
//
// ja: lfuCacheG はジェネリック版の LFU エビクションポリシーを実装したキャッシュです
type lfuCacheG[K comparable, V any] struct {
	mu          sync.RWMutex
	maxCapacity int
	defaultTTL  time.Duration
	items       map[K]*lfuEntryG[K, V]
	freqHeap    *minHeapG[K, V]
}

// lfuEntryG represents a generic LFU cache entry.
//
// ja: lfuEntryG はジェネリック版 LFU キャッシュのエントリです
type lfuEntryG[K comparable, V any] struct {
	entryG[K, V]

	// Index in the heap
	//
	// ja: ヒープ内のインデックス
	index int
}

// newLFUCacheG creates a new generic LFU cache instance.
//
// ja: newLFUCacheG は新しいジェネリック版 LFU キャッシュインスタンスを作成します
func newLFUCacheG[K comparable, V any](cfg *config) *lfuCacheG[K, V] {
	mh := &minHeapG[K, V]{}
	heap.Init(mh)

	return &lfuCacheG[K, V]{
		mu:          sync.RWMutex{},
		maxCapacity: cfg.maxCapacity,
		defaultTTL:  cfg.defaultTTL,
		items:       make(map[K]*lfuEntryG[K, V]),
		freqHeap:    mh,
	}
}

// Get retrieves the value for the specified key.
//
// ja: Get は指定されたキーの値を取得します
func (c *lfuCacheG[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	ent, exists := c.items[key]
	if !exists {
		return zero, false
	}

	// TTL check
	//
	// ja: TTL チェック
	if ent.isExpired() {
		c.removeEntry(key)
		return zero, false
	}

	// Increase access frequency
	//
	// ja: アクセス頻度を増加
	ent.frequency++
	heap.Fix(c.freqHeap, ent.index)

	return ent.value, true
}

// Set stores the specified key and value in the cache.
//
// ja: Set は指定されたキーと値をキャッシュに保存します
func (c *lfuCacheG[K, V]) Set(key K, value V) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores the specified key and value with TTL in the cache.
//
// ja: SetWithTTL は指定されたキーと値を TTL 付きでキャッシュに保存します
func (c *lfuCacheG[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check for existing entry
	// ja: 既存のエントリをチェック
	if ent, exists := c.items[key]; exists {
		// Update existing entry
		// ja: 既存エントリを更新
		ent.value = value

		ent.frequency++
		if ttl > 0 {
			ent.expiration = time.Now().Add(ttl)
		} else {
			ent.expiration = time.Time{}
		}

		heap.Fix(c.freqHeap, ent.index)

		return
	}

	// Capacity check and eviction
	// ja: 容量チェックとエビクション
	if len(c.items) >= c.maxCapacity {
		c.evictLeastFrequent()
	}

	// Create new entry
	// ja: 新規エントリを作成
	ent := &lfuEntryG[K, V]{
		entryG: entryG[K, V]{
			key:        key,
			value:      value,
			frequency:  1,
			expiration: time.Time{},
		},
		index: -1,
	}
	if ttl > 0 {
		ent.expiration = time.Now().Add(ttl)
	}

	// Add to heap
	// ja: ヒープに追加
	heap.Push(c.freqHeap, ent)
	c.items[key] = ent
}

// Delete removes the specified key from the cache.
//
// ja: Delete は指定されたキーをキャッシュから削除します
func (c *lfuCacheG[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.items[key]; exists {
		c.removeEntry(key)
	}
}

// Clear removes all entries from the cache.
//
// ja: Clear はキャッシュの全エントリを削除します
func (c *lfuCacheG[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]*lfuEntryG[K, V])
	mh := &minHeapG[K, V]{}
	heap.Init(mh)
	c.freqHeap = mh
}

// Size returns the current number of cache entries.
//
// ja: Size は現在のキャッシュエントリ数を返します
func (c *lfuCacheG[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// evictLeastFrequent removes the least frequently used entry.
//
// ja: evictLeastFrequent は最も使用頻度の低いエントリを削除します
func (c *lfuCacheG[K, V]) evictLeastFrequent() {
	if c.freqHeap.Len() > 0 {
		entInterface := heap.Pop(c.freqHeap)

		ent, ok := entInterface.(*lfuEntryG[K, V])
		if !ok || _testAssert_lfuCacheG_evictLeastFrequent {
			panic(fmt.Errorf("x is not a *lfuEntryG[K, V]: %w", ErrAssertionFailed))
		}

		delete(c.items, ent.key)
	}
}

// removeEntry removes an entry.
//
// ja: removeEntry はエントリを削除します
func (c *lfuCacheG[K, V]) removeEntry(key K) {
	ent := c.items[key]
	heap.Remove(c.freqHeap, ent.index)
	delete(c.items, key)
}

// minHeapG is a generic implementation of a min-heap.
//
// ja: minHeapG はジェネリック版最小ヒープの実装です
type minHeapG[K comparable, V any] []*lfuEntryG[K, V]

func (h *minHeapG[K, V]) Len() int { return len(*h) }

func (h *minHeapG[K, V]) Less(i, j int) bool {
	// When frequencies are equal, prioritize older entries
	// ja: 頻度が同じ場合は、より古いエントリを優先
	if (*h)[i].frequency == (*h)[j].frequency {
		return i < j
	}

	return (*h)[i].frequency < (*h)[j].frequency
}

func (h *minHeapG[K, V]) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].index = i
	(*h)[j].index = j
}

func (h *minHeapG[K, V]) Push(x interface{}) {
	ent, ok := x.(*lfuEntryG[K, V])
	if !ok || _testAssert_minHeapG_Push {
		panic(fmt.Errorf("x is not a *lfuEntryG[K, V]: %w", ErrAssertionFailed))
	}

	ent.index = len(*h)
	*h = append(*h, ent)
}

func (h *minHeapG[K, V]) Pop() interface{} {
	old := *h
	n := len(old)
	ent := old[n-1]
	*h = old[0 : n-1]
	ent.index = -1

	return ent
}

var (
	_testAssert_lfuCacheG_evictLeastFrequent bool
	_testAssert_minHeapG_Push                bool
)
