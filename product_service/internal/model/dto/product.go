package dto

import (
	"math/big"
	"product_service/internal/model"

	"github.com/google/uuid"
)

type ProductPage struct {
	Products []model.Product `json:"products"`
	Page     int             `json:"page"`
	Limit    int             `json:"limit"`
}

type Products struct {
	Products []model.Product `json:"products"`
}

type CreateProduct struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       *big.Float  `json:"price" swaggertype:"number"`
	CategoryIds []uuid.UUID `json:"categories"`
}
type UpdateProduct struct {
	Id          uuid.UUID   `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       *big.Float  `json:"price"`
	CategoryIds []uuid.UUID `json:"categories"`
}

type CreateProductRequest struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       *big.Float  `json:"price"`
	CategoryIds []uuid.UUID `json:"categories"`
}

type GetBatchItem struct {
	Id          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       *big.Float `json:"price"`
	SellerName  string     `json:""`
}

type GetBatchResponse struct {
	Products []GetBatchItem `json:"products"`
}

type GetBatchRequest struct {
	ProductIds []uuid.UUID `json:"product_ids"`
}
