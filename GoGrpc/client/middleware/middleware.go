package middleware

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// UnaryClientInterceptor is a gRPC client unary interceptor for adding metadata.
func UnaryClientInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return fmt.Errorf("failed to get metadata from outgoing context")
	}
	v, ok := md["auto_token"]
	if !ok || len(v) == 0 {
		return fmt.Errorf("auto_token not found in metadata")
	}
	fmt.Printf("auto_token from metadata: %s\n", v[0])
	return invoker(ctx, method, req, reply, cc, opts...)
}

func StreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get metadata from outgoing context")
	}
	v, ok := md["auto_token"]
	if !ok || len(v) == 0 {
		return nil, fmt.Errorf("auto_token not found in metadata")
	}
	fmt.Printf("auto_token from metadata: %s\n", v[0])
	return streamer(ctx, desc, cc, method, opts...)
}
