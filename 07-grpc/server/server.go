package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/tkmagesh/cisco-advgo-oct-2025/07-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	fmt.Printf("[AppService.GeneratePrimes] start = %d & end = %d\n", start, end)

	for no := start; no <= end; no++ {
		if isPrime(no) {
			resp := &proto.PrimeResponse{
				PrimeNo: no,
			}
			fmt.Printf("[AppService.GeneratePrimes] sending prime no : %d\n", no)
			if err := serverStream.Send(resp); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
	fmt.Println("[AppService.GeneratePrimes] All the prime numbers are sent!")
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

func (asi *AppServiceImpl) Aggregate(serverStream proto.AppService_AggregateServer) error {
	fmt.Println("[AppService.Aggregate] delaying receiving data....")
	time.Sleep(5 * time.Second)
	var min, max, count, sum, avg int64
	min = math.MaxInt64
	for {
		req, err := serverStream.Recv()
		if err == io.EOF {
			avg = sum / count
			resp := &proto.AggregateResponse{
				Min:     min,
				Max:     max,
				Sum:     sum,
				Count:   count,
				Average: avg,
			}
			fmt.Println("[AppService.Aggregate] sending aggregate response")
			if err := serverStream.SendAndClose(resp); err != nil {
				log.Fatalln(err)
			}
			return nil
		}
		if err != nil {
			log.Fatalln(err)
		}
		no := req.GetNo()
		fmt.Println("[AppService.Aggregate] received no =", no)
		if min > no {
			min = no
		}
		if max < no {
			max = no
		}
		count += 1
		sum += no
	}
	return nil
}

func (asi *AppServiceImpl) Greet(serverStream proto.AppService_GreetServer) error {
	for {
		greetReq, err := serverStream.Recv()
		if code := status.Code(err); code == codes.Unavailable {
			fmt.Println("Client connection closed")
			break
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		person := greetReq.GetPerson()
		firstName := person.GetFirstName()
		lastName := person.GetLastName()
		log.Printf("Received greet request for %q and %q\n", firstName, lastName)
		message := fmt.Sprintf("Hi %s %s, Have a nice day!", firstName, lastName)
		time.Sleep(2 * time.Second)
		log.Printf("Sending response : %q\n", message)
		greetResp := &proto.GreetResponse{
			Message: message,
		}
		if err := serverStream.Send(greetResp); err != nil {
			if code := status.Code(err); code == codes.Unavailable {
				fmt.Println("Client connection closed")
				break
			}
		}
	}
	return nil
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
