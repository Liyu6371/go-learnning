package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"server/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type GreetServer struct {
	pb.UnimplementedGreeterServer
}

func (s *GreetServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

func (s *GreetServer) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	replays := []string{"AA", "BB", "CC"}
	for _, r := range replays {
		replay := "Hello_" + in.GetName() + "_" + r
		if err := stream.Send(&pb.HelloResponse{Reply: replay}); err != nil {
			fmt.Printf("stream Send replay: %s error: %s\n", replay, err)
			continue
		}
	}
	return nil
}

func (s *GreetServer) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
	replay := "Hello"
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.HelloResponse{Reply: replay})
		}
		if err != nil {
			fmt.Printf("stream Recv error: %s\n", err)
			return err
		}
		replay += "_" + req.GetName()
	}
}
func (s *GreetServer) LotsOfGreetingsAndReplies(stream pb.Greeter_LotsOfGreetingsAndRepliesServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream.Recv error: %v", err)
		}
		reply := "Hello_" + in.GetName() + "_By_GreetServer"
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return fmt.Errorf("stream.Send error: %v", err)
		}
	}
}

func UnaryServerInterceptorfunc(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get metadata from context")
	}
	v, ok := md["AuthToken"]
	if !ok {
		return nil, fmt.Errorf("AuthToken not provided in metadata")
	}
	if len(v) == 0 {
		return nil, fmt.Errorf("AuthToken is empty")
	}
	if v[0] != "TestAuthToken" {
		return nil, fmt.Errorf("invalid AuthToken: %s", v[0])
	}
	m, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func StreamServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return fmt.Errorf("failed to get metadata from context")
	}
	v, ok := md["AuthToken"]
	if !ok {
		return fmt.Errorf("AuthToken not provided in metadata")
	}
	if len(v) == 0 {
		return fmt.Errorf("AuthToken is empty")
	}
	if v[0] != "TestAuthToken" {
		return fmt.Errorf("invalid AuthToken: %s", v[0])
	}
	err := handler(srv, ss)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("new listen error: %s\n", err)
		return
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptorfunc),
		grpc.StreamInterceptor(StreamServerInterceptor),
	)
	pb.RegisterGreeterServer(server, &GreetServer{})
	err = server.Serve(listen)
	if err != nil {
		fmt.Printf("server error: %s\n", err)
		return
	}
}
