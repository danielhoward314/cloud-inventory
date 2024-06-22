package services

import (
	"bytes"
	"context"
	"crypto/tls"
	"html/template"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/gomail.v2"

	"github.com/danielhoward314/cloud-inventory/backend/cmd/config"
	"github.com/danielhoward314/cloud-inventory/backend/dao"
	accountpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/account"
)

const (
	emailFrom = "no.reply@CloudInventory.com"
)

func NewAccountService(cfg *config.APIConfig) accountpb.AccountServiceServer {
	smtpDialer := gomail.NewDialer(cfg.GetSMTPHost(), cfg.GetSMTPPort(), "", "")
	smtpDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return &accountService{
		datastore:             cfg.GetDatastore(),
		registrationDatastore: cfg.GetRegistrationDatastore(),
		smtpDialer:            smtpDialer,
	}
}

// accountService implements the account gRPC service
type accountService struct {
	accountpb.UnimplementedAccountServiceServer
	smtpDialer            *gomail.Dialer
	registrationDatastore dao.RegistrationDatastore
	datastore             *dao.Datastore
}

// Signup creates a new organization and admin, and triggers primary admin email verification
func (as *accountService) Signup(ctx context.Context, request *accountpb.SignupRequest) (*accountpb.SignupResponse, error) {
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
	administratorID, err := as.datastore.Administrators.Create(administrator)
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

	// Create a new email message.
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
