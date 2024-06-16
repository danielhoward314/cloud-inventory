package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	hellopb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/hello"
)

func main() {
	helloServiceAddr := "localhost:50051"
	conn, err := grpc.Dial(helloServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to hello service: %v", err)
	}
	defer conn.Close()
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	err = hellopb.RegisterMyServiceHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalf("failed to register the order server: %v", err)
	}

	// Set up CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},         // Allow only this origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},  // Allow specific methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // Allow specific headers
		AllowCredentials: true,                                      // Allow credentials
	}).Handler(mux)
	// start listening to requests from the gateway server
	addr := ":8080"
	fmt.Println("API gateway server is running on " + addr)
	if err = http.ListenAndServe(addr, corsHandler); err != nil {
		log.Fatal("gateway server closed abruptly: ", err)
	}
}
