package cachez

import (
	"container/heap"
	"sync"
	"time"
)

// lfuCache is a cache implementing the LFU eviction policy.
//
// ja: lfuCache は LFU エビクションポリシーを実装したキャッシュです
type lfuCache struct {
	mu          sync.RWMutex
	maxCapacity int
	defaultTTL  time.Duration
	items       map[string]*lfuEntry
	freqHeap    *minHeap
}

// lfuEntry represents an LFU cache entry.
//
// ja: lfuEntry は LFU キャッシュのエントリです
type lfuEntry struct {
	entry

	// Index in the heap
	//
	// ja: ヒープ内のインデックス
	index int
}

// newLFUCache creates a new LFU cache instance.
//
// ja: newLFUCache は新しい LFU キャッシュインスタンスを作成します
func newLFUCache(cfg *config) *lfuCache {
	mh := &minHeap{}
	heap.Init(mh)

	return &lfuCache{
		mu:          sync.RWMutex{},
		maxCapacity: cfg.maxCapacity,
		defaultTTL:  cfg.defaultTTL,
		items:       make(map[string]*lfuEntry),
		freqHeap:    mh,
	}
}

// Get retrieves the value for the specified key.
//
// ja: Get は指定されたキーの値を取得します
func (c *lfuCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ent, exists := c.items[key]
	if !exists {
		return nil, false
	}

	// TTL check
	//
	// ja: TTL チェック
	if ent.isExpired() {
		c.removeEntry(key)
		return nil, false
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
func (c *lfuCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores the specified key and value with TTL in the cache.
//
// ja: SetWithTTL は指定されたキーと値を TTL 付きでキャッシュに保存します
func (c *lfuCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
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
	ent := &lfuEntry{
		entry: entry{
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
func (c *lfuCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.items[key]; exists {
		c.removeEntry(key)
	}
}

// Clear removes all entries from the cache.
//
// ja: Clear はキャッシュの全エントリを削除します
func (c *lfuCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*lfuEntry)
	mh := &minHeap{}
	heap.Init(mh)
	c.freqHeap = mh
}

// Size returns the current number of cache entries.
//
// ja: Size は現在のキャッシュエントリ数を返します
func (c *lfuCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items)
}

// evictLeastFrequent removes the least frequently used entry.
//
// ja: evictLeastFrequent は最も使用頻度の低いエントリを削除します
func (c *lfuCache) evictLeastFrequent() {
	if c.freqHeap.Len() > 0 {
		entInterface := heap.Pop(c.freqHeap)

		ent, ok := entInterface.(*lfuEntry)
		if !ok {
			return
		}

		delete(c.items, ent.key)
	}
}

// removeEntry removes an entry.
//
// ja: removeEntry はエントリを削除します
func (c *lfuCache) removeEntry(key string) {
	ent := c.items[key]
	heap.Remove(c.freqHeap, ent.index)
	delete(c.items, key)
}

// minHeap is an implementation of a min-heap.
//
// ja: minHeap は最小ヒープの実装です
type minHeap []*lfuEntry

func (h *minHeap) Len() int { return len(*h) }

func (h *minHeap) Less(i, j int) bool {
	// When frequencies are equal, prioritize older entries
	// ja: 頻度が同じ場合は、より古いエントリを優先
	if (*h)[i].frequency == (*h)[j].frequency {
		return i < j
	}

	return (*h)[i].frequency < (*h)[j].frequency
}

func (h *minHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].index = i
	(*h)[j].index = j
}

func (h *minHeap) Push(x interface{}) {
	ent, ok := x.(*lfuEntry)
	if !ok {
		return
	}

	ent.index = len(*h)
	*h = append(*h, ent)
}

func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	ent := old[n-1]
	*h = old[0 : n-1]
	ent.index = -1

	return ent
}
