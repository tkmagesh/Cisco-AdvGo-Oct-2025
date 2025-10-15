package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/tkmagesh/cisco-advgo-oct-2025/07-grpc/proto"
	"google.golang.org/grpc"
)

type AppServiceImpl struct {
	proto.UnimplementedAppServiceServer
}

// AppServiceServer interface implementation
func (asi *AppServiceImpl) Add(ctx context.Context, req *proto.AddRequest) (*proto.AddResponse, error) {

	// extract the data from the message payload
	x := req.GetX()
	y := req.GetY()

	// log the request
	log.Printf("[AppService.Add] x = %d & y = %d\n", x, y)

	// simulate a time consuming operation
	tickerCh := time.Tick(100 * time.Millisecond)
	result := y
	fmt.Print("[AppService.Add] processing .")
	for range x {
		select {
		case <-tickerCh:
			fmt.Print(".")
			result += 1
		case <-ctx.Done():
			fmt.Println("\n request timed out!")
			return nil, errors.New("timeout occurred")
		}
	}
	// process the request
	// result := x + y

	// create the response
	res := &proto.AddResponse{
		Result: result,
	}

	fmt.Println()
	fmt.Println("[AppService.Add] sending response...")
	// return the response
	return res, nil
}

func (asi *AppServiceImpl) GeneratePrimes(req *proto.PrimeRequest, serverStream proto.AppService_GeneratePrimesServer) error {
	start := req.GetStart()
	end := req.GetEnd()

	log.Printf("[AppService.GeneratePrimes] start = %d & end = %d\n", start, end)

	for no := start; no <= end; no++ {
		if isPrime(no) {
			resp := &proto.PrimeResponse{
				PrimeNo: no,
			}
			log.Printf("[AppService.GeneratePrimes] sending prime no : %d\n", no)
			if err := serverStream.Send(resp); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
	log.Println("[AppService.GeneratePrimes] All the prime numbers are sent!")
	return nil
}

func isPrime(no int64) bool {
	for i := int64(2); i <= (no / 2); i++ {
		if no%i == 0 {
			return false
		}
	}
	return true
}

func main() {
	asi := &AppServiceImpl{}
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterAppServiceServer(grpcServer, asi)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
