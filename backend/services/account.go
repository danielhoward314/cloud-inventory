package services

import (
	"context"
	"errors"
	"fmt"

	accountpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/account"
	"golang.org/x/crypto/bcrypt"
)

func NewAccountService() accountpb.AccountServiceServer {
	return &accountService{}
}

// accountService implements the account gRPC service
type accountService struct {
	accountpb.UnimplementedAccountServiceServer
}

// hash hashes the given password using bcrypt
func hashPassword(password string) (string, error) {
	// Generate a hashed password with a default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// checkPasswordHash compares the given password with the hashed password.
func checkPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Signup accepts email + password or OAuth 2.0 social sign-in details for account creation
func (as *accountService) Signup(ctx context.Context, in *accountpb.SignupRequest) (*accountpb.SignupResponse, error) {
	hash, err := hashPassword(in.GetPassword())
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("bad")
	}

	return &accountpb.SignupResponse{
		Email:       in.GetEmail(),
		Password:    in.GetPassword(),
		Hash:        hash,
		HashMatches: checkPasswordHash(in.GetPassword(), hash),
	}, nil
}
