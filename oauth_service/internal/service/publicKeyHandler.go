package service

import (
	"encoding/json"
	"net/http"
	"oauth2_provider/internal/jwt"
)

func GetPublicKeyHandler(kp jwt.PublicKeyProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := kp.GetPublic()
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		encoder := json.NewEncoder(w)

		encoder.Encode(jwt.JWK{N: key.N, E: key.E})
	}
}
