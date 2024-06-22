package postgres

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"

	"github.com/danielhoward314/cloud-inventory/backend/dao"
	"github.com/danielhoward314/cloud-inventory/backend/dao/postgres/queries"
)

type administrators struct {
	db *sql.DB
}

// NewAdministrators returns an instance implementing the Administrators interface
func NewAdministrators(db *sql.DB) dao.Administrators {
	return &administrators{db: db}
}

func (o *administrators) Create(administrator *dao.Administrator) (string, error) {
	if administrator.Email == "" {
		return "", errors.New("invalid administrator email")
	}
	if administrator.DisplayName == "" {
		return "", errors.New("invalid administrator display name")
	}
	if administrator.OrganizationID == "" {
		return "", errors.New("invalid administrator organization id")
	}
	var id string
	err := o.db.QueryRow(queries.AdministratorsInsert, administrator.Email, administrator.DisplayName, administrator.OrganizationID).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

// func (o *administrators) Read(id string) (*dao.Administrator, error) {
// 	return nil, nil
// }

// func (o *administrators) Update(*dao.Administrator) (*dao.Administrator, error) {
// 	return nil, nil
// }

// func (o *administrators) Delete(id string) (*dao.Administrator, error) {
// 	return nil, nil
// }
