package usecase

import (
	"context"
	"product_service/internal/model"
	"product_service/internal/model/dto"
	"product_service/internal/oauth"

	"github.com/google/uuid"
)

type Product interface {
	Create(ctx context.Context, createProduct dto.CreateProduct, sellerId uuid.UUID) (err error)
	GetById(ctx context.Context, id uuid.UUID) (product *model.Product, err error)
	GetBacth(ctx context.Context, ids []uuid.UUID) (product []dto.GetBatchItem, err error)
	Update(ctx context.Context, id uuid.UUID, product dto.CreateProduct, user oauth.OauthUser) (err error)
	GetPage(ctx context.Context, page int, limit int) (products []model.Product, err error)
	GetPageBySellerId(ctx context.Context, page int, limit int, id uuid.UUID) (products []model.Product, err error)
}

type Seller interface {
	GetAll(ctx context.Context) (sellers []model.Seller, err error)
	Create(ctx context.Context, createSellerDto dto.CreateSeller, user oauth.OauthUser) (err error)
	Update(ctx context.Context, sellerId uuid.UUID, updateSellerDto dto.UpdateSeller, user oauth.OauthUser) (err error)
	GetById(ctx context.Context, id uuid.UUID) (seller *model.Seller, err error)
	GetProductPage(ctx context.Context, id uuid.UUID, page int, limit int) (products []model.Product, err error)
	AddProduct(ctx context.Context, createProduct dto.CreateProduct, sellerId uuid.UUID) (err error)
}
