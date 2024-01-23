package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"oauth2_provider/internal/model"
	"oauth2_provider/internal/service"

	"github.com/google/uuid"
)

type ClientServer struct {
	ClientService service.ClientService
}

func (cs *ClientServer) SaveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "supports post", http.StatusMethodNotAllowed)
		return
	}

	var createClient model.CreateClientRequest

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(data, &createClient)
	if err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	appType, ok := model.ValidAppType[createClient.Type]
	if !ok {
		http.Error(w, "app type not supported", http.StatusBadRequest)
		return
	}

	client := model.Client{Id: uuid.New(), Secret: uuid.New().String(), Domain: createClient.Domain, Type: appType}

	if _, err := url.Parse(client.Domain); err != nil {
		http.Error(w, "invalid domain", http.StatusBadRequest)
	}

	err = cs.ClientService.Save(r.Context(), &client)

	if err != nil {
		http.Error(w, "error creating client", http.StatusBadRequest)
		return
	}

	json, _ := json.Marshal(model.CreateClientResponse{ClientId: client.Id, ClientSecret: client.Secret})

	w.Write(json)

}
