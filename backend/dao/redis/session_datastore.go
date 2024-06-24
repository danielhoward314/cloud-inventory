package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/danielhoward314/cloud-inventory/backend/dao"
	ciJWT "github.com/danielhoward314/cloud-inventory/backend/jwt"
)

const (
	// sessionTokenTTL = time.Hour * 168
	// TODO: change back to 1 week after jwt validation logic is complete
	sessionTokenTTL = time.Minute * 3
)

type sessionDatastore struct {
	sessionJWTSecret string
	client           *redis.Client
}

// NewSessionDatastore returns a redis implementation for the session key-value datastore
func NewSessionDatastore(client *redis.Client, sessionJWTSecret string) dao.SessionDatastore {
	return &sessionDatastore{
		client:           client,
		sessionJWTSecret: sessionJWTSecret,
	}
}

func (sds *sessionDatastore) Create(session *dao.Session) (string, error) {
	if session.OrganizationID == "" {
		return "", errors.New("invalid organization id")
	}
	if session.AdministratorID == "" {
		return "", errors.New("invalid administrator id")
	}
	claimsData := make(map[string]interface{})
	claimsData[ciJWT.OrganizationIDKey] = session.OrganizationID
	claimsData[ciJWT.AdministratorIDKey] = session.AdministratorID
	jwt, err := ciJWT.GenerateJWT(sds.sessionJWTSecret, ciJWT.Session, claimsData)
	if err != nil {
		return "", err
	}
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	status := sds.client.Set(context.Background(), jwt, sessionJSON, sessionTokenTTL)
	if status.Err() != nil {
		return "", err
	}
	return jwt, nil
}

func (sds *sessionDatastore) Read(jwt string) (*dao.Session, error) {
	sessionJSON, err := sds.client.Get(context.Background(), jwt).Result()
	if err != nil {
		return nil, err
	}
	var session dao.Session
	err = json.Unmarshal([]byte(sessionJSON), &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (sds *sessionDatastore) Delete(jwt string) error {
	cmdStatus := sds.client.Del(context.Background(), jwt)
	if cmdStatus.Err() != nil {
		return cmdStatus.Err()
	}
	return nil
}
