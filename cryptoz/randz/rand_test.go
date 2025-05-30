package randz_test

import (
	"bytes"
	crypto_rand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"testing"

	randz "github.com/hakadoriya/z.go/cryptoz/randz"
)

var ExampleReader = bytes.NewBufferString("" +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production." +
	"This is a test Reader for example. Please use crypto/rand's Reader in production.",
)

func Example() {
	r := randz.NewStringReader(randz.WithStringReaderRandReader(ExampleReader))
	s, err := randz.ReadString(r, 128)
	if err != nil {
		log.Printf("(*randz.Reader).ReadString: %v", err)
		return
	}

	fmt.Printf("very secure random string: %s", s)
	// Output: very secure random string: Wqr1gr1gjg2n12gUnjmn0gox0gn6jvyunugSunj1ng31ngl07y2xv0jwmn1gUnjmn0grwgy0xm3l2rxwuWqr1gr1gjg2n12gUnjmn0gox0gn6jvyunugSunj1ng31ngl
}

func TestCreateCodeVerifier(t *testing.T) {
	t.Parallel()

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		actual, err := randz.ReadString(randz.StringReader, 128)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		t.Logf("✅: cv=%s", actual)

		backup := randz.StringReader
		t.Cleanup(func() { randz.StringReader = backup })
		randz.StringReader = bytes.NewBuffer(nil)
		if _, err := randz.ReadString(randz.StringReader, 128); err == nil {
			t.Errorf("❌: err == nil: %v", err)
		}
	})

	t.Run("success()", func(t *testing.T) {
		t.Parallel()
		const expect = "wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz012345wxyz0123"
		r := randz.NewStringReader(randz.WithStringReaderRandReader(bytes.NewBufferString(
			"01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567",
		)), randz.WithStringReaderSource(randz.DefaultStringReaderSource))
		actual, err := randz.ReadString(r, 128)
		if err != nil {
			t.Errorf("❌: err != nil: %v", err)
		}
		if expect != actual {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})

	t.Run("failure(io.EOF)", func(t *testing.T) {
		t.Parallel()
		r := randz.NewStringReader(randz.WithStringReaderRandReader(bytes.NewReader(nil)))
		_, actual := randz.ReadString(r, 128)
		if actual == nil {
			t.Errorf("❌: err == nil")
		}
		expect := io.EOF
		if !errors.Is(actual, expect) {
			t.Errorf("❌: expect != actual: %v != %v", expect, actual)
		}
	})
}

func BenchmarkGenerateRandomString(b *testing.B) {
	b.ResetTimer()

	b.Run("r.Read/mrand.Reader", func(b *testing.B) {
		r := randz.NewStringReader(randz.WithStringReaderRandReader(mrand.New(mrand.NewSource(0))))
		buf := make([]byte, 128)

		for range b.N {
			_, _ = r.Read(buf)
		}
	})

	b.Run("r.Read/crypto_rand.Reader", func(b *testing.B) {
		r := randz.NewStringReader(randz.WithStringReaderRandReader(crypto_rand.Reader))
		buf := make([]byte, 128)

		for range b.N {
			_, _ = r.Read(buf)
		}
	})
}
