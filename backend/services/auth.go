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
	datastore      *dao.Datastore
	tokenDatastore dao.TokenDatastore
	smtpDialer     *gomail.Dialer
}

func NewAuthService(
	datastore *dao.Datastore,
	tokenDatastore dao.TokenDatastore,
	smtpDialer *gomail.Dialer,
) authpb.AuthServiceServer {
	return &authService{
		datastore:      datastore,
		tokenDatastore: tokenDatastore,
		smtpDialer:     smtpDialer,
	}
}

// ValidateSession validates admin ui session data submitted via a JWT in the request
func (as *authService) ValidateSession(ctx context.Context, request *authpb.ValidateSessionRequest) (*authpb.ValidateSessionResponse, error) {
	sessionTokenData, err := as.tokenDatastore.Read(request.Jwt)
	if err != nil {
		if err == redis.Nil {
			return nil, status.Errorf(codes.NotFound, "session data not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read session data: %s", err.Error())
	}
	administrator, err := as.datastore.Administrators.Read(sessionTokenData.AdministratorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "administrator not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read administrator data: %s", err.Error())
	}
	if !administrator.Verified {
		return nil, status.Errorf(codes.PermissionDenied, "email not verified")
	}
	err = as.tokenDatastore.Decode(ciJWT.Access, request.Jwt, ciJWT.AdminUISession)
	if err != nil {
		if err.Error() == ciJWT.TokenExpiredError {
			return nil, status.Errorf(codes.Unauthenticated, "access token has expired, use refresh token to request another")
		}
		if err.Error() == ciJWT.InvalidTokenError {
			return nil, status.Errorf(codes.PermissionDenied, "invalid access token")
		}
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
	adminUIAccessToken, err := as.tokenDatastore.Create(
		&dao.TokenData{
			AdministratorID:   administrator.ID,
			OrganizationID:    administrator.OrganizationID,
			AuthorizationRole: administrator.AuthorizationRole,
			TokenType:         ciJWT.Access,
			ClaimsType:        ciJWT.AdminUISession,
		},
		ciJWT.Access,
		ciJWT.AdminUISession,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err.Error())
	}
	adminUIRefreshToken, err := as.tokenDatastore.Create(
		&dao.TokenData{
			AdministratorID:   administrator.ID,
			OrganizationID:    administrator.OrganizationID,
			AuthorizationRole: administrator.AuthorizationRole,
			TokenType:         ciJWT.Refresh,
			ClaimsType:        ciJWT.AdminUISession,
		},
		ciJWT.Refresh,
		ciJWT.AdminUISession,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err.Error())
	}
	apiAccessToken, err := as.tokenDatastore.Create(
		&dao.TokenData{
			AdministratorID:   administrator.ID,
			OrganizationID:    administrator.OrganizationID,
			AuthorizationRole: administrator.AuthorizationRole,
			TokenType:         ciJWT.Access,
			ClaimsType:        ciJWT.APIAuthorization,
		},
		ciJWT.Access,
		ciJWT.APIAuthorization,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err.Error())
	}
	apiRefreshToken, err := as.tokenDatastore.Create(
		&dao.TokenData{
			AdministratorID:   administrator.ID,
			OrganizationID:    administrator.OrganizationID,
			AuthorizationRole: administrator.AuthorizationRole,
			TokenType:         ciJWT.Refresh,
			ClaimsType:        ciJWT.APIAuthorization,
		},
		ciJWT.Refresh,
		ciJWT.APIAuthorization,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session: %s", err.Error())
	}
	return &authpb.LoginResponse{
		AdministratorId:     administrator.ID,
		OrganizationId:      administrator.OrganizationID,
		AdministratorName:   administrator.DisplayName,
		OrganizationName:    organization.Name,
		BillingPlan:         organization.BillingPlanType,
		AdminUiAccessToken:  adminUIAccessToken,
		AdminUiRefreshToken: adminUIRefreshToken,
		ApiAccessToken:      apiAccessToken,
		ApiRefreshToken:     apiRefreshToken,
	}, nil
}

// RefreshToken takes in a refesh JWT of a given claims type and, if valid, returns a new access JWT of the same claims type
func (as *authService) RefreshToken(ctx context.Context, request *authpb.RefreshTokenRequest) (*authpb.RefreshTokenResponse, error) {
	if len(request.Jwt) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: request.jwt")
	}
	claimsType, err := ciJWT.GetClaimsTypeFromProtoEnum(request.ClaimsType)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid argument: %s", err.Error())
	}
	refreshTokenData, err := as.tokenDatastore.Read(request.Jwt)
	if err != nil {
		if err == redis.Nil {
			return nil, status.Errorf(codes.NotFound, "refresh token data not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read refresh token data: %s", err.Error())
	}
	administrator, err := as.datastore.Administrators.Read(refreshTokenData.AdministratorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "administrator not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read administrator data: %s", err.Error())
	}
	if !administrator.Verified {
		return nil, status.Errorf(codes.PermissionDenied, "email not verified")
	}
	err = as.tokenDatastore.Decode(ciJWT.Refresh, request.Jwt, claimsType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate refresh JWT: %s", err.Error())
	}
	// use the same tokenData as what was used for refresh token
	// tokenType hard-coded to access, since a refresh is always used to get an access token
	accessJWT, err := as.tokenDatastore.Create(refreshTokenData, ciJWT.Access, claimsType)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access JWT: %s", err.Error())
	}
	return &authpb.RefreshTokenResponse{Jwt: accessJWT}, nil
}
