//nolint:paralleltest
package envz

import (
	"strconv"
	"testing"
)

func TestInt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 2000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual, err := Int(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = 0
		actual, err := Int(TEST_ENV_KEY)
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
		actual, err := Int(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestIntOrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 2000000000
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(expect, 10))
		actual := IntOrDefault(TEST_ENV_KEY, 1)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 2000000000
		actual := IntOrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustInt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 2000000000
		t.Setenv(TEST_ENV_KEY, strconv.Itoa(expect))
		actual := MustInt(TEST_ENV_KEY)
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
		_ = MustInt(TEST_ENV_KEY)
	})
}
