package randz

import (
	crypto_rand "crypto/rand"
	"fmt"
	"io"
)

const DefaultStringReaderSource = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

//nolint:gochecknoglobals
var (
	StringReader = NewStringReader()
)

type stringReader struct {
	sourceString string
	randReader   io.Reader
}

type StringReaderOption interface {
	apply(r *stringReader)
}

type stringReaderOptionFunc func(r *stringReader)

func (f stringReaderOptionFunc) apply(r *stringReader) {
	f(r)
}

func WithStringReaderSource(source string) StringReaderOption {
	return stringReaderOptionFunc(func(r *stringReader) {
		r.sourceString = source
	})
}

func WithStringReaderRandReader(random io.Reader) StringReaderOption {
	return stringReaderOptionFunc(func(r *stringReader) {
		r.randReader = random
	})
}

func NewStringReader(opts ...StringReaderOption) io.Reader {
	r := &stringReader{
		sourceString: DefaultStringReaderSource,
		randReader:   crypto_rand.Reader,
	}

	for _, opt := range opts {
		opt.apply(r)
	}

	return r
}

func (r *stringReader) Read(p []byte) (n int, err error) {
	n, err = io.ReadFull(r.randReader, p)
	if err != nil {
		return n, fmt.Errorf("io.ReadFull: %w", err)
	}

	randomSourceLength := len(r.sourceString)
	for i := range p {
		p[i] = r.sourceString[int(p[i])%randomSourceLength]
	}

	return n, nil
}

func ReadString(random io.Reader, length int) (string, error) {
	b := make([]byte, length)

	if _, err := io.ReadFull(random, b); err != nil {
		return "", fmt.Errorf("io.ReadFull: %w", err)
	}

	return string(b), nil
}
