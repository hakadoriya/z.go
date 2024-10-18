//nolint:paralleltest
package envz

import (
	"strconv"
	"testing"
	"time"
)

func TestSecond(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 30 * time.Second
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(int64(expect.Seconds()), 10))
		actual, err := Second(TEST_ENV_KEY)
		if err != nil {
			t.Errorf("❌: Env: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure", func(t *testing.T) {
		const expect = 0
		t.Setenv(TEST_ENV_KEY, "test string")
		actual, err := Second(TEST_ENV_KEY)
		if err == nil {
			t.Errorf("❌: Env: err == nil")
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestSecondOrDefault(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 30 * time.Second
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(int64(expect.Seconds()), 10))
		actual := SecondOrDefault(TEST_ENV_KEY, 1)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("success(default)", func(t *testing.T) {
		const expect = 30 * time.Second
		actual := SecondOrDefault(TEST_ENV_KEY, expect)
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func TestMustSecond(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const expect = 30 * time.Second
		t.Setenv(TEST_ENV_KEY, strconv.FormatInt(int64(expect.Seconds()), 10))
		actual := MustSecond(TEST_ENV_KEY)
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
		_ = MustSecond(TEST_ENV_KEY)
	})
}
