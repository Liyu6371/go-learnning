package grpcserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"server/config"
	"server/middleware"
	"server/pb"
	"strings"

	"google.golang.org/grpc"
)

type GrpcGreeterServer struct {
	cfg *config.GrpcServerConfig
	pb.UnimplementedGreeterServer
}

func NewGrpcGreeterServer(c *config.GrpcServerConfig) *GrpcGreeterServer {
	return &GrpcGreeterServer{cfg: c}
}

func (s *GrpcGreeterServer) Run(ctx context.Context, cs *ConsulServer) {
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.Addr, s.cfg.Port))
	if err != nil {
		fmt.Printf("failed to listen: %s", err)
		return
	}
	defer listen.Close()
	server := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryServerInterceptorfunc),
		grpc.StreamInterceptor(middleware.StreamServerInterceptor),
	)
	pb.RegisterGreeterServer(server, s)
	go func() {
		if err := server.Serve(listen); err != nil {
			if err == grpc.ErrServerStopped {
				fmt.Printf("grpc server: %s stopped: %s\n", s.cfg.Name, err)
				return
			}
			fmt.Printf("server %s failed to service: %s\n", s.cfg.Name, err)
		}
	}()
	if cs == nil {
		fmt.Printf("Consul client is nil, skipping service registration for %s\n", s.cfg.Name)
		return
	}
	// Register service to Consul
	err = cs.RegisterService(s.cfg.Name, s.cfg.Addr, s.cfg.Port, []string{"grpc_greeter"})
	if err != nil {
		fmt.Printf("failed to register service to consul: %s\n", err)
		return
	}
	defer func() {
		if err := cs.DeregisterService(s.cfg.Name, s.cfg.Addr, s.cfg.Port); err != nil {
			fmt.Printf("failed to deregister service from consul: %s\n", err)
			return
		}
	}()
	<-ctx.Done()
	fmt.Printf(" gRPC server: %s catch ctx.Done signal\n", s.cfg.Name)
	server.GracefulStop()
	fmt.Printf(" gRPC server: %s stopped\n", s.cfg.Name)
}

func (s *GrpcGreeterServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	msg := "Hello_" + in.GetName() + "_By_Server_" + s.cfg.Name
	return &pb.HelloResponse{Reply: msg}, nil
}

func (s *GrpcGreeterServer) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	replies := []string{"Alice", "Bob", "Charlie"}
	for _, reply := range replies {
		msg := "Hello_" + in.GetName() + "_" + reply + "_By_Server_" + s.cfg.Name + "_" + reply
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
			msg := "Hello_" + strings.Join(greetings, "_") + "_By_Server_" + s.cfg.Name
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
		msg := "Hello_" + in.GetName() + "_By_Server_" + s.cfg.Name
		if err := stream.Send(&pb.HelloResponse{Reply: msg}); err != nil {
			return fmt.Errorf("failed to send greeting: %s", err)
		}
	}
}
