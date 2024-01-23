package usecase

import (
	"context"
	"store/cart_service/internal/model"
	"store/cart_service/internal/model/dto"
	"store/cart_service/internal/oauth"

	"github.com/google/uuid"
)

type Cart interface {
	AddOrUpdate(ctx context.Context, entries []dto.AddToCartRequest, user oauth.OauthUser) (err error)
	GetCart(ctx context.Context, userId uuid.UUID) (entries []model.CartEntry, err error)
	Remove(ctx context.Context, product_ids []uuid.UUID, user oauth.OauthUser) (err error)
	GetProducts(ctx context.Context, user oauth.OauthUser) (products []model.ProductFull, err error)
}
