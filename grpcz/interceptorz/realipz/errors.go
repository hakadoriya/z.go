package realipz

import "errors"

// ErrMetadataNotFound は, context 内にメタデータが見つからない場合のエラーです.
var ErrMetadataNotFound = errors.New("realip: metadata not found")
