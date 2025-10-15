package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"

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
	// doServerStreaming(ctx, clientPxy)
	doClientStreaming(ctx, clientPxy)
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

func doClientStreaming(ctx context.Context, clientPxy proto.AppServiceClient) {
	count := rand.Intn(20)
	fmt.Printf("[Aggreate] : sending %d numbers\n", count)

	clientStream, err := clientPxy.Aggregate(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	for range count {
		no := rand.Int63n(60)
		fmt.Printf("[Aggreate] : sending no : %d\n", no)
		req := &proto.AggregateRequest{
			No: no,
		}
		if err := clientStream.Send(req); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Printf("[Aggreate] : All data sent!\n")
	resp, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("[Aggreate] : Received aggregates!\n")
	fmt.Println("min :", resp.GetMin())
	fmt.Println("max :", resp.GetMax())
	fmt.Println("count :", resp.GetCount())
	fmt.Println("sum :", resp.GetSum())
	fmt.Println("average :", resp.GetAverage())

}
