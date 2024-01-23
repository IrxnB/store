package model

import "github.com/google/uuid"

type Seller struct {
	OwnerId uuid.UUID `json:"id"`
	Name    string    `json:"name"`
}
