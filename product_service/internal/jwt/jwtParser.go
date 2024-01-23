package jwt

import (
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	ClientId  string = "client_id"
	UserId    string = "user_id"
	UserRoles string = "roles"
	TokenId   string = "token_id"
	CreatedAt string = "created_at"
	ExpiresAt string = "expires_at"
	Scope     string = "scope"
	TokenType string = "type"
)

type PublicKeyProvider interface {
	GetPublic() (*rsa.PublicKey, error)
}

type JwtParser struct {
	kp PublicKeyProvider
}

func NewJwtParser(kp PublicKeyProvider) *JwtParser {
	return &JwtParser{kp: kp}
}

func (jp JwtParser) ParseToken(tokenStr string) (*Token, error) {
	token := Token{}
	var err error

	tok, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("validation error")
		}
		return jp.kp.GetPublic()
	})

	tokenClaims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}

	if err != nil {
		return nil, err
	}

	clientId, ok := tokenClaims[ClientId]
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}
	userId, ok := tokenClaims[UserId]
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}
	expiresAt, ok := tokenClaims[ExpiresAt]
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}
	scope, ok := tokenClaims[Scope]
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}

	clientIdStr, ok := clientId.(string)
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}
	userIdStr, ok := userId.(string)
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}
	token.ClientId, err = uuid.Parse(clientIdStr)
	if err != nil {
		return nil, fmt.Errorf("wrong token")
	}
	token.UserId, err = uuid.Parse(userIdStr)
	if err != nil {
		return nil, fmt.Errorf("wrong token")
	}
	expiresAtFloat, ok := expiresAt.(float64)
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}
	token.ExpiresAt = int64(expiresAtFloat)
	token.Scope, ok = scope.(string)
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}
	roles, ok := tokenClaims[UserRoles]
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}

	rolesStr, ok := roles.(string)
	if !ok {
		return nil, fmt.Errorf("wrong token")
	}

	tokenId, ok := tokenClaims[TokenId]
	if ok {
		idStr, ok := tokenId.(string)
		if !ok {
			return nil, fmt.Errorf("wrong token")
		}
		token.TokenId, err = uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("wrong token")
		}
	}

	token.UserRoleNames = strings.Split(rolesStr, " ")

	return &token, nil
}
