package config

import (
	"database/sql"
	"errors"
	"os"
	"strconv"

	"github.com/danielhoward314/cloud-inventory/backend/dao"
	ciPostgres "github.com/danielhoward314/cloud-inventory/backend/dao/postgres"
	ciRedis "github.com/danielhoward314/cloud-inventory/backend/dao/redis"
	"github.com/go-redis/redis/v8"
)

type APIConfig struct {
	datastore             *dao.Datastore
	jwtSecret             string
	registrationDatastore dao.RegistrationDatastore
	smtpHost              string
	smtpPort              int
}

func NewAPIConfig(db *sql.DB, redisClient *redis.Client, jwtSecret string) (*APIConfig, error) {
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		return nil, errors.New("error: SMTP_HOST is empty")
	}
	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		return nil, errors.New("error: SMTP_PORT is empty")
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return nil, errors.New("error: invalid SMTP_PORT")
	}
	return &APIConfig{
		datastore:             ciPostgres.NewDatastore(db),
		jwtSecret:             jwtSecret,
		registrationDatastore: ciRedis.NewRegistrationDatastore(redisClient),
		smtpHost:              smtpHost,
		smtpPort:              smtpPort,
	}, nil
}

func (cfg *APIConfig) GetDatastore() *dao.Datastore {
	return cfg.datastore
}

func (cfg *APIConfig) GetJWTSecret() string {
	return cfg.jwtSecret
}

func (cfg *APIConfig) GetRegistrationDatastore() dao.RegistrationDatastore {
	return cfg.registrationDatastore
}

func (cfg *APIConfig) GetSMTPHost() string {
	return cfg.smtpHost
}

func (cfg *APIConfig) GetSMTPPort() int {
	return cfg.smtpPort
}
