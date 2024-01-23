package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"product_service/internal/jwt"
)

type JwkWebApi interface {
	GetJWK() (*jwt.JWK, error)
}
type JwkWebApiImpl struct {
	url string
}

func NewJwkWebApi(baseUrl string) (*JwkWebApiImpl, error) {
	_, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	return &JwkWebApiImpl{url: baseUrl}, nil
}

func (jwks *JwkWebApiImpl) GetJWK() (*jwt.JWK, error) {
	request, err := url.JoinPath(jwks.url, "/public-key")
	if err != nil {
		return nil, fmt.Errorf("intertal")
	}
	resp, err := http.Get(request)
	if err != nil {
		return nil, fmt.Errorf("intertal")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("intertal")
	}

	var jwk jwt.JWK

	err = json.Unmarshal(body, &jwk)
	if err != nil {
		return nil, fmt.Errorf("intertal")
	}

	return &jwk, nil
}
