package services

import (
	"bytes"
	"context"
	"crypto/tls"
	"database/sql"
	"html/template"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/gomail.v2"

	"github.com/danielhoward314/cloud-inventory/backend/cmd/config"
	"github.com/danielhoward314/cloud-inventory/backend/dao"
	accountpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/account"
	"github.com/go-redis/redis/v8"
)

const (
	emailFrom = "no.reply@CloudInventory.com"
)

// accountService implements the account gRPC service
type accountService struct {
	accountpb.UnimplementedAccountServiceServer
	datastore             *dao.Datastore
	registrationDatastore dao.RegistrationDatastore
	sessionDatastore      dao.SessionDatastore
	smtpDialer            *gomail.Dialer
}

func NewAccountService(cfg *config.APIConfig) accountpb.AccountServiceServer {
	smtpDialer := gomail.NewDialer(cfg.GetSMTPHost(), cfg.GetSMTPPort(), "", "")
	smtpDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &accountService{
		datastore:             cfg.GetDatastore(),
		registrationDatastore: cfg.GetRegistrationDatastore(),
		sessionDatastore:      cfg.GetSessionDatastore(),
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
		Email:          request.PrimaryAdministratorEmail,
		DisplayName:    request.PrimaryAdministratorName,
		OrganizationID: organizationID,
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

// Verify validates email verification codes, updates the administrators.verified column & creates a user session
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
	jwt, err := as.sessionDatastore.Create(&dao.Session{
		OrganizationID:  administrator.OrganizationID,
		AdministratorID: administrator.ID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate jwt: %s", err.Error())
	}
	return &accountpb.VerificationResponse{
		Jwt: jwt,
	}, nil
}
