package cachez

import (
	"container/list"
	"sync"
	"time"
)

// lruCache is a cache implementing the LRU eviction policy.
//
// ja: lruCache は LRU エビクションポリシーを実装したキャッシュです
type lruCache struct {
	mu          sync.RWMutex
	maxCapacity int
	defaultTTL  time.Duration
	items       map[string]*list.Element
	evictList   *list.List
}

// newLRUCache creates a new LRU cache instance.
//
// ja: newLRUCache は新しい LRU キャッシュインスタンスを作成します
func newLRUCache(opts *Options) *lruCache {
	return &lruCache{
		mu:          sync.RWMutex{},
		maxCapacity: opts.MaxCapacity,
		defaultTTL:  opts.DefaultTTL,
		items:       make(map[string]*list.Element),
		evictList:   list.New(),
	}
}

// Get retrieves the value for the specified key.
//
// ja: Get は指定されたキーの値を取得します
func (c *lruCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, exists := c.items[key]
	if !exists {
		return nil, false
	}

	ent, ok := elem.Value.(*entry)
	if !ok {
		return nil, false
	}

	if ent.isExpired() {
		c.removeElement(elem)
		return nil, false
	}

	// Move accessed entry to front
	// ja: アクセスされたエントリを最前面に移動
	c.evictList.MoveToFront(elem)

	return ent.value, true
}

// Set stores the specified key and value in the cache.
//
// ja: Set は指定されたキーと値をキャッシュに保存します
func (c *lruCache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores the specified key and value with TTL in the cache.
//
// ja: SetWithTTL は指定されたキーと値を TTL 付きでキャッシュに保存します
func (c *lruCache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check for existing entry
	if elem, exists := c.items[key]; exists {
		// Update existing entry
		ent, ok := elem.Value.(*entry)
		if !ok {
			return
		}

		ent.value = value
		if ttl > 0 {
			ent.expiration = time.Now().Add(ttl)
		} else {
			ent.expiration = time.Time{}
		}

		c.evictList.MoveToFront(elem)

		return
	}

	// Create new entry
	ent := &entry{
		key:        key,
		value:      value,
		expiration: time.Time{},
		frequency:  0,
	}
	if ttl > 0 {
		ent.expiration = time.Now().Add(ttl)
	}

	// Capacity check and eviction
	if c.evictList.Len() >= c.maxCapacity {
		c.evictOldest()
	}

	// Add new entry
	elem := c.evictList.PushFront(ent)
	c.items[key] = elem
}

// Delete removes the specified key from the cache.
//
// ja: Delete は指定されたキーをキャッシュから削除します
func (c *lruCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.removeElement(elem)
	}
}

// Clear removes all entries from the cache.
//
// ja: Clear はキャッシュの全エントリを削除します
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.evictList = list.New()
}

// Size returns the current number of cache entries.
//
// ja: Size は現在のキャッシュエントリ数を返します
func (c *lruCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.evictList.Len()
}

// evictOldest removes the oldest entry.
//
// ja: evictOldest は最も古いエントリを削除します
func (c *lruCache) evictOldest() {
	elem := c.evictList.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}

// removeElement removes an entry from the list and map.
//
// ja: removeElement はリストとマップからエントリを削除します
func (c *lruCache) removeElement(elem *list.Element) {
	c.evictList.Remove(elem)

	ent, ok := elem.Value.(*entry)
	if !ok {
		return
	}

	delete(c.items, ent.key)
}
