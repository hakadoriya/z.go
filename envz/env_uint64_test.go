//nolint:paralleltest
package envz

import (
	"strconv"
	"testing"
)

func TestUint64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatUint(expect, 10))
		actual, err := Uint64(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := Uint64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-format)", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := Uint64(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestUint64OrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatUint(expect, 10))
		actual := Uint64OrDefault(TEST_ENV_KEY, 1)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 100000000000
		actual := Uint64OrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustUint64(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 100000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatUint(expect, 10))
		actual := MustUint64(TEST_ENV_KEY)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Errorf("❌: recover: err == nil")
			}
		}()
		_ = MustUint64(TEST_ENV_KEY)
	})
}
