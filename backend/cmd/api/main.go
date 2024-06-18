package main

import (
	"log"
	"net"

	accountpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/account"
	"github.com/danielhoward314/cloud-inventory/backend/services"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	accountSvc := services.NewAccountService()
	accountpb.RegisterAccountServiceServer(s, accountSvc)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
