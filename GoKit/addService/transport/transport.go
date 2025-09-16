package transport

import (
	"addService/pb"
	"context"

	"github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.UnimplementedAddServiceServer
	SumHandler    grpc.Handler
	ConcatHandler grpc.Handler
}

func (g *grpcServer) Sum(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	_, resp, err := g.SumHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.SumResponse), nil
}

func (g *grpcServer) Concat(ctx context.Context, req *pb.ConcatRequest) (*pb.ConcatResponse, error) {
	_, resp, err := g.ConcatHandler.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ConcatResponse), nil
}
