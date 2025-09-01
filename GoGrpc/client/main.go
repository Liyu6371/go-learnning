package main

import (
	"client/middleware"
	"client/pb"
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func ProcessMd(ctx context.Context) context.Context {
	subCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs("auto_token", "test_auto_token"))
	return subCtx
}

func RunSayHello(ctx context.Context, client pb.GreeterClient) {
	// fmt.Println("---- 单一调用 ----")
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		fmt.Printf("call SayHello error: %s\n", err)
		return
	}
	fmt.Printf("response from server: %s\n", r.GetReply())
}

func RunLotsOfReplies(ctx context.Context, client pb.GreeterClient) {
	fmt.Println("---- 服务端流式传输 ----")
	// 服务端流式传输
	replyStream, err := client.LotsOfReplies(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil {
		fmt.Printf("call LotsOfReplies error: %s\n", err)
		return
	}
	for {
		reply, err := replyStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("call LotsOfReplies error: %s\n", err)
			return
		}
		fmt.Printf("response from server: %s\n", reply.GetReply())
	}
}

func RunLotsOfGreetings(ctx context.Context, client pb.GreeterClient) {
	fmt.Println("---- 客户端流式传输 ----")
	// 客户端流式传输
	reqStream, err := client.LotsOfGreetings(ctx)
	if err != nil {
		fmt.Printf("call LotsOfGreetings error: %s\n", err)
		return
	}
	names := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	for _, name := range names {
		if err := reqStream.Send(&pb.HelloRequest{Name: name}); err != nil {
			fmt.Printf("send request error: %s\n", err)
			continue
		}
	}
	resp, err := reqStream.CloseAndRecv()
	if err != nil {
		fmt.Printf("receive response error: %s\n", err)
		return
	}
	fmt.Printf("response from server: %s\n", resp.GetReply())
}

func RunLotsOfGreetingsAndReplies(ctx context.Context, client pb.GreeterClient) {
	names := []string{"Alice", "Bob", "Charlie", "David", "Eve"}
	fmt.Println("---- 双向流式传输 ----")
	reqStream2, err := client.LotsOfGreetingsAndReplies(ctx)
	if err != nil {
		fmt.Printf("call LotsOfGreetingsAndReplies error: %s\n", err)
		return
	}
	waitC := make(chan struct{})
	go func() {
		for {
			in, err := reqStream2.Recv()
			if err == io.EOF {
				close(waitC)
				return
			}
			if err != nil {
				fmt.Printf("Failed to receive a note : %s\n", err)
				close(waitC)
				return
			}
			fmt.Printf("Got message from server: %s\n", in.GetReply())
		}
	}()
	for _, name := range names {
		if err := reqStream2.Send(&pb.HelloRequest{Name: name}); err != nil {
			fmt.Printf("Failed to send a note : %s\n", err)
			continue
		}
	}
	reqStream2.CloseSend()
	<-waitC
}

func main() {
	conn, err := grpc.NewClient(
		"localhost:8972",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.UnaryClientInterceptor),
		grpc.WithStreamInterceptor(middleware.StreamClientInterceptor),
	)
	if err != nil {
		fmt.Printf("new grpc client error: %s\n", err)
		return
	}
	defer conn.Close()
	ctx := ProcessMd(context.Background())
	client := pb.NewGreeterClient(conn)
	RunSayHello(ctx, client)
	RunLotsOfReplies(ctx, client)
	RunLotsOfGreetings(ctx, client)
	RunLotsOfGreetingsAndReplies(ctx, client)
}
