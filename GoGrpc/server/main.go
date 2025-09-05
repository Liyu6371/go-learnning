package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"server/config"
	"server/grpcserver"
	"sync"
	"syscall"
)

// raiseNumOfGrpcServer 拉起一定数量的 gRPC 服务器
func raiseNumOfGrpcServer(ctx context.Context, cfg *config.Config) {
	fmt.Println("Raising gRPC servers...")

	wg := sync.WaitGroup{}
	consulServer, err := grpcserver.NewConsulServer(&cfg.Consul)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, grpcCfg := range cfg.GrpcServers {
		grpcServerInst := grpcserver.NewGrpcGreeterServer(&grpcCfg)
		if grpcServerInst == nil {
			fmt.Printf("Failed to create gRPC server instance for %s\n", grpcCfg.Name)
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			grpcServerInst.Run(ctx, consulServer)
		}()
	}
	wg.Wait()
	fmt.Println("All gRPC servers have been shut down.")
}

func main() {
	cfg := config.DefaultConfig()
	if cfg == nil {
		fmt.Println("Failed to load default configuration")
		return
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		raiseNumOfGrpcServer(ctx, cfg)
	}()
	<-signalCh
	cancel()
	wg.Wait()
}
