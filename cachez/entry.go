package cachez

import (
	"time"
)

// entryG represents a generic cache entry struct.
//
// ja: entryG はジェネリックキャッシュエントリを表す構造体です
type entryG[K comparable, V any] struct {
	key        K
	value      V
	expiration time.Time
	// Access frequency for LFU
	//
	// ja: LFU 用のアクセス頻度
	frequency int
}

// isExpired checks whether the entry has expired.
//
// ja: isExpired はエントリが期限切れかどうかをチェックします
func (e *entryG[K, V]) isExpired() bool {
	if e.expiration.IsZero() {
		return false
	}

	return time.Now().After(e.expiration)
}