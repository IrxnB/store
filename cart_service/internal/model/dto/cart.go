package dto

import (
	"store/cart_service/internal/model"

	"github.com/google/uuid"
)

type AddToCartRequest struct {
	ProductId uuid.UUID `json:"product_id"`
	Ammount   int       `json:"ammount"`
}

type AddToCart struct {
	Requests []AddToCartRequest `json:"batch"`
}

type GetBatchRequest struct {
	ProductIds []uuid.UUID `json:"product_ids"`
}

type CheckExistResponse struct {
	Products []model.Product `json:"products"`
}

type GetBatchResponse struct {
	Products []model.ProductFull `json:"products"`
}

type RemoveFromCart struct {
	ProductIds []uuid.UUID `json:"product_ids"`
}
