package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	accountpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/account"
	authpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/auth"
)

func main() {
	apiHost := os.Getenv("API_HOST")
	if len(apiHost) == 0 {
		apiHost = "localhost"
	}
	apiPort := os.Getenv("API_PORT")
	if len(apiPort) == 0 {
		apiPort = "50051"
	}
	apiAddr := apiHost + ":" + apiPort
	conn, err := grpc.Dial(apiAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to hello service: %v", err)
	}
	defer conn.Close()
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	err = accountpb.RegisterAccountServiceHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalf("failed to register the account server: %v", err)
	}
	err = authpb.RegisterAuthServiceHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatalf("failed to register the auth server: %v", err)
	}

	corsEnv := os.Getenv("CORS_ALLOW_LIST")
	if len(corsEnv) == 0 {
		corsEnv = "http://localhost:5173"
	}
	corsAllowList := strings.Split(corsEnv, ",")
	// Set up CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   corsAllowList,                                       // Allow only this origin
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"}, // Allow specific methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"},           // Allow specific headers
		AllowCredentials: true,                                                // Allow credentials
	}).Handler(mux)

	gatewayAddr := "[::]" + ":" + "8080"
	if err = http.ListenAndServe(gatewayAddr, corsHandler); err != nil {
		log.Fatal("gateway server closed abruptly: ", err)
	}
}
