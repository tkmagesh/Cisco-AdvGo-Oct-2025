package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tkmagesh/cisco-advgo-oct-2025/07-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	options := grpc.WithTransportCredentials(insecure.NewCredentials())
	clientConn, err := grpc.NewClient("localhost:50051", options)
	if err != nil {
		log.Fatalln(err)
	}
	clientPxy := proto.NewAppServiceClient(clientConn)
	ctx := context.Background()
	addReq := &proto.AddRequest{
		X: 100,
		Y: 200,
	}
	addRes, err := clientPxy.Add(ctx, addReq)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("[Add Response] Result :", addRes.GetResult())

}
