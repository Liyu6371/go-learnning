package resolver

import "google.golang.org/grpc/resolver"

const (
	scheme   = "mygrpc"
	endpoint = "localgrpc.service.com"
)

var addrs = []string{
	"127.0.0.1:50051",
	"127.0.0.1:50052",
}

// type Builder interface {
// 	// Build creates a new resolver for the given target.
// 	//
// 	// gRPC dial calls Build synchronously, and fails if the returned error is
// 	// not nil.
// 	Build(target Target, cc ClientConn, opts BuildOptions) (Resolver, error)
// 	// Scheme returns the scheme supported by this resolver.  Scheme is defined
// 	// at https://github.com/grpc/grpc/blob/master/doc/naming.md.  The returned
// 	// string should not contain uppercase characters, as they will not match
// 	// the parsed target's scheme as defined in RFC 3986.
// 	Scheme() string
// }

type localGrpcResolver struct {
	target   resolver.Target
	cc       resolver.ClientConn
	addStore map[string][]string
}

func (l *localGrpcResolver) ResolveNow(o resolver.ResolveNowOptions) {
	addrStrs, ok := l.addStore[l.target.Endpoint()]
	if !ok {
		return
	}
	addrList := make([]resolver.Address, len(addrStrs))
	for i, v := range addrStrs {
		addrList[i] = resolver.Address{Addr: v}
	}
	l.cc.UpdateState(resolver.State{Addresses: addrList})
}

func (l *localGrpcResolver) Close() {}

// localGrpcResolverBuilder 自定义解析器 Builder
type localGrpcResolverBuilder struct{}

func (l *localGrpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &localGrpcResolver{
		target: target,
		cc:     cc,
		addStore: map[string][]string{
			endpoint: addrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (l *localGrpcResolverBuilder) Scheme() string {
	return scheme
}

func init() {
	resolver.Register(&localGrpcResolverBuilder{})
}
