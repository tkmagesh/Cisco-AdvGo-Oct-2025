package main

import (
	"context"
	"log"
	"net"

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

	// process the request
	result := x + y

	// create the response
	res := &proto.AddResponse{
		Result: result,
	}

	// return the response
	return res, nil
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
