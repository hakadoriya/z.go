package realipz

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"slices"
	"strings"

	"github.com/hakadoriya/z.go/netipz"
)

const (
	HeaderXForwardedFor = "X-Forwarded-For"
	HeaderXRealIP       = "X-Real-IP"
)

func DefaultSetRealIPFrom() []netip.Prefix {
	return []netip.Prefix{
		netipz.LoopbackAddress,
		netipz.LinkLocalAddress,
		netipz.PrivateIPAddressClassA,
		netipz.PrivateIPAddressClassB,
		netipz.PrivateIPAddressClassC,
	}
}

type XRealIPSource interface {
	RemoteAddr() string
	Get(key string) string
}

// NewXRealIPSourceHTTPRequest returns XRealIPSource from http.Request.
//
// NOTE: Use this function only in special cases where you want to get X-Real-IP without using middleware or interceptor.
// NOTE: For normal use, you don't need to use this function externally. Consider using NewXRealIPMiddleware, NewXRealIPUnaryServerInterceptor, or NewXRealIPStreamServerInterceptor first.
func NewXRealIPSourceHTTPRequest(r *http.Request) XRealIPSource {
	return &xRealIPSourceHTTPRequest{r: r}
}

type xRealIPSourceHTTPRequest struct {
	r *http.Request
}

func (r *xRealIPSourceHTTPRequest) RemoteAddr() string {
	return r.r.RemoteAddr
}

func (r *xRealIPSourceHTTPRequest) Get(key string) string {
	s := r.r.Header.Get(key)
	return s
}

// XRealIP gets the X-Real-IP value from real_ip_header.
//
// NOTE: Use this function only in special cases where you want to get X-Real-IP without using middleware or interceptor.
// NOTE: For normal use, you don't need to use this function externally. Consider using NewXRealIPMiddleware, NewXRealIPUnaryServerInterceptor, or NewXRealIPStreamServerInterceptor first.
//
// When real_ip_header is X-Forwarded-For and has a value like:
//
//	X-Forwarded-For: <SpoofingIP>, <ClientIP>, <ProxyIP>, <Proxy2IP>
//
// XRealIP returns <ClientIP>.
//
// NOTE: Parameter naming follows NGINX configuration naming conventions.
//
// Example usage:
//
//	realip := realip.XRealIP(
//		r,
//		realip.DefaultSetRealIPFrom(),
//		realip.HeaderXForwardedFor,
//		true,
//	)
//
//nolint:revive,stylecheck
func XRealIP(r XRealIPSource, set_real_ip_from []netip.Prefix, real_ip_header string, real_ip_recursive bool) (netip.Addr, error) {
	xff := strings.Split(r.Get(real_ip_header), ",")

	// NOTE: If real_ip_recursive=off, return X-Forwarded-For tail value.
	if !real_ip_recursive {
		tail := strings.TrimSpace(xff[len(xff)-1])
		ip, err := netip.ParseAddr(tail)
		if err != nil {
			return netipz.Zero, fmt.Errorf("X-Forwarded-For tail=%s, netip.ParseAddr: %w", tail, err)
		}
		return ip, nil
	}

	xRealIP := netipz.Zero
	for idx := len(xff) - 1; idx >= 0; idx-- {
		ip, err := netip.ParseAddr(strings.TrimSpace(xff[idx]))
		// NOTE: If invalid ip, treat previous loop ip as X-Real-IP.
		if err != nil {
			break
		}

		xRealIP = ip

		// NOTE: If set_real_ip_from does not contain ip, treat this loop ip as X-Real-IP.
		if !slices.ContainsFunc(set_real_ip_from, func(prefix netip.Prefix) bool {
			return prefix.Contains(xRealIP)
		}) {
			break
		}
	}

	if xRealIP.IsValid() {
		return xRealIP, nil
	}

	// NOTE: If X-Forwarded-For is invalid csv that has invalid IP string, return RemoteAddr as X-Real-IP.
	ipStr, _, _ := net.SplitHostPort(r.RemoteAddr())
	ip, err := netip.ParseAddr(ipStr)
	if err != nil {
		return netipz.Zero, fmt.Errorf("remoteAddr=%s, netip.ParseAddr: %w", r.RemoteAddr(), err)
	}
	return ip, nil
}

type contextXRealIPKey struct{}

func FromContext(ctx context.Context) (xRealIP netip.Addr) {
	v, ok := ctx.Value(contextXRealIPKey{}).(netip.Addr)
	if ok {
		return v
	}

	return netipz.Zero
}

func WithContext(parent context.Context, xRealIP netip.Addr) context.Context {
	return context.WithValue(parent, contextXRealIPKey{}, xRealIP)
}

type (
	newXRealIPConfig struct {
		clientIPAddressHeader string
	}

	NewXRealIPOption interface {
		apply(cfg *newXRealIPConfig)
	}

	newXRealIPOptionClientIPAddressHeader string
)

func (f newXRealIPOptionClientIPAddressHeader) apply(cfg *newXRealIPConfig) {
	cfg.clientIPAddressHeader = string(f)
}

func WithNewXRealIPOptionClientIPAddressHeader(header string) NewXRealIPOption { //nolint:ireturn
	return newXRealIPOptionClientIPAddressHeader(header)
}

// NewXRealIPMiddleware returns a middleware that adds X-Real-IP header and sets X-Real-IP value in ctx.
//
// When set_real_ip_from is X-Forwarded-For and has a value like:
//
//	X-Forwarded-For: <SpoofingIP>, <ClientIP>, <ProxyIP>, <Proxy2IP>
//
// realip middleware sets <ClientIP> as X-Real-IP header.
//
// NOTE: Parameter naming follows NGINX configuration naming conventions.
//
// Example usage:
//
//	realip.NewXRealIPMiddleware(
//		realip.DefaultSetRealIPFrom(),
//		realip.HeaderXForwardedFor,
//		true,
//	)
//
//nolint:revive,stylecheck
func NewXRealIPMiddleware(set_real_ip_from []netip.Prefix, real_ip_header string, real_ip_recursive bool, opts ...NewXRealIPOption) func(next http.Handler) http.Handler {
	c := &newXRealIPConfig{
		clientIPAddressHeader: HeaderXRealIP,
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			xRealIP, err := XRealIP(NewXRealIPSourceHTTPRequest(r), set_real_ip_from, real_ip_header, real_ip_recursive)
			if err != nil {
				next.ServeHTTP(rw, r)
				return
			}

			if r.Header.Get(c.clientIPAddressHeader) == "" {
				r.Header.Set(c.clientIPAddressHeader, xRealIP.String())
			}

			next.ServeHTTP(rw, r.WithContext(
				// add xRealIP to context
				WithContext(r.Context(), xRealIP),
			))
		})
	}
}
