package main

import (
	"context"
	"log"
	"net"

	hellopb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/hello"

	"google.golang.org/grpc"
)

// server is used to implement myservice.MyServiceServer.
type server struct {
	hellopb.MyServiceServer
}

// SayHello implements myservice.MyServiceServer
func (s *server) SayHello(ctx context.Context, in *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &hellopb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hellopb.RegisterMyServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
