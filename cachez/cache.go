package cachez

import (
	"time"
)

// EvictionPolicy represents the cache eviction policy type.
//
// ja: EvictionPolicy はキャッシュのエビクションポリシーを表す型です
type EvictionPolicy int

const (
	// LRU is the Least Recently Used eviction policy.
	//
	// ja: LRU は Least Recently Used エビクションポリシーです
	LRU EvictionPolicy = iota
	// LFU is the Least Frequently Used eviction policy.
	//
	// ja: LFU は Least Frequently Used エビクションポリシーです
	LFU
)

// Cache is the cache interface.
//
// ja: Cache はキャッシュインターフェースです
type Cache interface {
	// Get retrieves the value for the specified key.
	//
	// ja: Get は指定されたキーの値を取得します
	Get(key string) (interface{}, bool)
	// Set stores the specified key and value in the cache.
	//
	// ja: Set は指定されたキーと値をキャッシュに保存します
	Set(key string, value interface{})
	// SetWithTTL stores the specified key and value with TTL in the cache.
	//
	// ja: SetWithTTL は指定されたキーと値を TTL 付きでキャッシュに保存します
	SetWithTTL(key string, value interface{}, ttl time.Duration)
	// Delete removes the specified key from the cache.
	//
	// ja: Delete は指定されたキーをキャッシュから削除します
	Delete(key string)
	// Clear removes all entries from the cache.
	//
	// ja: Clear はキャッシュの全エントリを削除します
	Clear()
	// Size returns the current number of cache entries.
	//
	// ja: Size は現在のキャッシュエントリ数を返します
	Size() int
}

// Options represents cache configuration options.
//
// ja: Options はキャッシュの設定オプションです
type Options struct {
	// MaxCapacity is the maximum capacity of the cache.
	//
	// ja: MaxCapacity はキャッシュの最大容量です
	MaxCapacity int
	// EvictionPolicy is the eviction policy.
	//
	// ja: EvictionPolicy はエビクションポリシーです
	EvictionPolicy EvictionPolicy
	// DefaultTTL is the default TTL.
	//
	// ja: DefaultTTL はデフォルトの TTL です
	DefaultTTL time.Duration
}

// DefaultOptions returns default options.
//
// ja: DefaultOptions はデフォルトのオプションを返します
func DefaultOptions() *Options {
	return &Options{
		MaxCapacity:    1000, //nolint:mnd
		EvictionPolicy: LRU,
		// No TTL by default
		//
		// ja: デフォルトの TTL は 0
		DefaultTTL: 0,
	}
}

// New creates a new cache instance.
//
// ja: New は新しいキャッシュインスタンスを作成します
func New(opts *Options) Cache {
	if opts == nil {
		opts = DefaultOptions()
	}

	if opts.MaxCapacity <= 0 {
		opts.MaxCapacity = 1000
	}

	switch opts.EvictionPolicy {
	case LFU:
		return newLFUCache(opts)
	case LRU:
		fallthrough
	default:
		return newLRUCache(opts)
	}
}
