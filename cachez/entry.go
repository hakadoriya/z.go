package cachez

import (
	"time"
)

// entry represents a cache entry struct.
//
// ja: entry はキャッシュエントリを表す構造体です
type entry struct {
	key        string
	value      interface{}
	expiration time.Time
	// Access frequency for LFU
	//
	// ja: LFU 用のアクセス頻度
	frequency int
}

// isExpired checks whether the entry has expired.
//
// ja: isExpired はエントリが期限切れかどうかをチェックします
func (e *entry) isExpired() bool {
	if e.expiration.IsZero() {
		return false
	}

	return time.Now().After(e.expiration)
}
