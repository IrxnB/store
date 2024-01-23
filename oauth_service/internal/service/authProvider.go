package service

import (
	"context"
	"fmt"
	"oauth2_provider/internal/model"
	"time"

	"github.com/google/uuid"
)

type AuthProvider interface {
	Authorize(ctx context.Context,
		clientId uuid.UUID,
		credentials model.Credentials,
		state string,
		responseType string,
		redirectBase string,
		scope string) (redirectUrl string, err error)
	Token(ctx context.Context,
		clientId uuid.UUID,
		clientSecret string,
		authCode string) (token string, err error)
}

type AuthProviderImpl struct {
	ClientService  ClientService
	UserService    UserService
	JwtProvider    JwtProvider
	SessionService SessionService
}

func (ap *AuthProviderImpl) Authorize(ctx context.Context,
	clientId uuid.UUID,
	credentials model.Credentials,
	state string,
	responseType string,
	redirectBase string,
	scope string) (redirectUrl string, err error) {
	rt := model.ResponseType(responseType)

	client, err := ap.ClientService.GetById(ctx, clientId)
	if err != nil {
		return "", err
	}

	if _, ok := model.ValidAppType[string(client.Type)]; client.Type == model.ClientSide && rt == model.Code || !ok {
		return "", err
	}

	user, err := ap.UserService.GetWithValidation(ctx, credentials)
	if err != nil {
		return "", err
	}

	roleNames := make([]string, 0, len(user.Roles))
	for _, r := range user.Roles {
		roleNames = append(roleNames, r.Name)
	}

	token, err := ap.JwtProvider.GenerateToken(user.Id, client.Id, scope, rt, roleNames)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s?%s=%s&state=%s", redirectBase, responseType, token, state), nil
}

func (ap *AuthProviderImpl) Token(ctx context.Context,
	clientId uuid.UUID,
	clientSecret string,
	authCode string) (accessToken string, err error) {

	client, err := ap.ClientService.GetWithValidation(ctx, clientId, clientSecret)

	if err != nil {
		return "", err
	}

	codeClaims, err := ap.JwtProvider.ParseToken(authCode)

	if codeClaims.TokenId == uuid.Nil {
		return "", fmt.Errorf("forbidden")
	}

	if err != nil {
		return "", err
	}

	if codeClaims.ClientId != client.Id {
		return "", err
	}

	if expAt := time.Unix(codeClaims.ExpiresAt, 0); expAt.Before(time.Now()) {
		return "", fmt.Errorf("expired")
	}

	token, err := ap.JwtProvider.GenerateAccessFromAuthCode(codeClaims)
	if err != nil {
		return "", err
	}

	if expAt := time.Unix(codeClaims.ExpiresAt, 0); expAt.Before(time.Now()) {
		return "", fmt.Errorf("internal")
	}

	err = ap.SessionService.Save(ctx, codeClaims.TokenId, codeClaims.UserId, codeClaims.ClientId)
	if err != nil {
		return "", fmt.Errorf("already used")
	}

	return token, nil

}
