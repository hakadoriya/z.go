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

// Cache is the generic cache interface.
//
// ja: Cache はジェネリックキャッシュインターフェースです
type Cache[K comparable, V any] interface {
	// Get retrieves the value for the specified key.
	//
	// ja: Get は指定されたキーの値を取得します
	Get(key K) (V, bool)
	// Set stores the specified key and value in the cache.
	//
	// ja: Set は指定されたキーと値をキャッシュに保存します
	Set(key K, value V)
	// SetWithTTL stores the specified key and value with TTL in the cache.
	//
	// ja: SetWithTTL は指定されたキーと値を TTL 付きでキャッシュに保存します
	SetWithTTL(key K, value V, ttl time.Duration)
	// Delete removes the specified key from the cache.
	//
	// ja: Delete は指定されたキーをキャッシュから削除します
	Delete(key K)
	// Clear removes all entries from the cache.
	//
	// ja: Clear はキャッシュの全エントリを削除します
	Clear()
	// Size returns the current number of cache entries.
	//
	// ja: Size は現在のキャッシュエントリ数を返します
	Size() int
}

// config represents internal cache configuration.
//
// ja: config は内部キャッシュ設定を表します
type config struct {
	maxCapacity    int
	evictionPolicy EvictionPolicy
	defaultTTL     time.Duration
}

// Option is a functional option for configuring the cache.
//
// ja: Option はキャッシュを設定するための関数オプションです
type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) {
	f(c)
}

// WithMaxCapacity sets the maximum capacity of the cache.
//
// ja: WithMaxCapacity はキャッシュの最大容量を設定します
func WithMaxCapacity(capacity int) Option {
	return optionFunc(func(c *config) {
		if capacity > 0 {
			c.maxCapacity = capacity
		}
	})
}

// WithEvictionPolicy sets the eviction policy.
//
// ja: WithEvictionPolicy はエビクションポリシーを設定します
func WithEvictionPolicy(policy EvictionPolicy) Option {
	return optionFunc(func(c *config) {
		c.evictionPolicy = policy
	})
}

// WithDefaultTTL sets the default TTL for cache entries.
//
// ja: WithDefaultTTL はキャッシュエントリのデフォルト TTL を設定します
func WithDefaultTTL(ttl time.Duration) Option {
	return optionFunc(func(c *config) {
		c.defaultTTL = ttl
	})
}

// New creates a new generic cache instance with functional options.
//
// ja: New は関数オプションを使用して新しいジェネリックキャッシュインスタンスを作成します
func New[K comparable, V any](opts ...Option) Cache[K, V] {
	// Default configuration
	cfg := &config{
		maxCapacity:    100, //nolint:mnd
		evictionPolicy: LRU,
		defaultTTL:     0, // No TTL by default // ja: デフォルトは無期限
	}

	// Apply options
	for _, opt := range opts {
		opt.apply(cfg)
	}

	// Create cache based on eviction policy
	switch cfg.evictionPolicy {
	case LFU:
		return newLFUCacheG[K, V](cfg)
	case LRU:
		fallthrough
	default:
		return newLRUCacheG[K, V](cfg)
	}
}
