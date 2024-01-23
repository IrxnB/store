package model

import (
	"math/big"

	"github.com/google/uuid"
)

type Product struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProductFull struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SellerName  string     `json:"seller_name"`
	Price       *big.Float `json:"price" swaggertype:"number"`
}
