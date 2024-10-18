//nolint:paralleltest
package envz

import (
	"strconv"
	"testing"
)

func TestInt32(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 2000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual, err := Int32(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: Int32: %v", err)
		}
		if expect != int(actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := Int32(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: Int32: err == nil")
		}
		if expect != int(actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-format)", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := Int32(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: Int32: err == nil")
		}
		if expect != int(actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestInt32OrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 2000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual := Int32OrDefault(TEST_ENV_KEY, 1)
		if expect != int(actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 2000000000
		actual := Int32OrDefault(TEST_ENV_KEY, expect)
		if expect != int(actual) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustInt32(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 2000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual := MustInt32(TEST_ENV_KEY)
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
		_ = MustInt32(TEST_ENV_KEY)
	})
}
