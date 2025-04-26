package realipz

import (
	"context"
	"net/netip"

	"github.com/hakadoriya/z.go/realipz"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type (
	newXRealIPConfig struct {
		clientIPAddressHeader string
	}

	NewXRealIPOption interface {
		apply(cfg *newXRealIPConfig)
	}

	newXRealIPServerInterceptorOptionClientIPAddressHeader string
)

func (f newXRealIPServerInterceptorOptionClientIPAddressHeader) apply(cfg *newXRealIPConfig) {
	cfg.clientIPAddressHeader = string(f)
}

func WithNewXRealIPServerInterceptorOptionClientIPAddressHeader(header string) NewXRealIPOption { //nolint:ireturn
	return newXRealIPServerInterceptorOptionClientIPAddressHeader(header)
}

// NewXRealIPSourceGRPCContext returns XRealIPSource from gRPC context.
//
// NOTE: Use this function only in special cases where you want to get X-Real-IP without using middleware or interceptor.
// NOTE: For normal use, you don't need to use this function externally. Consider using NewXRealIPMiddleware, NewXRealIPUnaryServerInterceptor, or NewXRealIPStreamServerInterceptor first.
func NewXRealIPSourceGRPCContext(ctx context.Context) (realipz.XRealIPSource, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrMetadataNotFound
	}

	var remoteAddr string
	if pr, ok := peer.FromContext(ctx); ok {
		remoteAddr = pr.Addr.String()
	}

	return &xRealIPSourceGRPCContext{
		ctx:        ctx,
		remoteAddr: remoteAddr,
		md:         md,
	}, nil
}

type xRealIPSourceGRPCContext struct {
	//nolint:containedctx // NOTE: ロガーを使いたいため
	ctx        context.Context
	remoteAddr string
	md         metadata.MD
}

func (r *xRealIPSourceGRPCContext) RemoteAddr() string {
	return r.remoteAddr
}

func (r *xRealIPSourceGRPCContext) Get(key string) string {
	s := r.md.Get(key)
	if len(s) == 0 {
		return ""
	}
	return s[0]
}

// NewXRealIPUnaryServerInterceptor returns a gRPC UnaryServerInterceptor.
//
// When set_real_ip_from is X-Forwarded-For and has a value like:
//
//	X-Forwarded-For: <SpoofingIP>, <ClientIP>, <ProxyIP>, <Proxy2IP>
//
// NewXRealIPUnaryServerInterceptor sets <ClientIP> as X-Real-IP header.
//
// NOTE: Parameter naming follows NGINX configuration naming conventions.
//
// Example usage:
//
//	realip.NewXRealIPUnaryServerInterceptor(
//		realip.DefaultSetRealIPFrom(),
//		realip.HeaderXForwardedFor,
//		true,
//	)
//
//nolint:revive,stylecheck
func NewXRealIPUnaryServerInterceptor(set_real_ip_from []netip.Prefix, real_ip_header string, real_ip_recursive bool, opts ...NewXRealIPOption) grpc.UnaryServerInterceptor {
	c := &newXRealIPConfig{
		clientIPAddressHeader: realipz.HeaderXRealIP,
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		xRealIPSource, err := NewXRealIPSourceGRPCContext(ctx)
		if err != nil {
			return handler(ctx, req)
		}

		xRealIP, err := realipz.XRealIP(xRealIPSource, set_real_ip_from, real_ip_header, real_ip_recursive)
		if err != nil {
			return handler(ctx, req)
		}

		if xRealIPSource.Get(c.clientIPAddressHeader) == "" {
			ctx = metadata.AppendToOutgoingContext(ctx, c.clientIPAddressHeader, xRealIP.String())
		}

		return handler(
			// add xRealIP to context
			realipz.WithContext(ctx, xRealIP),
			req,
		)
	}
}

// NewXRealIPStreamServerInterceptor returns a gRPC StreamServerInterceptor.
//
// When set_real_ip_from is X-Forwarded-For and has a value like:
//
//	X-Forwarded-For: <SpoofingIP>, <ClientIP>, <ProxyIP>, <Proxy2IP>
//
// NewXRealIPStreamServerInterceptor sets <ClientIP> as X-Real-IP header.
//
// NOTE: Parameter naming follows NGINX configuration naming conventions.
//
// Example usage:
//
//	realip.NewXRealIPStreamServerInterceptor(
//		realip.DefaultSetRealIPFrom(),
//		realip.HeaderXForwardedFor,
//		true,
//	)
//
//nolint:revive,stylecheck
func NewXRealIPStreamServerInterceptor(set_real_ip_from []netip.Prefix, real_ip_header string, real_ip_recursive bool, opts ...NewXRealIPOption) grpc.StreamServerInterceptor {
	c := &newXRealIPConfig{
		clientIPAddressHeader: realipz.HeaderXRealIP,
	}

	for _, opt := range opts {
		opt.apply(c)
	}

	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		xRealIPSource, err := NewXRealIPSourceGRPCContext(stream.Context())
		if err != nil {
			return handler(srv, stream)
		}

		xRealIP, err := realipz.XRealIP(xRealIPSource, set_real_ip_from, real_ip_header, real_ip_recursive)
		if err != nil {
			return handler(srv, stream)
		}

		ctx := stream.Context()
		if xRealIPSource.Get(c.clientIPAddressHeader) == "" {
			ctx = metadata.AppendToOutgoingContext(ctx, c.clientIPAddressHeader, xRealIP.String())
		}

		return handler(
			srv,
			&xRealIPStreamServerStream{
				ctx:          realipz.WithContext(ctx, xRealIP),
				ServerStream: stream,
			},
		)
	}
}

type xRealIPStreamServerStream struct {
	grpc.ServerStream

	//nolint:containedctx // NOTE: このパッケージは grpc の ServerStream をラップするため, メソッド Context() では context.Context を返す必要がある
	ctx context.Context
}

func (s *xRealIPStreamServerStream) Context() context.Context {
	return s.ctx
}
