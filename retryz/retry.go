package retryz

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type (
	Jitter       func(duration time.Duration) (durationWithJitter time.Duration)
	jitterConfig struct {
		minJitter time.Duration
		maxJitter time.Duration
		rnd       *rand.Rand
	}
	JitterOption func(j *jitterConfig)
)

func WithDefaultJitterRange(minJitter, maxJitter time.Duration) JitterOption {
	return func(j *jitterConfig) {
		j.minJitter = minJitter
		j.maxJitter = maxJitter
	}
}

func WithDefaultJitterRand(rnd *rand.Rand) JitterOption {
	return func(j *jitterConfig) {
		j.rnd = rnd
	}
}

func DefaultJitter(opts ...JitterOption) Jitter {
	const (
		defaultMinJitter = 1 * time.Millisecond
		defaultMaxJitter = 100 * time.Millisecond
	)

	j := &jitterConfig{
		minJitter: defaultMinJitter,
		maxJitter: defaultMaxJitter,
		rnd:       nil,
	}

	for _, opt := range opts {
		opt(j)
	}

	return func(duration time.Duration) (durationWithJitter time.Duration) {
		if j.rnd == nil {
			return time.Duration(int64(duration) + rand.Int63n(int64(j.minJitter)+int64(j.maxJitter)) - int64(j.minJitter)) //nolint:gosec
		}
		return time.Duration(int64(duration) + j.rnd.Int63n(int64(j.minJitter)+int64(j.maxJitter)) - int64(j.minJitter))
	}
}

type Backoff func(initialInterval time.Duration, retries int) (intervalForThisRetry time.Duration)

func DefaultBackoff() Backoff {
	return func(initialInterval time.Duration, retries int) (intervalForThisRetry time.Duration) {
		return time.Duration(int64(initialInterval) << retries)
	}
}

type Config struct {
	initialInterval time.Duration
	maxInterval     time.Duration
	maxRetries      int
	timeout         time.Duration
	backoff         Backoff
	jitter          Jitter
}

const Infinite = -1

