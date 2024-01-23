package repo

import (
	"context"
	"store/cart_service/internal/model"
	"store/cart_service/internal/model/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Cart struct {
	db *pgxpool.Pool
}

func NewCart(db *pgxpool.Pool) (*Cart, error) {
	cart := &Cart{db: db}
	err := cart.db.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (c *Cart) Upsert(ctx context.Context, userId uuid.UUID, entries []dto.AddToCartRequest) (err error) {
	panic("not implemented")
}

func (c *Cart) CheckExist(ctx context.Context, productIds []uuid.UUID) (exist map[uuid.UUID]interface{}, err error) {
	panic("not implemented")
}

func (c *Cart) CacheProducts(ctx context.Context, products []model.Product) (err error) {
	panic("not implemented")
}

func (c *Cart) GetCart(ctx context.Context, userId uuid.UUID) (entries []model.CartEntry, err error) {
	panic("not implemented")
}

func (c *Cart) RemoveFromCart(ctx context.Context, userId uuid.UUID, productIds []uuid.UUID) (err error) {
	panic("not implemented")
}

func (c *Cart) GetCardIds(ctx context.Context, userId uuid.UUID) (ids []uuid.UUID, err error) {
	panic("not implemented")
}

func (c *Cart) UpdateCache(ctx context.Context, products []model.Product) (err error) {
	panic("not implemented")
}
