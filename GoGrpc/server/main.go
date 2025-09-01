package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"server/middleware"
	"server/pb"

	"google.golang.org/grpc"
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

func main() {
	listen, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("new listen error: %s\n", err)
		return
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryServerInterceptorfunc),
		grpc.StreamInterceptor(middleware.StreamServerInterceptor),
	)
	pb.RegisterGreeterServer(server, &GreetServer{})
	err = server.Serve(listen)
	if err != nil {
		fmt.Printf("server error: %s\n", err)
		return
	}
}
