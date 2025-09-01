package middleware

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// StreamServerInterceptor is a gRPC server stream interceptor for validating metadata.
func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return fmt.Errorf("failed to get metadata from context")
	}
	v, ok := md["auto_token"]
	if !ok {
		return fmt.Errorf("auto_token not provided in metadata")
	}
	if len(v) == 0 {
		return fmt.Errorf("auto_token is empty")
	}
	if v[0] != "test_auto_token" {
		return fmt.Errorf("invalid auto_token: %s", v[0])
	}
	err := handler(srv, ss)
	if err != nil {
		return err
	}
	return nil
}

// UnaryServerInterceptorfunc is a gRPC server unary interceptor for validating metadata.
func UnaryServerInterceptorfunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get metadata from context")
	}
	v, ok := md["auto_token"]
	if !ok {
		return nil, fmt.Errorf("auto_token not provided in metadata")
	}
	if len(v) == 0 {
		return nil, fmt.Errorf("auto_token is empty")
	}
	if v[0] != "test_auto_token" {
		return nil, fmt.Errorf("invalid auto_token: %s", v[0])
	}
	m, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	return m, nil
}