func NewConfig(initialInterval, maxInterval time.Duration, opts ...Option) *Config {
	const (
		defaultInitialInterval = 1 * time.Second
		defaultMaxInterval     = 30 * time.Second
	)

	if initialInterval == 0 {
		initialInterval = defaultInitialInterval
	}

	if maxInterval == 0 {
		maxInterval = defaultMaxInterval
	}

	c := &Config{
		initialInterval: initialInterval,
		maxInterval:     maxInterval,
		maxRetries:      Infinite,
		timeout:         0,
		backoff:         nil,
		jitter:          nil,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

type Option func(c *Config)

func WithBackoff(backoffFunc Backoff) Option {
	return func(c *Config) {
		c.backoff = backoffFunc
	}
}

func WithJitter(jitterFunc Jitter) Option {
	return func(c *Config) {
		c.jitter = jitterFunc
	}
}

func WithMaxRetries(maxRetries int) Option {
	return func(c *Config) {
		c.maxRetries = maxRetries
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *Config) {
		c.timeout = timeout
	}
}

// WARNING: Retryer should not be used across goroutines. Generate Retryer from Config for each goroutine.
type Retryer struct {
	ctx    context.Context //nolint:containedctx // WARNING: Retryer should not be used across goroutines. Generate Retryer from Config for each goroutine.
	cancel context.CancelFunc
	config *Config
	// variables
	interval time.Duration
	retries  int
	reason   error
}

func (c *Config) Build(ctx context.Context) *Retryer {
	var cancel context.CancelFunc
	if c.timeout > 0 {
		// NOTE: use context.WithTimeoutCause to distinguish from context.DeadlineExceeded
		ctx, cancel = context.WithTimeoutCause(ctx, c.timeout, ErrTimeoutExceeded)
	}
	copied := *c
	return &Retryer{
		ctx:      ctx,
		cancel:   cancel,
		config:   &copied,
		interval: 0,
		retries:  0,
		reason:   nil,
	}
}

// New returns *retry.Retryer. *retry.Retryer provides methods to facilitate retry execution.
//
// WARNING: Retryer should not be used across goroutines. Generate Retryer from Config for each goroutine.
//
// Is used as follows:
//
//	ctx := context.Background()
//	c := retry.NewConfig(10*time.Millisecond, 500*time.Millisecond, retry.WithMaxRetries(10))
//	r := retry.New(ctx, c)
//
//	for r.Retry() {
//		if r.Retries() == 0 {
//			fmt.Printf("FIRSTTRY! time=%s retries=%d retryAfter=%s\n", time.Now(), r.Retries(), r.RetryAfter())
//			continue
//		}
//		fmt.Printf("RETRYING! time=%s retries=%d retryAfter=%s\n", time.Now(), r.Retries(), r.RetryAfter())
//	}
//	fmt.Printf("COMPLETE! time=%s retries=%d error=%v\n", time.Now(), r.Retries(), r.Err())
//
// Then, yields the following results:
//
//	FIRSTTRY! time=2009-11-10 23:00:00 +0000 UTC m=+0.000000001 retries=0 retryAfter=10ms
//	RETRYING! time=2009-11-10 23:00:00.01 +0000 UTC m=+0.010000001 retries=1 retryAfter=20ms
//	RETRYING! time=2009-11-10 23:00:00.03 +0000 UTC m=+0.030000001 retries=2 retryAfter=40ms
//	RETRYING! time=2009-11-10 23:00:00.07 +0000 UTC m=+0.070000001 retries=3 retryAfter=80ms
//	RETRYING! time=2009-11-10 23:00:00.15 +0000 UTC m=+0.150000001 retries=4 retryAfter=160ms
//	RETRYING! time=2009-11-10 23:00:00.31 +0000 UTC m=+0.310000001 retries=5 retryAfter=320ms
//	RETRYING! time=2009-11-10 23:00:00.63 +0000 UTC m=+0.630000001 retries=6 retryAfter=500ms
//	RETRYING! time=2009-11-10 23:00:01.13 +0000 UTC m=+1.130000001 retries=7 retryAfter=500ms
//	RETRYING! time=2009-11-10 23:00:01.63 +0000 UTC m=+1.630000001 retries=8 retryAfter=500ms
//	RETRYING! time=2009-11-10 23:00:02.13 +0000 UTC m=+2.130000001 retries=9 retryAfter=500ms
//	RETRYING! time=2009-11-10 23:00:02.63 +0000 UTC m=+2.630000001 retries=10 retryAfter=500ms
//	COMPLETE! time=2009-11-10 23:00:02.63 +0000 UTC m=+2.630000001 retries=10 error=retry: reached max retries
//
// If the maximum count of attempts is not given via retry.WithMaxRetries(),
// *retry.Retryer that retry.New() returned will retry infinitely many times.
func New(ctx context.Context, config *Config) *Retryer {
	return config.Build(ctx)
}

func (r *Retryer) MaxRetries() (retries int) {
	return r.config.maxRetries
}

func (r *Retryer) Retries() (retries int) {
	return r.retries - 1
}

func (r *Retryer) RetryAfter() (retryAfter time.Duration) {
	return r.interval
}

func (r *Retryer) Err() (reason error) {
	return r.reason
}

const (
	ErrUnretryableErrorPrefix = "retryz: unretryable error"
)

var (
	ErrMaxRetriesExceeded = errors.New("retryz: max retries exceeded")
	ErrTimeoutExceeded    = errors.New("retryz: timeout exceeded")
)

func (r *Retryer) Retry() bool {
	if 0 <= r.MaxRetries() && r.MaxRetries() <= r.Retries() {
		r.reason = fmt.Errorf("maxRetries=%d: %w", r.config.maxRetries, ErrMaxRetriesExceeded)
		return false
	}

	select {
	case <-r.ctx.Done():
		if err := r.ctx.Err(); err != nil {
			r.reason = fmt.Errorf("ctx.Err: %w", err)
		}
		if err := context.Cause(r.ctx); err != nil {
			r.reason = fmt.Errorf("%w, context.Cause: %w", r.reason, err)
		}
		if r.cancel != nil {
			r.cancel()
		}
		return false
	case <-time.After(r.RetryAfter()):
		r.increment()
		return true
	}
}

type doConfig struct {
	errorHandler      func(ctx context.Context, r *Retryer, err error)
	unretryableErrors []error
	retryableErrors   []error
}

type DoOption interface {
	apply(c *doConfig)
}

type doOptionFunc func(c *doConfig)

func (f doOptionFunc) apply(c *doConfig) { f(c) }

func WithErrorHandler(f func(ctx context.Context, r *Retryer, err error)) DoOption {
	return doOptionFunc(func(c *doConfig) {
		c.errorHandler = f
	})
}

func WithUnretryableErrors(errs ...error) DoOption {
	return doOptionFunc(func(c *doConfig) {
		c.unretryableErrors = append(c.unretryableErrors, errs...)
	})
}

func WithRetryableErrors(errs ...error) DoOption {
	return doOptionFunc(func(c *doConfig) {
		c.retryableErrors = append(c.retryableErrors, errs...)
	})
}

//nolint:cyclop
func (r *Retryer) Do(f func(ctx context.Context) error, opts ...DoOption) error {
	c := new(doConfig)

	for _, opt := range opts {
		opt.apply(c)
	}

	var err error
LabelRetry:
	for r.Retry() {
		err = f(r.ctx)
		if errors.Is(err, nil) {
			return nil
		}
		if c.errorHandler != nil {
			c.errorHandler(r.ctx, r, err)
		}
		if len(c.unretryableErrors) > 0 {
			for _, unretryableErr := range c.unretryableErrors {
				if errors.Is(err, unretryableErr) {
					return fmt.Errorf(ErrUnretryableErrorPrefix+": %w", err)
				}
			}
			// continue LabelRetry NOTE: Do not continue here (considering the case where `c.retryableErrors` is set).
		}
		if len(c.retryableErrors) > 0 {
			for _, retryableErr := range c.retryableErrors {
				if errors.Is(err, retryableErr) {
					continue LabelRetry
				}
			}
			return fmt.Errorf(ErrUnretryableErrorPrefix+": %w", err)
		}
	}

	return fmt.Errorf("%s: %w", r.Err().Error(), err)
}

func (r *Retryer) getInitialInterval() time.Duration {
	return r.config.initialInterval
}

func (r *Retryer) truncateAtMaxInterval(d time.Duration) time.Duration {
	if d > r.config.maxInterval {
		return r.config.maxInterval
	}

	return d
}

func (r *Retryer) increment() {
	if r.config.backoff == nil {
		r.config.backoff = DefaultBackoff()
	}

	r.interval = r.truncateAtMaxInterval(r.config.backoff(r.getInitialInterval(), r.retries))

	if r.config.jitter == nil {
		r.config.jitter = DefaultJitter()
	}
	r.interval = r.config.jitter(r.interval)

	r.retries++
}
