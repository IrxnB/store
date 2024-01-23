package model

import "github.com/google/uuid"

type CartEntry struct {
	ProductId uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Ammount   int       `json:"ammount"`
}
