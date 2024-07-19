package dao

import (
	ciJWT "github.com/danielhoward314/cloud-inventory/backend/jwt"
)

type TokenData struct {
	OrganizationID    string           `json:"organization_id"`
	AdministratorID   string           `json:"administrator_id"`
	AuthorizationRole string           `json:"authorization_role"`
	TokenType         ciJWT.TokenType  `json:"token_type"`
	ClaimsType        ciJWT.ClaimsType `json:"claims_type"`
}

// TokenDatastore defines the interface for access token operations in a key-value datastore
type TokenDatastore interface {
	Create(tokenData *TokenData, tokenType ciJWT.TokenType, claimsType ciJWT.ClaimsType) (string, error)
	Read(token string) (*TokenData, error)
	Decode(tokenType ciJWT.TokenType, tokenString string, claimsType ciJWT.ClaimsType) error
	Delete(token string) error
}
