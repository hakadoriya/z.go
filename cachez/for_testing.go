package cachez

// for LFU
//
//nolint:gochecknoglobals,revive
var (
	_testAssert_lfuCacheG_evictLeastFrequent bool
	_testAssert_minHeapG_Push                bool
)

// for LRU
//
//nolint:gochecknoglobals,revive
var (
	_testAssert_lruCacheG_Get           bool
	_testAssert_lruCacheG_SetWithTTL    bool
	_testAssert_lruCacheG_removeElement bool
)
