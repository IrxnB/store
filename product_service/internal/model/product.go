package model

import (
	"math/big"

	"github.com/google/uuid"
)

type Product struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       *big.Float `json:"price" swaggertype:"number"`
	Categories  []Category `json:"categories"`
	Seller      Seller     `json:"seller"`
}

type Category struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
