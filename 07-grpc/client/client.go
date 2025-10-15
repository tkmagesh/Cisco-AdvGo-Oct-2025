package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tkmagesh/cisco-advgo-oct-2025/07-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	options := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.NewClient("localhost:50051", options)
	if err != nil {
		log.Fatalln(err)
	}
	clientPxy := proto.NewAppServiceClient(clientConn)
	ctx := context.Background()
	doRequestResponse(ctx, clientPxy)
}

func doRequestResponse(ctx context.Context, clientPxy proto.AppServiceClient) {
	// timeout after sometime
	/*
		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
	*/
	addReq := &proto.AddRequest{
		X: 100,
		Y: 200,
	}
	addRes, err := clientPxy.Add(ctx, addReq)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			fmt.Println("[Add] timed out!")
			log.Fatalln(err)
		}
	}
	fmt.Println("[Add Response] Result :", addRes.GetResult())
}
