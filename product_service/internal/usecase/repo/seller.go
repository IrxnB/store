package repo

import (
	"context"
	"product_service/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Seller struct {
	db *pgxpool.Pool
}

func NewSeller(db *pgxpool.Pool) (sellerRepo *Seller, err error) {
	err = db.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	return &Seller{db: db}, nil
}

func (s *Seller) Save(ctx context.Context, seller model.Seller) (err error) {
	_, err = s.db.Exec(ctx, `
		INSERT INTO seller(id, name)
		VALUES ($1, $2)
	`, seller.OwnerId, seller.Name)
	if err != nil {
		return err
	}
	return nil
}

func (s *Seller) GetAll(ctx context.Context) (sellers []model.Seller, err error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, name
		FROM seller
	`)

	if err != nil {
		return nil, err
	}

	sellers = make([]model.Seller, 0)

	for rows.Next() {
		var cur model.Seller

		err = rows.Scan(&cur.OwnerId, &cur.Name)
		if err != nil {
			return nil, err
		}

		sellers = append(sellers, cur)
	}
	return sellers, nil
}

func (s *Seller) Update(ctx context.Context, seller model.Seller) (err error) {
	panic("not implemented")
}

func (s *Seller) GetById(ctx context.Context, id uuid.UUID) (seller *model.Seller, err error) {
	seller = &model.Seller{}
	err = s.db.QueryRow(ctx, `
		SELECT id, name
		FROM seller
		WHERE id = $1
	`, id).Scan(&seller.OwnerId, &seller.Name)

	if err != nil {
		return nil, err
	}
	return
}
