package usecase

import (
	"context"
	"store/cart_service/internal/model"
	"store/cart_service/internal/model/dto"
	"store/cart_service/internal/oauth"
	"store/cart_service/internal/usecase/api"
	"store/cart_service/internal/usecase/repo"

	"github.com/google/uuid"
)

type CartUsecase struct {
	repo       *repo.Cart
	productApi *api.Product
}

func NewCartUsecase(repo *repo.Cart, productApi *api.Product) *CartUsecase {
	return &CartUsecase{repo: repo, productApi: productApi}
}

func (uc *CartUsecase) AddOrUpdate(ctx context.Context, entries []dto.AddToCartRequest, user oauth.OauthUser) (err error) {
	allIds := make([]uuid.UUID, 0, len(entries))
	entriesMap := make(map[uuid.UUID]dto.AddToCartRequest)
	for _, e := range entries {
		entriesMap[e.ProductId] = e
		allIds = append(allIds, e.ProductId)
	}

	cachedIds, err := uc.repo.CheckExist(ctx, allIds)
	if err != nil {
		return err
	}

	toUpdate := make([]dto.AddToCartRequest, 0, len(entries))
	toFetch := make([]uuid.UUID, 0, len(entries)-len(cachedIds))
	for _, e := range entries {
		if _, ok := cachedIds[e.ProductId]; ok {
			toUpdate = append(toUpdate, e)
		} else {
			toFetch = append(toFetch, e.ProductId)
		}
	}

	fetched, err := uc.productApi.ExistBatch(ctx, toFetch)
	if err != nil {
		return err
	}

	uc.repo.CacheProducts(ctx, fetched)

	for _, p := range fetched {
		toUpdate = append(toUpdate, entriesMap[p.Id])
	}

	_ = uc.repo.Upsert(ctx, user.Id, toUpdate)
	return nil
}
func (uc *CartUsecase) GetCart(ctx context.Context, userId uuid.UUID) (entries []model.CartEntry, err error) {
	return uc.repo.GetCart(ctx, userId)
}
func (uc *CartUsecase) Remove(ctx context.Context, productIds []uuid.UUID, user oauth.OauthUser) (err error) {
	return uc.repo.RemoveFromCart(ctx, user.Id, productIds)
}
func (uc *CartUsecase) GetProducts(ctx context.Context, user oauth.OauthUser) (products []model.ProductFull, err error) {
	ids, err := uc.repo.GetCardIds(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	products, err = uc.productApi.GetBatch(ctx, ids)
	if err != nil {
		return nil, err
	}

	updateCache := make([]model.Product, 0, len(products))

	for _, p := range products {
		updateCache = append(updateCache, model.Product{Id: p.Id, Name: p.Name})
	}

	uc.repo.UpdateCache(ctx, updateCache)
	return products, nil
}
