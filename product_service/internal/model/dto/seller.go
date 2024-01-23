package dto

import "product_service/internal/model"

type CreateSeller struct {
	Name string `json:"name"`
}

type UpdateSeller struct {
	Name string `json:"name"`
}

type SellerList struct {
	Sellers []model.Seller `json:"sellers"`
}
