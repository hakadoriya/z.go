package ioz_test

import (
	"bytes"
	"testing"

	ioz "github.com/hakadoriya/z.go/ioz"
	"github.com/hakadoriya/z.go/testingz/assertz"
)

func TestWriteFunc_Write(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		buf := new(bytes.Buffer)
		w := ioz.WriteFunc(func(p []byte) (n int, err error) {
			return buf.Write(append([]byte("prefix "), p...))
		})
		_, _ = w.Write([]byte("test"))
		assertz.Equal(t, "prefix test", buf.String())
	})
}
