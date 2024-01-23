package usecase

import (
	"context"
	"fmt"
	"product_service/internal/model"
	"product_service/internal/model/dto"
	"product_service/internal/usecase/repo"

	"product_service/internal/oauth"

	"github.com/google/uuid"
)

type ProductImpl struct {
	repo *repo.Product
}

func NewProduct(repo *repo.Product) (productImpl *ProductImpl) {
	return &ProductImpl{repo: repo}
}

func (p *ProductImpl) Create(ctx context.Context, createProduct dto.CreateProduct, sellerId uuid.UUID) (err error) {
	return p.repo.Save(ctx, createProduct, sellerId)
}

func (p *ProductImpl) GetById(ctx context.Context, id uuid.UUID) (product *model.Product, err error) {
	return p.repo.GetById(ctx, id)
}
func (p *ProductImpl) Update(ctx context.Context, id uuid.UUID, product dto.CreateProduct, user oauth.OauthUser) (err error) {
	if !isSeller(user) {
		return fmt.Errorf("forbidden")
	}

	fromDb, err := p.repo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("not found")
	}

	if fromDb.Id != user.Id {
		return fmt.Errorf("forbidden")
	}

	fromDb.Name = product.Name
	fromDb.Description = product.Description
	fromDb.Price = product.Price
	fromDb.Categories = make([]model.Category, len(product.CategoryIds))

	for _, cid := range product.CategoryIds {
		fromDb.Categories = append(fromDb.Categories, model.Category{Id: cid})
	}

	return p.repo.Update(ctx, *fromDb)
}
func (p *ProductImpl) GetPage(ctx context.Context, page int, limit int) (products []model.Product, err error) {
	return p.repo.GetPage(ctx, page, limit)
}

func (p *ProductImpl) GetPageBySellerId(ctx context.Context, page int, limit int, id uuid.UUID) (products []model.Product, err error) {
	return p.repo.GetPageBySellerId(ctx, page, limit, id)
}

func (p *ProductImpl) GetBacth(ctx context.Context, ids []uuid.UUID) (product []dto.GetBatchItem, err error) {
	return p.repo.GetBatch(ctx, ids)
}

func isSeller(user oauth.OauthUser) bool {
	return user.HasRole("seller")
}
