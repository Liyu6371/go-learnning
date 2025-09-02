package grpcserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"server/middleware"
	"server/pb"
	"strings"

	"google.golang.org/grpc"
)

type GrpcGreeterServer struct {
	Name    string
	Address string
	Port    int
	pb.UnimplementedGreeterServer
}

func NewGrpcGreeterServer(name, addr string, port int) *GrpcGreeterServer {
	return &GrpcGreeterServer{
		Name:    name,
		Address: addr,
		Port:    port,
	}
}

func (s *GrpcGreeterServer) Run(ctx context.Context) {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Address, s.Port))
	if err != nil {
		fmt.Printf("failed to listen: %s", err)
		return
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryServerInterceptorfunc),
		grpc.StreamInterceptor(middleware.StreamServerInterceptor),
	)
	pb.RegisterGreeterServer(server, s)
	go func() {
		if err := server.Serve(listen); err != nil {
			if err == grpc.ErrServerStopped {
				fmt.Printf("grpc server: %s stopped: %s\n", s.Name, err)
				return
			}
			fmt.Printf("server %s failed to service: %s\n", s.Name, err)
		}
	}()
	<-ctx.Done()
	fmt.Printf(" gRPC server: %s catch ctx.Done signal\n", s.Name)
	server.GracefulStop()
	fmt.Printf(" gRPC server: %s stopped\n", s.Name)
}

func (s *GrpcGreeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	msg := "Hello_" + in.GetName() + "_By_Server_" + s.Name
	return &pb.HelloResponse{Reply: msg}, nil
}

func (s *GrpcGreeterServer) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	replies := []string{"Alice", "Bob", "Charlie"}
	for _, reply := range replies {
		msg := "Hello_" + in.GetName() + "_" + reply + "_By_Server_" + s.Name + "_" + reply
		if err := stream.Send(&pb.HelloResponse{Reply: msg}); err != nil {
			return err
		}
	}
	return nil
}

func (s *GrpcGreeterServer) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
	greetings := []string{}
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			msg := "Hello_" + strings.Join(greetings, "_") + "_By_Server_" + s.Name
			return stream.SendAndClose(&pb.HelloResponse{Reply: msg})
		}
		if err != nil {
			return fmt.Errorf("failed to receive greeting: %s", err)
		}
		greetings = append(greetings, in.GetName())
	}
}

func (s *GrpcGreeterServer) LotsOfGreetingsAndReplies(stream pb.Greeter_LotsOfGreetingsAndRepliesServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to receive greeting: %s", err)
		}
		msg := "Hello_" + in.GetName() + "_By_Server_" + s.Name
		if err := stream.Send(&pb.HelloResponse{Reply: msg}); err != nil {
			return fmt.Errorf("failed to send greeting: %s", err)
		}
	}
}
