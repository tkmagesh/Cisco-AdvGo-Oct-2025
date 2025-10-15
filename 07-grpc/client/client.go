package main

import (
	"context"
	"fmt"
	"io"
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
	// doRequestResponse(ctx, clientPxy)
	doServerStreaming(ctx, clientPxy)
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

func doServerStreaming(ctx context.Context, clientPxy proto.AppServiceClient) {
	primeReq := &proto.PrimeRequest{
		Start: 2,
		End:   100,
	}
	fmt.Println("[Server Streaming] Generate Primes - sending request")
	clientStream, err := clientPxy.GeneratePrimes(ctx, primeReq)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		resp, err := clientStream.Recv()
		if err == io.EOF {
			fmt.Println("[Server Streaming] All the responses received")
			return
		}
		if err != nil {
			log.Fatalln(err)
		}
		primeNo := resp.GetPrimeNo()
		fmt.Printf("[Server Streaming] Prime No : %d\n", primeNo)
	}

}
