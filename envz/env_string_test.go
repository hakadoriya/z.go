//nolint:paralleltest
package envz

import (
	"testing"
)

const TEST_ENV_KEY = "TEST_ENV_KEY"

func TestString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = "test string"
		t.Setenv(TEST_ENV_KEY, expect)
		actual, err := String(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		const expect = ""
		actual, err := String(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestStringOrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = "test string"
		t.Setenv(TEST_ENV_KEY, expect)
		actual := StringOrDefault(TEST_ENV_KEY, "default")
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = "default"
		actual := StringOrDefault(TEST_ENV_KEY, "default")
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = "test string"
		t.Setenv(TEST_ENV_KEY, expect)
		actual := MustString(TEST_ENV_KEY)
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
		_ = MustString(TEST_ENV_KEY)
	})
}
