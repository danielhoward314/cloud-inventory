package services

import (
	"bytes"
	"context"
	"database/sql"
	"html/template"
	"log/slog"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/gomail.v2"

	"github.com/danielhoward314/cloud-inventory/backend/dao"
	"github.com/danielhoward314/cloud-inventory/backend/dao/postgres"
	ciJWT "github.com/danielhoward314/cloud-inventory/backend/jwt"
	accountpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/account"
)

const (
	emailFrom = "no.reply@CloudInventory.com"
)

// accountService implements the account gRPC service
type accountService struct {
	accountpb.UnimplementedAccountServiceServer
	datastore             *dao.Datastore
	registrationDatastore dao.RegistrationDatastore
	tokenDatastore        dao.TokenDatastore
	smtpDialer            *gomail.Dialer
}

func NewAccountService(
	datastore *dao.Datastore,
	registrationDatastore dao.RegistrationDatastore,
	tokenDatastore dao.TokenDatastore,
	smtpDialer *gomail.Dialer,
) accountpb.AccountServiceServer {
	return &accountService{
		datastore:             datastore,
		registrationDatastore: registrationDatastore,
		tokenDatastore:        tokenDatastore,
		smtpDialer:            smtpDialer,
	}
}

// Signup creates a new organization and admin, and triggers primary admin email verification
func (as *accountService) Signup(ctx context.Context, request *accountpb.SignupRequest) (*accountpb.SignupResponse, error) {
	if request.OrganizationName == "" {
		slog.Error("invalid organization name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid organization name")
	}
	if request.PrimaryAdministratorEmail == "" {
		slog.Error("invalid primary administrator email")
		return nil, status.Errorf(codes.InvalidArgument, "invalid primary administrator email")
	}
	if request.PrimaryAdministratorName == "" {
		slog.Error("invalid primary administrator name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid primary administrator name")
	}
	if request.PrimaryAdministratorCleartextPassword == "" {
		slog.Error("invalid primary administrator cleartext password")
		return nil, status.Errorf(codes.InvalidArgument, "invalid primary administrator cleartext password")
	}
	organization := &dao.Organization{
		Name:                      request.OrganizationName,
		PrimaryAdministratorEmail: request.PrimaryAdministratorEmail,
	}
	organizationID, err := as.datastore.Organizations.Create(organization)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create organization")
	}
	administrator := &dao.Administrator{
		Email:             request.PrimaryAdministratorEmail,
		DisplayName:       request.PrimaryAdministratorName,
		OrganizationID:    organizationID,
		AuthorizationRole: postgres.PrimaryAdmin,
	}
	administratorID, err := as.datastore.Administrators.Create(administrator, request.PrimaryAdministratorCleartextPassword)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create administrator")
	}
	token, emailCode, err := as.registrationDatastore.Create(&dao.Registration{
		OrganizationID:  organizationID,
		AdministratorID: administratorID,
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to create registration")
	}
	emailTemplateData := struct {
		Code string
	}{
		Code: emailCode,
	}
	tmpl, err := template.ParseFiles("templates/verify_email.html")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse registration email template %s", err.Error())
	}
	var body bytes.Buffer
	err = tmpl.Execute(&body, emailTemplateData)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to executing registration email template %s", err.Error())
	}

	m := gomail.NewMessage()
	m.SetHeader("From", emailFrom)
	m.SetHeader("To", request.PrimaryAdministratorEmail)
	m.SetHeader("Subject", "Cloud Inventory: Verify your email")
	m.SetBody("text/html", body.String())
	err = as.smtpDialer.DialAndSend(m)
	if err != nil {
		slog.Error(err.Error())
		return nil, status.Errorf(codes.Internal, "failed to send administrator email verification email %s", err.Error())
	}
	return &accountpb.SignupResponse{
		Token: token,
	}, nil
}

// Verify validates email verification codes, updates the administrators.verified column & creates admin UI & API JWTs
func (as *accountService) Verify(ctx context.Context, request *accountpb.VerificationRequest) (*accountpb.VerificationResponse, error) {
	registration, err := as.registrationDatastore.Read(request.Token)
	if err != nil {
		if err == redis.Nil {
			return nil, status.Errorf(codes.NotFound, "registration token not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read registration token: %s", err.Error())
	}
	if registration.EmailCode != request.VerificationCode {
		return nil, status.Errorf(codes.PermissionDenied, "verification code not authorized")
	}
	administrator, err := as.datastore.Administrators.Read(registration.AdministratorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "administrator not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read administrator data: %s", err.Error())
	}
	administrator.Verified = true
	err = as.datastore.Administrators.Update(administrator)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update administrator data: %s", err.Error())
	}
	err = as.registrationDatastore.Delete(request.Token)
	if err != nil {
		// non-fatal error, the registration data has a short TTL
		slog.Warn("failed to delete registration data")
	}
	adminUIAccessToken, err := as.tokenDatastore.Create(
		&dao.TokenData{
			OrganizationID:    administrator.OrganizationID,
			AdministratorID:   administrator.ID,
			AuthorizationRole: administrator.AuthorizationRole,
			TokenType:         ciJWT.Access,
			ClaimsType:        ciJWT.AdminUISession,
		},
		ciJWT.Access,
		ciJWT.AdminUISession,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate jwt: %s", err.Error())
	}
	adminUIRefreshToken, err := as.tokenDatastore.Create(
		&dao.TokenData{
			OrganizationID:    administrator.OrganizationID,
			AdministratorID:   administrator.ID,
			AuthorizationRole: administrator.AuthorizationRole,
			TokenType:         ciJWT.Refresh,
			ClaimsType:        ciJWT.AdminUISession,
		},
		ciJWT.Refresh,
		ciJWT.AdminUISession,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate jwt: %s", err.Error())
	}
	apiAccessToken, err := as.tokenDatastore.Create(
		&dao.TokenData{
			OrganizationID:    administrator.OrganizationID,
			AdministratorID:   administrator.ID,
			AuthorizationRole: administrator.AuthorizationRole,
			TokenType:         ciJWT.Access,
			ClaimsType:        ciJWT.APIAuthorization,
		},
		ciJWT.Access,
		ciJWT.APIAuthorization,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate jwt: %s", err.Error())
	}
	apiRefreshToken, err := as.tokenDatastore.Create(
		&dao.TokenData{
			OrganizationID:    administrator.OrganizationID,
			AdministratorID:   administrator.ID,
			AuthorizationRole: administrator.AuthorizationRole,
			TokenType:         ciJWT.Refresh,
			ClaimsType:        ciJWT.APIAuthorization,
		},
		ciJWT.Refresh,
		ciJWT.APIAuthorization,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate jwt: %s", err.Error())
	}
	return &accountpb.VerificationResponse{
		AdminUiAccessToken:  adminUIAccessToken,
		AdminUiRefreshToken: adminUIRefreshToken,
		ApiAccessToken:      apiAccessToken,
		ApiRefreshToken:     apiRefreshToken,
	}, nil
}
