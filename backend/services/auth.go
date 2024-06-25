package services

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/gomail.v2"

	"github.com/danielhoward314/cloud-inventory/backend/dao"
	ciJWT "github.com/danielhoward314/cloud-inventory/backend/jwt"
	"github.com/danielhoward314/cloud-inventory/backend/passwords"
	authpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/auth"
)

// authService implements the account gRPC service
type authService struct {
	authpb.UnimplementedAuthServiceServer
	datastore        *dao.Datastore
	sessionDatastore dao.SessionDatastore
	sessionJWTSecret string
	smtpDialer       *gomail.Dialer
}

func NewAuthService(
	datastore *dao.Datastore,
	sessionDataStore dao.SessionDatastore,
	sessionJWTSecret string,
	smtpDialer *gomail.Dialer,
) authpb.AuthServiceServer {
	return &authService{
		datastore:        datastore,
		sessionDatastore: sessionDataStore,
		sessionJWTSecret: sessionJWTSecret,
		smtpDialer:       smtpDialer,
	}
}

// ValidateSession validates user session data submitted via a JWT in the request
func (as *authService) ValidateSession(ctx context.Context, request *authpb.ValidateSessionRequest) (*authpb.ValidateSessionResponse, error) {
	_, err := as.sessionDatastore.Read(request.Jwt)
	if err != nil {
		if err == redis.Nil {
			return nil, status.Errorf(codes.NotFound, "session data not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read session data: %s", err.Error())
	}
	err = ciJWT.DecodeJWT(as.sessionJWTSecret, request.Jwt, ciJWT.Session)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate session JWT: %s", err.Error())
	}
	return &authpb.ValidateSessionResponse{
		Jwt: request.Jwt,
	}, nil
}

func (as *authService) Login(ctx context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	administrator, err := as.datastore.Administrators.ReadByEmail(request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "administrator not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read administrator data: %s", err.Error())
	}
	if !administrator.Verified {
		return nil, status.Errorf(codes.PermissionDenied, "email not verified")
	}
	err = passwords.ValidateBCryptHashedPassword(administrator.PasswordHash, request.Password)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "authentication error")
	}
	organization, err := as.datastore.Organizations.Read(administrator.OrganizationID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "organization not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read organization data: %s", err.Error())
	}
	sessionJWT, err := as.sessionDatastore.Create(&dao.Session{
		AdministratorID: administrator.ID,
		OrganizationID:  administrator.OrganizationID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err.Error())
	}
	return &authpb.LoginResponse{
		AdministratorId:   administrator.ID,
		OrganizationId:    administrator.OrganizationID,
		AdministratorName: administrator.DisplayName,
		OrganizationName:  organization.Name,
		BillingPlan:       organization.BillingPlanType,
		Jwt:               sessionJWT,
	}, nil
}
