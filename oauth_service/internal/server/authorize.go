package server

import (
	"encoding/json"
	"io"
	"net/http"
	"oauth2_provider/internal/model"
	"oauth2_provider/internal/service"

	"github.com/google/uuid"
)

const (
	clientId      = "client_id"
	clientSecret  = "client_secret"
	response_type = "response_type"
	state         = "state"
	scope         = "scope"
	username      = "username"
	password      = "password"
	redirectUrl   = "redirect_url"
)

type OauthServer interface {
	HandleAuthorize(w http.ResponseWriter, r *http.Request)
	HandleToken(w http.ResponseWriter, r *http.Request)
	HandleValidate(w http.ResponseWriter, r *http.Request)
}

type OauthServerImpl struct {
	Ap service.AuthProvider
}

func (s *OauthServerImpl) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	clientId, err := uuid.Parse(q.Get(clientId))
	if err != nil {
		http.Error(w, "not found client", http.StatusBadRequest)
	}

	responseType := q.Get(response_type)
	redirectBase := q.Get(redirectUrl)
	scope := q.Get(scope)
	state := q.Get(state)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		return
	}

	var creds model.Credentials

	err = json.Unmarshal(body, &creds)
	if err != nil {
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		return
	}

	redirectUrl, err := s.Ap.Authorize(r.Context(), clientId, creds, state, responseType, redirectBase, scope)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	redirect(w, r, redirectUrl)
}

func (s *OauthServerImpl) HandleToken(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	clientIdStr := q.Get("client_id")
	clientSecret := q.Get("client_secret")
	code := q.Get("code")

	clientId, err := uuid.Parse(clientIdStr)
	if err != nil {
		http.Error(w, "not found client", http.StatusBadRequest)
		return
	}

	token, err := s.Ap.Token(r.Context(), clientId, clientSecret, code)
	if err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	json, err := json.Marshal(token)
	if err != nil {
		http.Error(w, "json marshaling error", http.StatusInternalServerError)
		return
	}

	w.Write(json)
}
func (s *OauthServerImpl) HandleValidate(w http.ResponseWriter, r *http.Request) {}

func redirect(w http.ResponseWriter, r *http.Request, url string) {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
