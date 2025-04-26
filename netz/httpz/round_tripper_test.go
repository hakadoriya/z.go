package httpz_test

import (
	"net/http"
	"testing"

	"github.com/hakadoriya/z.go/netz/httpz"
	"github.com/hakadoriya/z.go/testingz"
	"github.com/hakadoriya/z.go/testingz/assertz"
)

func TestRoundTripFunc_RoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		a := httpz.RoundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return nil, testingz.ErrTestError
		})

		_, err := a.RoundTrip(nil) //nolint:bodyclose
		assertz.ErrorIs(t, err, testingz.ErrTestError)
	})
}
