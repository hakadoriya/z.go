package retryz

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestRetryer_Retry(t *testing.T) {
	t.Parallel()

	t.Run("success,EXAMPLE", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 7
			maxInterval = 100 * time.Millisecond
		)
		ctx := context.Background()
		r := NewConfig(5*time.Millisecond, 100*time.Millisecond, WithMaxRetries(maxRetries)).Build(ctx)
		buf := bytes.NewBuffer(nil)
		for r.Retry() {
			fmt.Fprintf(buf, "[EXAMPLE] time=%s retries=%d/%d retryAfter=%s\n", time.Now(), r.Retries(), r.MaxRetries(), r.RetryAfter())
		}
		fmt.Fprintf(buf, "[EXAMPLE] time=%s If there is no difference of about %s in execution time between ↑ and ↓, it is OK.\n", time.Now(), maxInterval)
		fmt.Fprintf(buf, "[EXAMPLE] time=%s retries=%d/%d retryAfter=%s\n", time.Now(), r.Retries(), r.MaxRetries(), r.RetryAfter())
		t.Logf("✅: %s", buf)
	})

	t.Run("success,constant", func(t *testing.T) {
		t.Parallel()

		const maxRetries = 20
		ctx := context.Background()
		r := New(ctx, NewConfig(1*time.Microsecond, 999*time.Microsecond, WithMaxRetries(maxRetries), WithJitter(DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond), WithDefaultJitterRand(rand.New(rand.NewSource(0)))))))
		buf := bytes.NewBuffer(nil)
		for r.Retry() {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), r.MaxRetries(), r.RetryAfter())
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; retries=1/20 retryAfter=7.395152ms; retries=2/20 retryAfter=3.999827ms; retries=3/20 retryAfter=7.205794ms; retries=4/20 retryAfter=4.392202ms; retries=5/20 retryAfter=158.063µs; retries=6/20 retryAfter=5.044153ms; retries=7/20 retryAfter=6.550456ms; retries=8/20 retryAfter=10.150929ms; retries=9/20 retryAfter=3.149646ms; retries=10/20 retryAfter=1.942416ms; retries=11/20 retryAfter=7.975708ms; retries=12/20 retryAfter=10.258259ms; retries=13/20 retryAfter=1.884298ms; retries=14/20 retryAfter=10.98752ms; retries=15/20 retryAfter=7.115249ms; retries=16/20 retryAfter=4.980575ms; retries=17/20 retryAfter=4.528631ms; retries=18/20 retryAfter=1.917339ms; retries=19/20 retryAfter=6.163748ms; retries=20/20 retryAfter=10.440706ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Log("✅: actual: " + buf.String())
	})

	t.Run("success,jitter", func(t *testing.T) {
		t.Parallel()

		const maxRetries = 20
		ctx := context.Background()
		r := New(ctx, NewConfig(1*time.Microsecond, 999*time.Microsecond, WithMaxRetries(maxRetries), WithJitter(DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond), WithDefaultJitterRand(rand.New(rand.NewSource(0)))))))
		buf := bytes.NewBuffer(nil)
		for r.Retry() {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), r.MaxRetries(), r.RetryAfter())
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; retries=1/20 retryAfter=7.395152ms; retries=2/20 retryAfter=3.999827ms; retries=3/20 retryAfter=7.205794ms; retries=4/20 retryAfter=4.392202ms; retries=5/20 retryAfter=158.063µs; retries=6/20 retryAfter=5.044153ms; retries=7/20 retryAfter=6.550456ms; retries=8/20 retryAfter=10.150929ms; retries=9/20 retryAfter=3.149646ms; retries=10/20 retryAfter=1.942416ms; retries=11/20 retryAfter=7.975708ms; retries=12/20 retryAfter=10.258259ms; retries=13/20 retryAfter=1.884298ms; retries=14/20 retryAfter=10.98752ms; retries=15/20 retryAfter=7.115249ms; retries=16/20 retryAfter=4.980575ms; retries=17/20 retryAfter=4.528631ms; retries=18/20 retryAfter=1.917339ms; retries=19/20 retryAfter=6.163748ms; retries=20/20 retryAfter=10.440706ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
	})

	t.Run("failure,ErrMaxRetriesExceeded", func(t *testing.T) {
		t.Parallel()

		const maxRetries = 1
		jitter := DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond))
		ctx := context.Background()
		r := New(ctx, NewConfig(0, 0, WithMaxRetries(maxRetries), WithBackoff(DefaultBackoff()), WithJitter(jitter)))
		r.Retry() // first
		r.Retry() // second
		r.Retry() // third
		actual, expected := r.Err(), ErrMaxRetriesExceeded
		if !errors.Is(actual, expected) {
			t.Errorf("❌: err != `%s`: %v", expected, actual)
		}
	})

	t.Run("failure,context.Canceled", func(t *testing.T) {
		t.Parallel()

		const maxRetries = 1
		jitter := DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond))
		ctx, cancel := context.WithCancel(context.Background())
		t.Cleanup(cancel)
		r := New(ctx, NewConfig(0, 0, WithMaxRetries(maxRetries), WithBackoff(DefaultBackoff()), WithJitter(jitter)))
		r.Retry() // first
		cancel()
		r.Retry() // second
		r.Retry() // third
		actual, expected := r.Err(), context.Canceled
		if !errors.Is(actual, expected) {
			t.Errorf("❌: err != `%s`: %v", expected, actual)
		}
	})

	t.Run("failure,context.DeadlineExceeded", func(t *testing.T) {
		t.Parallel()

		const maxRetries = 1
		jitter := DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		t.Cleanup(cancel)
		r := New(ctx, NewConfig(0, 0, WithMaxRetries(maxRetries), WithBackoff(DefaultBackoff()), WithJitter(jitter)))
		r.Retry() // first
		time.Sleep(10 * time.Millisecond)
		r.Retry() // second
		r.Retry() // third
		actual, expected := r.Err(), context.DeadlineExceeded
		if !errors.Is(actual, expected) {
			t.Errorf("❌: err != `%s`: %v", expected, actual)
		}
	})

	t.Run("failure,ErrTimeoutExceeded", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		r := New(ctx, NewConfig(0, 0, WithTimeout(10*time.Microsecond)))
		r.Retry() // first
		// time.Sleep(1 * time.Millisecond)
		r.Retry() // second
		actual, expected := r.Err(), ErrTimeoutExceeded
		t.Logf("📝: actual: %v", actual)
		if !errors.Is(actual, expected) {
			t.Errorf("❌: err != `%s`: %v", expected, actual)
		}
	})
}

