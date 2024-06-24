package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

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
	server := http.Server{
		Addr:         gatewayAddr,
		Handler:      corsHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start server in its own goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("gateway server failed to listen: %v", err)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until a signal is received
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	server.Shutdown(ctx)
}
