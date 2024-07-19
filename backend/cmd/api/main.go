package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	"gopkg.in/gomail.v2"

	ciPostgres "github.com/danielhoward314/cloud-inventory/backend/dao/postgres"
	ciRedis "github.com/danielhoward314/cloud-inventory/backend/dao/redis"
	accountpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/account"
	authpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/auth"
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

	// SMTP dialer instance
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		log.Fatal("error: SMTP_HOST is empty")
	}
	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		log.Fatal("error: SMTP_PORT is empty")
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatal("error: invalid SMTP_PORT")
	}
	smtpDialer := gomail.NewDialer(smtpHost, smtpPort, "", "")
	smtpDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// secret for JWT access token generation
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if accessTokenSecret == "" {
		log.Fatal("error: ACCESS_TOKEN_SECRET is empty")
	}

	// secret for JWT refresh token generation
	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if refreshTokenSecret == "" {
		log.Fatal("error: REFRESH_TOKEN_SECRET is empty")
	}

	datastore := ciPostgres.NewDatastore(db)
	registrationDatastore := ciRedis.NewRegistrationDatastore(redisClient)
	tokenDatastore := ciRedis.NewTokenDatastore(redisClient, accessTokenSecret, refreshTokenSecret)

	// dependency injection for each gRPC service
	accountSvc := services.NewAccountService(
		datastore,
		registrationDatastore,
		tokenDatastore,
		smtpDialer,
	)

	authSvc := services.NewAuthService(
		datastore,
		tokenDatastore,
		smtpDialer,
	)

	// register service layer implementations for gRPC service interfaces
	accountpb.RegisterAccountServiceServer(s, accountSvc)
	authpb.RegisterAuthServiceServer(s, authSvc)

	// start gRPC server
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
