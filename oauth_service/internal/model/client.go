package model

import "github.com/google/uuid"

type AppType string

const (
	ServerSide AppType = "server"
	ClientSide AppType = "client"
)

var (
	ValidAppType = map[string]AppType{
		string(ServerSide): ServerSide,
		string(ClientSide): ClientSide,
	}
)

type Client struct {
	Id     uuid.UUID
	Secret string
	Domain string
	Type   AppType
}

type CreateClientRequest struct {
	Domain string `json:"domain"`
	Type   string `json:"type"`
}

type CreateClientResponse struct {
	ClientId     uuid.UUID `json:"id"`
	ClientSecret string    `json:"secret"`
}
