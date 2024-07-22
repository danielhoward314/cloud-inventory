package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/danielhoward314/cloud-inventory/backend/dao"
	ciPostgres "github.com/danielhoward314/cloud-inventory/backend/dao/postgres"
	providerspb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/providers"
)

// providersService implements the providers gRPC service
type providersService struct {
	providerspb.UnimplementedProvidersServiceServer
	datastore *dao.Datastore
}

func NewProvidersService(
	datastore *dao.Datastore,
) providerspb.ProvidersServiceServer {
	return &providersService{
		datastore: datastore,
	}
}

func (ps *providersService) List(ctx context.Context, request *providerspb.ListRequest) (*providerspb.ListResponse, error) {
	if request.OrganizationId == "" {
		slog.Error("invalid organization id")
		return nil, status.Errorf(codes.InvalidArgument, "invalid organization id")
	}
	slog.Info("getting list of providers", "request.organization_id", request.OrganizationId)
	providers, err := ps.datastore.Providers.List(request.OrganizationId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "no providers found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to read provider data: %s", err.Error())
	}
	pbProviders := make([]*providerspb.Provider, 0, len(providers))
	for _, provider := range providers {
		pbProvider := &providerspb.Provider{
			Id:                 provider.ID,
			ExternalIdentifier: provider.ExternalIdentifier,
			Name:               provider.Name,
			ProviderName:       provider.ProviderName,
			OrganizationId:     provider.OrganizationID,
		}
		var err2 error
		slog.Info("unmarshaling provider.metadata json", "provider_name", provider.ProviderName)
		switch provider.ProviderName {
		case ciPostgres.AWS:
			var awsMetadata providerspb.Provider_AwsMetadata
			err2 = json.Unmarshal(provider.Metadata, &awsMetadata)
			if err2 != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal provider.metadata into AWS metadata: %s", err2.Error())
			}
			pbProvider.Metadata = &awsMetadata
		case ciPostgres.GCP:
			var gcpMetadata providerspb.Provider_GcpMetadata
			err2 = json.Unmarshal(provider.Metadata, &gcpMetadata)
			if err2 != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal provider.metadata into GCP metadata: %s", err2.Error())
			}
			pbProvider.Metadata = &gcpMetadata
		case ciPostgres.AZURE:
			var azureMetadata providerspb.Provider_AzureMetadata
			err2 = json.Unmarshal(provider.Metadata, &azureMetadata)
			if err2 != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal provider.metadata into Azure metadata: %s", err2.Error())
			}
			pbProvider.Metadata = &azureMetadata
		default:
			slog.Error("invalid provider_name")
			return nil, status.Errorf(codes.Internal, "invalid provider_name")
		}
		pbProviders = append(pbProviders, pbProvider)
	}
	return &providerspb.ListResponse{
		Providers: pbProviders,
	}, nil
}