func TestRetryer_Do(t *testing.T) {
	t.Parallel()

	t.Run("success,nil", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		ctx := context.Background()
		r := New(ctx, NewConfig(1*time.Microsecond, 999*time.Microsecond, WithMaxRetries(maxRetries), WithJitter(DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond), WithDefaultJitterRand(rand.New(rand.NewSource(0)))))))
		buf := bytes.NewBuffer(nil)
		err := r.Do(func(_ context.Context) error {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), r.MaxRetries(), r.RetryAfter())
			return nil
		})
		if err != nil {
			t.Errorf("❌: err != nil")
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
	})

	t.Run("success,WithUnretryableErrors", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		ctx := context.Background()
		r := New(ctx, NewConfig(1*time.Microsecond, 999*time.Microsecond, WithMaxRetries(maxRetries), WithJitter(DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond), WithDefaultJitterRand(rand.New(rand.NewSource(0)))))))
		buf := bytes.NewBuffer(nil)
		err := r.Do(func(_ context.Context) error {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), r.MaxRetries(), r.RetryAfter())
			return io.EOF
		}, WithUnretryableErrors(io.ErrUnexpectedEOF))
		if !errors.Is(err, io.EOF) {
			t.Errorf("❌: err(%s) != nil(%s)", err, io.EOF)
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; retries=1/20 retryAfter=7.395152ms; retries=2/20 retryAfter=3.999827ms; retries=3/20 retryAfter=7.205794ms; retries=4/20 retryAfter=4.392202ms; retries=5/20 retryAfter=158.063µs; retries=6/20 retryAfter=5.044153ms; retries=7/20 retryAfter=6.550456ms; retries=8/20 retryAfter=10.150929ms; retries=9/20 retryAfter=3.149646ms; retries=10/20 retryAfter=1.942416ms; retries=11/20 retryAfter=7.975708ms; retries=12/20 retryAfter=10.258259ms; retries=13/20 retryAfter=1.884298ms; retries=14/20 retryAfter=10.98752ms; retries=15/20 retryAfter=7.115249ms; retries=16/20 retryAfter=4.980575ms; retries=17/20 retryAfter=4.528631ms; retries=18/20 retryAfter=1.917339ms; retries=19/20 retryAfter=6.163748ms; retries=20/20 retryAfter=10.440706ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
	})

	t.Run("success,WithRetryableErrors", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		ctx := context.Background()
		r := New(ctx, NewConfig(1*time.Microsecond, 999*time.Microsecond, WithMaxRetries(maxRetries), WithJitter(DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond), WithDefaultJitterRand(rand.New(rand.NewSource(0)))))))
		buf := bytes.NewBuffer(nil)
		err := r.Do(func(_ context.Context) error {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), r.MaxRetries(), r.RetryAfter())
			return io.EOF
		}, WithRetryableErrors(io.EOF))
		if !errors.Is(err, io.EOF) {
			t.Errorf("❌: err(%s) != nil(%s)", err, io.EOF)
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; retries=1/20 retryAfter=7.395152ms; retries=2/20 retryAfter=3.999827ms; retries=3/20 retryAfter=7.205794ms; retries=4/20 retryAfter=4.392202ms; retries=5/20 retryAfter=158.063µs; retries=6/20 retryAfter=5.044153ms; retries=7/20 retryAfter=6.550456ms; retries=8/20 retryAfter=10.150929ms; retries=9/20 retryAfter=3.149646ms; retries=10/20 retryAfter=1.942416ms; retries=11/20 retryAfter=7.975708ms; retries=12/20 retryAfter=10.258259ms; retries=13/20 retryAfter=1.884298ms; retries=14/20 retryAfter=10.98752ms; retries=15/20 retryAfter=7.115249ms; retries=16/20 retryAfter=4.980575ms; retries=17/20 retryAfter=4.528631ms; retries=18/20 retryAfter=1.917339ms; retries=19/20 retryAfter=6.163748ms; retries=20/20 retryAfter=10.440706ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
	})

	t.Run("failure,reachedMaxRetries", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		ctx := context.Background()
		r := New(ctx, NewConfig(1*time.Microsecond, 999*time.Microsecond, WithMaxRetries(maxRetries), WithJitter(DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond), WithDefaultJitterRand(rand.New(rand.NewSource(0)))))))
		buf := bytes.NewBuffer(nil)
		expected, actual := io.ErrUnexpectedEOF, r.Do(func(_ context.Context) error {
			return io.ErrUnexpectedEOF
		}, WithErrorHandler(func(_ context.Context, r *Retryer, err error) {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), r.MaxRetries(), r.RetryAfter())
		}))
		if !errors.Is(actual, expected) {
			t.Errorf("❌: err != `%s`: %v", expected, actual)
		}
		if !strings.Contains(actual.Error(), ErrMaxRetriesExceeded.Error()) {
			t.Errorf("❌: err not contain `%v`: %v", ErrMaxRetriesExceeded, actual)
		}
		const expectedContent = `retries=0/20 retryAfter=8.166505ms; retries=1/20 retryAfter=7.395152ms; retries=2/20 retryAfter=3.999827ms; retries=3/20 retryAfter=7.205794ms; retries=4/20 retryAfter=4.392202ms; retries=5/20 retryAfter=158.063µs; retries=6/20 retryAfter=5.044153ms; retries=7/20 retryAfter=6.550456ms; retries=8/20 retryAfter=10.150929ms; retries=9/20 retryAfter=3.149646ms; retries=10/20 retryAfter=1.942416ms; retries=11/20 retryAfter=7.975708ms; retries=12/20 retryAfter=10.258259ms; retries=13/20 retryAfter=1.884298ms; retries=14/20 retryAfter=10.98752ms; retries=15/20 retryAfter=7.115249ms; retries=16/20 retryAfter=4.980575ms; retries=17/20 retryAfter=4.528631ms; retries=18/20 retryAfter=1.917339ms; retries=19/20 retryAfter=6.163748ms; retries=20/20 retryAfter=10.440706ms; `
		actualContent := buf.String()
		if expectedContent != actualContent {
			t.Errorf("❌: expect(%s) != actual(%s)", expectedContent, actualContent)
		}
		t.Logf("✅: actual: %s", buf)
	})

	t.Run("failure,WithUnretryableErrors", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		ctx := context.Background()
		r := New(ctx, NewConfig(1*time.Microsecond, 999*time.Microsecond, WithMaxRetries(maxRetries), WithJitter(DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond), WithDefaultJitterRand(rand.New(rand.NewSource(0)))))))
		buf := bytes.NewBuffer(nil)
		err := r.Do(func(_ context.Context) error {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), r.MaxRetries(), r.RetryAfter())
			return io.ErrUnexpectedEOF
		}, WithUnretryableErrors(io.ErrUnexpectedEOF))
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		const expectErr = "retry: unretryable error: unexpected EOF"
		if !strings.Contains(err.Error(), expectErr) {
			t.Errorf("❌: err not contain: `%s` != `%v`", expectErr, err)
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
	})

	t.Run("failure,WithRetryableErrors", func(t *testing.T) {
		t.Parallel()

		const (
			maxRetries  = 20
			maxInterval = 100 * time.Millisecond
		)
		ctx := context.Background()
		r := New(ctx, NewConfig(1*time.Microsecond, 999*time.Microsecond, WithMaxRetries(maxRetries), WithJitter(DefaultJitter(WithDefaultJitterRange(1*time.Millisecond, 10*time.Millisecond), WithDefaultJitterRand(rand.New(rand.NewSource(0)))))))
		buf := bytes.NewBuffer(nil)
		err := r.Do(func(_ context.Context) error {
			fmt.Fprintf(buf, "retries=%d/%d retryAfter=%s; ", r.Retries(), r.MaxRetries(), r.RetryAfter())
			return io.ErrUnexpectedEOF
		}, WithRetryableErrors(io.EOF))
		if err == nil {
			t.Errorf("❌: err == nil")
		}
		const expectErr = "retry: unretryable error: unexpected EOF"
		if !strings.Contains(err.Error(), expectErr) {
			t.Errorf("❌: err not contain: `%s` != `%v`", expectErr, err)
		}
		const expect = `retries=0/20 retryAfter=8.166505ms; `
		actual := buf.String()
		if expect != actual {
			t.Errorf("❌: expect(%s) != actual(%s)", expect, actual)
		}
		t.Logf("✅: actual: %s", buf)
	})
}
