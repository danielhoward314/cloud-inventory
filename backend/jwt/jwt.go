package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	authpb "github.com/danielhoward314/cloud-inventory/backend/protogen/golang/auth"
)

// ClaimsType is an enum representing different claims to be used in JWT generation
type ClaimsType int

const (
	AdminUISession ClaimsType = iota
	APIAuthorization
)

// TokenType is an enum representing different token types
type TokenType int

const (
	Access TokenType = iota
	Refresh
)

const (
	OrganizationIDKey     = "organization_id"
	AdministratorIDKey    = "administrator_id"
	AuthorizationRoleKey  = "authorization_role"
	TokenTypeKey          = "token_type"
	ClaimsTypeKey         = "claims_type"
	issuerClaimValue      = "cloud-inventory-api"
	uiAudienceClaimValue  = "cloud-inventory-ui"
	apiAudienceClaimValue = "cloud-inventory-api"
	TokenExpiredError     = "token expired"
	InvalidTokenError     = "invalid token"
)

const (
	adminUIAccessTokenExpiry           = 3 * time.Hour
	adminUIRefreshTokenExpiry          = 7 * 24 * time.Hour
	apiAuthorizationAccessTokenExpiry  = 15 * time.Minute
	apiAuthorizationRefreshTokenExpiry = 7 * 24 * time.Hour
)

type AdminUISessionClaims struct {
	OrganizationID    string `json:"organization_id"`
	AuthorizationRole string `json:"authorization_role"`
	jwt.StandardClaims
}

type APIAuthorizationClaims struct {
	OrganizationID    string `json:"organization_id"`
	AuthorizationRole string `json:"authorization_role"`
	jwt.StandardClaims
}

func GenerateJWT(secret string, tokenType TokenType, claimsType ClaimsType, claimsData map[string]interface{}) (string, error) {
	jwtID := uuid.NewString()
	// fields common across all claims types
	standardClaims := jwt.StandardClaims{
		Issuer:   issuerClaimValue,
		IssuedAt: time.Now().Unix(),
		Id:       jwtID,
	}
	switch claimsType {
	case AdminUISession:
		if tokenType == Access {
			standardClaims.ExpiresAt = time.Now().Add(adminUIAccessTokenExpiry).Unix()
		} else {
			standardClaims.ExpiresAt = time.Now().Add(adminUIRefreshTokenExpiry).Unix()
		}
		organizationID, ok := claimsData[OrganizationIDKey]
		if !ok {
			return "", errors.New("missing organization_id claims")
		}
		administratorID, ok := claimsData[AdministratorIDKey]
		if !ok {
			return "", errors.New("missing administrator_id claims")
		}
		authorizationRole, ok := claimsData[AuthorizationRoleKey]
		if !ok {
			return "", errors.New("missing authorization_role claims")
		}
		standardClaims.Subject = administratorID.(string)
		standardClaims.Audience = uiAudienceClaimValue
		claims := &AdminUISessionClaims{
			OrganizationID:    organizationID.(string),
			AuthorizationRole: authorizationRole.(string),
			StandardClaims:    standardClaims,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(secret))
		if err != nil {
			return "", err
		}
		return tokenString, nil
	case APIAuthorization:
		if tokenType == Access {
			standardClaims.ExpiresAt = time.Now().Add(apiAuthorizationAccessTokenExpiry).Unix()
		} else {
			standardClaims.ExpiresAt = time.Now().Add(apiAuthorizationRefreshTokenExpiry).Unix()
		}
		organizationID, ok := claimsData[OrganizationIDKey]
		if !ok {
			return "", errors.New("missing organization_id claims")
		}
		administratorID, ok := claimsData[AdministratorIDKey]
		if !ok {
			return "", errors.New("missing administrator_id claims")
		}
		authorizationRole, ok := claimsData[AuthorizationRoleKey]
		if !ok {
			return "", errors.New("missing authorization_role claims")
		}
		standardClaims.Subject = administratorID.(string)
		standardClaims.Audience = apiAudienceClaimValue
		claims := &APIAuthorizationClaims{
			OrganizationID:    organizationID.(string),
			AuthorizationRole: authorizationRole.(string),
			StandardClaims:    standardClaims,
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
	case AdminUISession:
		claims := &AdminUISessionClaims{}
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
			return errors.New(InvalidTokenError)
		}
		if claims.ExpiresAt < time.Now().Unix() {
			return errors.New(TokenExpiredError)
		}
		return nil
	}
	return errors.New("unknown claims type")
}

func GetClaimsTypeFromProtoEnum(ct authpb.ClaimsType) (ClaimsType, error) {
	switch ct {
	case authpb.ClaimsType_ADMIN_UI_SESSION:
		return AdminUISession, nil
	case authpb.ClaimsType_API_AUTHORIZATION:
		return APIAuthorization, nil
	default:
		return AdminUISession, errors.New("invalid claims type")
	}
}
