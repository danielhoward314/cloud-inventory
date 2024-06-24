package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// ClaimsType is an enum representing different claims to be used in JWT generation
type ClaimsType int

const (
	Session ClaimsType = iota
)

const (
	// OrganizationIDKey is the map key for claims data that requires an organization ID
	OrganizationIDKey = "organization_id"
	// AdministratorIDKey is the map key for claims data that requires an administrator ID
	AdministratorIDKey = "administrator_id"
	issuerClaimValue   = "cloud-inventory-api"
	audienceClaimValue = "cloud-inventory-ui"
)

type SessionClaims struct {
	OrganizationID string `json:"organization_id"`
	jwt.StandardClaims
}

func GenerateJWT(secret string, claimsType ClaimsType, claimsData map[string]interface{}) (string, error) {
	jwtID := uuid.NewString()
	switch claimsType {
	case Session:
		expirationTime := time.Now().Add(168 * time.Hour)
		organizationID, ok := claimsData[OrganizationIDKey]
		if !ok {
			return "", errors.New("missing organization_id claims")
		}
		administratorID, ok := claimsData[AdministratorIDKey]
		if !ok {
			return "", errors.New("missing administrator_id claims")
		}
		claims := &SessionClaims{
			OrganizationID: organizationID.(string),
			StandardClaims: jwt.StandardClaims{
				Issuer:    issuerClaimValue,
				Subject:   administratorID.(string),
				Audience:  audienceClaimValue,
				ExpiresAt: expirationTime.Unix(),
				IssuedAt:  time.Now().Unix(),
				Id:        jwtID,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			return "", err
		}
		return tokenString, nil
	}
	return "", errors.New("unknown claims type")
}

func DecodeJWT(secret string, tokenString string, claimsType ClaimsType) error {
	switch claimsType {
	case Session:
		claims := &SessionClaims{}
		parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil {
			return errors.New("failed to parse token")
		}
		if !parsedToken.Valid {
			return errors.New("invalid token")
		}
		if claims.ExpiresAt < time.Now().Unix() {
			return errors.New("session token expired")
		}
		return nil
	}
	return errors.New("unknown claims type")
}
