package sqlz

import "bytes"

type testUser struct {
	*bytes.Buffer // Anonymous

	UserID     int     `testdb:"user_id"`
	Username   string  `testdb:"username"`
	NullString *string `testdb:"null_string"`
	Hyphen     string  `testdb:"-"`
	NoTag      string
}

var (
	_testUserTableName = "test_user"
	_testUserColumns   = []string{"user_id", "username", "null_string"}
)

func (*testUser) TableName() string { return _testUserTableName }
func (*testUser) Columns() []string { return _testUserColumns }

func copySlice[T any](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

func ptr[T any](v T) *T {
	return &v
}
