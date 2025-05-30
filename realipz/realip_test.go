package realipz_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hakadoriya/z.go/realipz"
)

const testXForwardedFor = "127.0.0.1, 33.33.33.33, 10.1.1.1, 10.10.10.10, 10.100.100.100"

func TestXRealIP(t *testing.T) {
	t.Parallel()

	t.Run("failure,X-Forwarded-For_is_invalid_and_real_ip_recursive_off", func(t *testing.T) {
		t.Parallel()

		expect := "invalid IP"
		r, err := http.NewRequest(http.MethodPost, "http://realip", bytes.NewBufferString("test_request_body"))
		if err != nil {
			t.Fatalf("❌: err is not nil: %s", err)
		}
		r.Header.Set(realipz.HeaderXForwardedFor, "invalid")
		actual, err := realipz.XRealIP(realipz.NewXRealIPSourceHTTPRequest(r), realipz.DefaultSetRealIPFrom(), realipz.HeaderXForwardedFor, false)
		if err == nil {
			t.Errorf("❌: err is nil")
		}
		if expect != actual.String() {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})

	t.Run("failure,X-Forwarded-For_is_invalid_and_real_ip_recursive_on", func(t *testing.T) {
		t.Parallel()

		expect := "invalid IP"
		r, err := http.NewRequest(http.MethodPost, "http://realip", bytes.NewBufferString("test_request_body"))
		if err != nil {
			t.Fatalf("❌: err is not nil: %s", err)
		}
		r.Header.Set(realipz.HeaderXForwardedFor, "invalid")
		actual, err := realipz.XRealIP(realipz.NewXRealIPSourceHTTPRequest(r), realipz.DefaultSetRealIPFrom(), realipz.HeaderXForwardedFor, true)
		if err == nil {
			t.Errorf("❌: err is nil")
		}
		if expect != actual.String() {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})
}

func TestContextXRealIP(t *testing.T) {
	t.Parallel()
	expect := "invalid IP"
	actual := realipz.FromContext(context.Background())
	if expect != actual.String() {
		t.Errorf("❌: expect != actual: %s", actual)
	}
}

func TestNewXRealIPMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("success,normal", func(t *testing.T) {
		t.Parallel()

		header := http.CanonicalHeaderKey("X-Test-Real-IP")
		expect := "33.33.33.33"
		var actual string
		var actualFromCtx string
		actualResponse := &httptest.ResponseRecorder{}

		middleware := realipz.NewXRealIPMiddleware(
			realipz.DefaultSetRealIPFrom(),
			realipz.HeaderXForwardedFor,
			true,
			realipz.WithNewXRealIPOptionClientIPAddressHeader(header),
		)

		r := httptest.NewRequest(http.MethodPost, "http://realip", bytes.NewBufferString("test_request_body"))
		r.Header.Set(realipz.HeaderXForwardedFor, testXForwardedFor)

		middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get(header)
				actualFromCtx = realipz.FromContext(r.Context()).String()
			})).
			ServeHTTP(actualResponse, r)

		if expect != actualFromCtx {
			t.Errorf("❌: expect != actualFromCtx: %s", actual)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})

	t.Run("success,real_ip_header_is_not_X-Forwarded-For", func(t *testing.T) {
		t.Parallel()

		const testHeaderKey = "Test-Header-Key"

		expect := "33.33.33.33"
		var actual string
		var actualFromCtx string
		actualResponse := &httptest.ResponseRecorder{}

		middleware := realipz.NewXRealIPMiddleware(
			realipz.DefaultSetRealIPFrom(),
			testHeaderKey,
			true,
		)

		r := httptest.NewRequest(http.MethodPost, "http://realip", bytes.NewBufferString("test_request_body"))
		r.Header.Set(testHeaderKey, testXForwardedFor)

		middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get(realipz.HeaderXRealIP)
				actualFromCtx = realipz.FromContext(r.Context()).String()
			})).
			ServeHTTP(actualResponse, r)

		if expect != actualFromCtx {
			t.Errorf("❌: expect != actualFromCtx: %s", actual)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})

	t.Run("success,X-Forwarded-For_is_empty", func(t *testing.T) {
		t.Parallel()

		expect := "192.0.2.1"
		var actual string
		var actualFromCtx string
		actualResponse := &httptest.ResponseRecorder{}

		middleware := realipz.NewXRealIPMiddleware(
			realipz.DefaultSetRealIPFrom(),
			realipz.HeaderXForwardedFor,
			true,
		)

		r := httptest.NewRequest(http.MethodPost, "http://realip", bytes.NewBufferString("test_request_body"))
		r.Header.Set(realipz.HeaderXForwardedFor, "")

		middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get(realipz.HeaderXRealIP)
				actualFromCtx = realipz.FromContext(r.Context()).String()
			})).
			ServeHTTP(actualResponse, r)

		if expect != actualFromCtx {
			t.Errorf("❌: expect != actualFromCtx: %s", actual)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})

	t.Run("success,real_ip_recursive_off", func(t *testing.T) {
		t.Parallel()

		expect := "10.100.100.100"
		var actual string
		var actualFromCtx string
		actualResponse := &httptest.ResponseRecorder{}

		middleware := realipz.NewXRealIPMiddleware(
			realipz.DefaultSetRealIPFrom(),
			realipz.HeaderXForwardedFor,
			false,
		)

		r := httptest.NewRequest(http.MethodPost, "http://realip", bytes.NewBufferString("test_request_body"))
		r.Header.Set(realipz.HeaderXForwardedFor, testXForwardedFor)

		middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get(realipz.HeaderXRealIP)
				actualFromCtx = realipz.FromContext(r.Context()).String()
			})).
			ServeHTTP(actualResponse, r)

		if expect != actualFromCtx {
			t.Errorf("❌: expect != actualFromCtx: %s", actual)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})

	t.Run("failure,real_ip_recursive_off", func(t *testing.T) {
		t.Parallel()

		expect := "invalid IP"
		var actual string
		var actualFromCtx string
		actualResponse := &httptest.ResponseRecorder{}

		middleware := realipz.NewXRealIPMiddleware(
			realipz.DefaultSetRealIPFrom(),
			realipz.HeaderXForwardedFor,
			false,
		)

		r, err := http.NewRequest(http.MethodPost, "http://realip", bytes.NewBufferString("test_request_body"))
		if err != nil {
			t.Fatalf("❌: err is not nil: %s", err)
		}
		r.Header.Set(realipz.HeaderXForwardedFor, "invalid")

		middleware(http.HandlerFunc(
			func(rw http.ResponseWriter, r *http.Request) {
				actual = r.Header.Get(realipz.HeaderXRealIP)
				actualFromCtx = realipz.FromContext(r.Context()).String()
			})).
			ServeHTTP(actualResponse, r)

		if expect != actualFromCtx {
			t.Errorf("❌: expect != actualFromCtx: %s", actual)
		}
		if "" != actual {
			t.Errorf("❌: expect != actual: %s", actual)
		}
	})
}
