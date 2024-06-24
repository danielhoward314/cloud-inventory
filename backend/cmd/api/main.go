package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	"github.com/danielhoward314/cloud-inventory/backend/cmd/config"
	accountpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/account"
	"github.com/danielhoward314/cloud-inventory/backend/services"
)

func main() {
	// gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// postgres sql.DB instance
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	password := os.Getenv("POSTGRES_PASSWORD")
	sslMode := os.Getenv("POSTGRES_SSLMODE")
	user := os.Getenv("POSTGRES_USER")
	applicationDB := os.Getenv("POSTGRES_APPLICATION_DATABASE")
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		applicationDB,
		sslMode,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// redis client
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		log.Fatal("error: REDIS_HOST is empty")
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		log.Fatal("error: REDIS_PORT is empty")
	}
	redisAddr := redisHost + ":" + redisPort
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		DB:   0, // use default DB
	})

	// JWT secret for user sesssions
	jwtSecret := os.Getenv("JWT_SESSION_SECRET")
	if jwtSecret == "" {
		log.Fatal("error: JWT_SESSION_SECRET is empty")
	}

	// dependency injection
	cfg, err := config.NewAPIConfig(db, redisClient, jwtSecret)
	if err != nil {
		log.Fatalf(err.Error())
	}
	accountSvc := services.NewAccountService(cfg)

	// register service layer implementations for gRPC service interfaces
	accountpb.RegisterAccountServiceServer(s, accountSvc)

	// start gRPC server
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
