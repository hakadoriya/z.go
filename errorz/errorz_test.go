package errorz

import (
	"errors"
	"io"
	"net"
	"regexp"
	"testing"
)

var errTestError = errors.New("testingz: test error")

func TestContains(t *testing.T) {
	t.Parallel()
	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		err := (error)(nil)
		if Contains(err, "testingz: test error") {
			t.Errorf("❌: err not contain %s: %v", "testingz: test error", err)
		}
	})

	t.Run("success,ErrTestError", func(t *testing.T) {
		t.Parallel()
		err := errTestError
		if !Contains(err, "testingz: test error") {
			t.Errorf("❌: err not contain `%s`: %v", "testingz: test error", err)
		}
	})
}

func TestHasPrefix(t *testing.T) {
	t.Parallel()
	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		err := (error)(nil)
		if HasPrefix(err, "testingz: test error") {
			t.Errorf("❌: err not has prefix %s: %v", "testingz: test error", err)
		}
	})

	t.Run("success,ErrTestError", func(t *testing.T) {
		t.Parallel()
		err := errTestError
		if !HasPrefix(err, "testingz: test error") {
			t.Errorf("❌: err not has prefix `%s`: %v", "testingz: test error", err)
		}
	})
}

func TestHasSuffix(t *testing.T) {
	t.Parallel()
	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		err := (error)(nil)
		if HasSuffix(err, "testingz: test error") {
			t.Errorf("❌: err not has suffix %s: %v", "testingz: test error", err)
		}
	})

	t.Run("success,ErrTestError", func(t *testing.T) {
		t.Parallel()
		err := errTestError
		if !HasSuffix(err, "testingz: test error") {
			t.Errorf("❌: err not has suffix `%s`: %v", "testingz: test error", err)
		}
	})
}

func TestMatchRegex(t *testing.T) {
	t.Parallel()
	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()
		err := (error)(nil)
		if MatchRegex(err, regexp.MustCompile("testingz: +test error")) {
			t.Errorf("❌: err not has suffix %s: %v", "testingz: test error", err)
		}
	})

	t.Run("success,ErrTestError", func(t *testing.T) {
		t.Parallel()
		err := errTestError
		if !MatchRegex(err, regexp.MustCompile("testingz: +test error")) {
			t.Errorf("❌: err not has suffix `%s`: %v", "testingz: test error", err)
		}
	})
}

type testTimeoutError struct {
	err     error
	timeout bool
}

func (e testTimeoutError) Error() string {
	return e.err.Error()
}

func (e testTimeoutError) Timeout() bool {
	return e.timeout
}

func (e testTimeoutError) Temporary() bool {
	return false
}

func TestIsNetTimeout(t *testing.T) {
	t.Parallel()

	t.Run("success,true", func(t *testing.T) {
		t.Parallel()
		err := testTimeoutError{err: net.ErrClosed, timeout: true}
		if !IsNetTimeout(err) {
			t.Errorf("❌: err is net timeout: %v", err)
		}
	})

	t.Run("success,false,net.Error", func(t *testing.T) {
		t.Parallel()
		err := testTimeoutError{err: net.ErrClosed, timeout: false}
		if IsNetTimeout(err) {
			t.Errorf("❌: err is net timeout: %v", err)
		}
	})

	t.Run("success,false,error", func(t *testing.T) {
		t.Parallel()
		err := io.EOF
		if IsNetTimeout(err) {
			t.Errorf("❌: err is not net timeout: %v", err)
		}
	})
}
