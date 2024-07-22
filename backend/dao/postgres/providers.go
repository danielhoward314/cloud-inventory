package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/danielhoward314/cloud-inventory/backend/dao"
	"github.com/danielhoward314/cloud-inventory/backend/dao/postgres/queries"
)

type providers struct {
	db *sql.DB
}

// ProviderName is a type alias representing the postgres ENUM for the provider_name column
type ProviderName string

const (
	// AWS is a string for the AWS provider_name ENUM
	AWS = "AWS"
	// GCP is a string for the GCP provider_name ENUM
	GCP = "GCP"
	// AZURE is a string for the AZURE provider_name ENUM
	AZURE = "AZURE"
)

// NewProviders returns an instance implementing the Providers interface
func NewProviders(db *sql.DB) dao.Providers {
	return &providers{db: db}
}

func (p *providers) List(organizationID string) ([]*dao.Provider, error) {
	rows, err := p.db.Query(queries.ProvidersSelectByOrganizationID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var providers []*dao.Provider
	for rows.Next() {
		row := new(dao.Provider)
		err2 := rows.Scan(
			&row.ID,
			&row.ExternalIdentifier,
			&row.Name,
			&row.ProviderName,
			&row.Metadata,
			&row.OrganizationID,
		)
		if err2 != nil {
			return nil, err2
		}
		providers = append(providers, row)
	}
	return providers, nil
}
