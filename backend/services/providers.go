package services

import (
	"context"

	"github.com/danielhoward314/cloud-inventory/backend/dao"
	providerspb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/providers"
)

const (
	accessTokenKey = "access_token"
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
	return &providerspb.ListResponse{Test: "test response"}, nil
}
