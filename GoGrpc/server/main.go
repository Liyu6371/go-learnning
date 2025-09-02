package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"server/grpcserver"
	"sync"
	"syscall"
)

func RaiseNumOfGrpcServer(ctx context.Context, num int, addr []string, port []int) {
	fmt.Println("Raising gRPC servers...")
	if len(addr) != num || len(port) != num {
		fmt.Println("Invalid input: addr and port slices must have the same length as num")
		return
	}
	wg := sync.WaitGroup{}
	for i := 0; i < num; i++ {
		serverAddr := addr[i]
		serverPort := port[i]
		serverName := fmt.Sprintf("grpc_server_%d", i+1)

		grpcServerInst := grpcserver.NewGrpcGreeterServer(serverName, serverAddr, serverPort)
		if grpcServerInst == nil {
			fmt.Printf("Failed to create gRPC server instance for %s\n", serverName)
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			grpcServerInst.Run(ctx)
		}()
	}
	wg.Wait()
	fmt.Println("All gRPC servers have been shut down.")
}

func main() {
	addr := []string{"127.0.0.1", "127.0.0.1"}
	port := []int{50051, 50052}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		RaiseNumOfGrpcServer(ctx, 2, addr, port)
	}()
	<-signalCh
	cancel()
	wg.Wait()
}
