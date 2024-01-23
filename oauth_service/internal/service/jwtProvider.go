package service

import (
	"log"
	"oauth2_provider/internal/encoding"
	oauth_jwt "oauth2_provider/internal/jwt"
	"oauth2_provider/internal/model"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	tokenExpiration = time.Minute * 3600
)

type JwtProvider interface {
	GenerateToken(userId, clientId uuid.UUID, scope string, tokenType model.ResponseType, roleNames []string) (string, error)
	ParseToken(tokenStr string) (*oauth_jwt.Token, error)
	GenerateAccessFromAuthCode(authCode *oauth_jwt.Token) (token string, err error)
}

type JwtProviderImpl struct {
	keyStorage *encoding.KeyFileStorage
	jwtParser  *oauth_jwt.JwtParser
}

func NewJwtProvider(keyStorage encoding.KeyFileStorage) *JwtProviderImpl {
	return &JwtProviderImpl{keyStorage: &keyStorage, jwtParser: oauth_jwt.NewJwtParser(&keyStorage)}
}

type TokenClaims jwt.MapClaims

func (jp *JwtProviderImpl) GenerateToken(userId, clientId uuid.UUID, scope string, tokenType model.ResponseType, roleNames []string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)
	claims[oauth_jwt.UserId] = userId
	claims["client_id"] = clientId

	now := time.Now()
	claims[oauth_jwt.ExpiresAt] = now.Add(tokenExpiration).Unix()
	claims[oauth_jwt.TokenType] = tokenType
	claims[oauth_jwt.UserRoles] = strings.Join(roleNames, " ")
	claims[oauth_jwt.Scope] = scope
	if tokenType == model.Code {
		claims[oauth_jwt.TokenId] = uuid.New()
	}

	private, err := jp.keyStorage.GetPrivate()
	if err != nil {
		log.Fatal(err)
	}

	tokenStr, err := token.SignedString(private)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (jp *JwtProviderImpl) ParseToken(tokenStr string) (*oauth_jwt.Token, error) {
	return jp.jwtParser.ParseToken(tokenStr)
}

func (jp *JwtProviderImpl) GenerateAccessFromAuthCode(authCode *oauth_jwt.Token) (token string, err error) {
	token, err = jp.GenerateToken(authCode.UserId, authCode.ClientId, authCode.Scope, model.Token, authCode.UserRoleNames)
	return
}
