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
	// doClientStreaming(ctx, clientPxy)
	doBidirectionalStream(ctx, clientPxy)
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

func doBidirectionalStream(ctx context.Context, clientPxy proto.AppServiceClient) {

	clientStream, err := clientPxy.Greet(ctx)

	if err != nil {
		log.Fatalln(err)
	}
	go sendRequests(ctx, clientStream)
	done := make(chan struct{})
	go func() {
		fmt.Println("Press ENTER to cancel")
		fmt.Scanln()
		clientStream.CloseSend()
		close(done)
	}()
	go recvResponse(ctx, clientStream)
	// return done
	<-done
}

func sendRequests(ctx context.Context, clientStream proto.AppService_GreetClient) {
	persons := []*proto.PersonName{
		{FirstName: "Magesh", LastName: "Kuppan"},
		{FirstName: "Suresh", LastName: "Kannan"},
		{FirstName: "Ramesh", LastName: "Jayaraman"},
		{FirstName: "Rajesh", LastName: "Pandit"},
		{FirstName: "Ganesh", LastName: "Kumar"},
	}

	// done := make(chan struct{})

	for _, person := range persons {
		req := &proto.GreetRequest{
			Person: person,
		}
		log.Printf("Sending Person : %s %s\n", person.FirstName, person.LastName)
		if err := clientStream.Send(req); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func recvResponse(ctx context.Context, clientStream proto.AppService_GreetClient) {
	for {
		res, err := clientStream.Recv()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(res.GetMessage())
	}
}
