package model

import "github.com/google/uuid"

type Role struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
