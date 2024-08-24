//nolint:paralleltest
package envz

import (
	"strconv"
	"testing"
)

func TestBool(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = true
		t.Setenv(TEST_ENV_KEY, strconv.FormatBool(expect))
		actual, err := Bool(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(env-not-set)", func(t *testing.T) {
		const expect = false
		actual, err := Bool(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(fail-to-format)", func(t *testing.T) {
		const expect = false
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := Bool(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestBoolOrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = true
		t.Setenv(TEST_ENV_KEY, strconv.FormatBool(expect))
		actual := BoolOrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = true
		actual := BoolOrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustBool(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = true
		t.Setenv(TEST_ENV_KEY, strconv.FormatBool(expect))
		actual := MustBool(TEST_ENV_KEY)
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
		_ = MustBool(TEST_ENV_KEY)
	})
}
