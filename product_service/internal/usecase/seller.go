package usecase

import (
	"context"
	"fmt"
	"product_service/internal/model"
	"product_service/internal/model/dto"
	"product_service/internal/oauth"
	"product_service/internal/usecase/repo"

	"github.com/google/uuid"
)

type SellerImpl struct {
	repo *repo.Seller
	ps   Product
}

func NewSeller(repo *repo.Seller, prod Product) *SellerImpl {
	return &SellerImpl{repo: repo, ps: prod}
}

func (s *SellerImpl) Create(ctx context.Context, createSellerDto dto.CreateSeller, user oauth.OauthUser) (err error) {
	seller := model.Seller{Name: createSellerDto.Name, OwnerId: user.Id}
	err = s.repo.Save(ctx, seller)
	return
}

func (s *SellerImpl) GetAll(ctx context.Context) (sellers []model.Seller, err error) {
	return s.repo.GetAll(ctx)
}

func (s *SellerImpl) Update(ctx context.Context, sellerId uuid.UUID, updateSellerDto dto.UpdateSeller, user oauth.OauthUser) (err error) {
	fromDb, err := s.repo.GetById(ctx, sellerId)

	if err != nil {
		return fmt.Errorf("not found")
	}

	if fromDb.OwnerId != user.Id {
		return fmt.Errorf("forbidden")
	}

	seller := model.Seller{OwnerId: sellerId, Name: updateSellerDto.Name}

	return s.repo.Update(ctx, seller)

}
func (s *SellerImpl) GetById(ctx context.Context, id uuid.UUID) (seller *model.Seller, err error) {
	seller, err = s.repo.GetById(ctx, id)
	return
}

func (s *SellerImpl) AddProduct(ctx context.Context, createProduct dto.CreateProduct, sellerId uuid.UUID) (err error) {
	product := model.Product{
		Name:        createProduct.Name,
		Description: createProduct.Description,
		Price:       createProduct.Price,
		Seller:      model.Seller{OwnerId: sellerId},
		Categories:  make([]model.Category, 0, len(createProduct.CategoryIds)),
	}

	for _, cid := range createProduct.CategoryIds {
		product.Categories = append(product.Categories, model.Category{Id: cid})
	}

	return s.ps.Create(ctx, createProduct, sellerId)
}

func (s *SellerImpl) GetProductPage(ctx context.Context, id uuid.UUID, page int, limit int) (products []model.Product, err error) {
	return s.ps.GetPageBySellerId(ctx, page, limit, id)
}
