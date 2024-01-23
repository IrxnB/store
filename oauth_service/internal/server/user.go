package server

import (
	"encoding/json"
	"io"
	"net/http"
	"oauth2_provider/internal/model"
	"oauth2_provider/internal/service"

	"github.com/google/uuid"
)

type UserServer struct {
	UserService service.UserService
}

func (us *UserServer) SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "supports post", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}

	var creds model.Credentials

	err = json.Unmarshal(body, &creds)

	if err != nil {
		http.Error(w, "wrong body", http.StatusBadRequest)
		return
	}

	user := model.User{Id: uuid.New(), Username: creds.Username, PasswordHash: creds.Password, Roles: make([]model.Role, 1)}

	user.Roles = append(user.Roles, model.Role{Name: "user"})

	err = us.UserService.Save(r.Context(), &user)

	if err != nil {
		http.Error(w, "error creating client", http.StatusBadRequest)
		return
	}

	json, _ := json.Marshal(user.Id)

	w.Write(json)
}
