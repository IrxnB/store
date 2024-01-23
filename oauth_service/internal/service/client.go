package service

import (
	"context"
	"fmt"
	"oauth2_provider/internal/encoding"
	"oauth2_provider/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ClientService interface {
	GetById(ctx context.Context, id uuid.UUID) (*model.Client, error)
	GetWithValidation(ctx context.Context, id uuid.UUID, secret string) (*model.Client, error)
	Save(ctx context.Context, client *model.Client) error
}

type ClientServiceImpl struct {
	Db *pgxpool.Pool
}

func (cs *ClientServiceImpl) GetById(ctx context.Context, id uuid.UUID) (*model.Client, error) {
	var client model.Client

	err := cs.Db.QueryRow(ctx, `
		SELECT id, secret, domain, type
		FROM client
		WHERE id = $1
	`, id).Scan(&client.Id, &client.Secret, &client.Domain, &client.Type)

	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (cs *ClientServiceImpl) GetWithValidation(ctx context.Context, id uuid.UUID, secret string) (*model.Client, error) {

	fromDb, err := cs.GetById(ctx, id)

	if err != nil {
		return nil, err
	}

	if !encoding.Check(secret, fromDb.Secret) {
		return nil, fmt.Errorf("notFound")
	}

	return fromDb, nil
}

func (cs *ClientServiceImpl) Save(ctx context.Context, client *model.Client) error {

	secretEnc, err := encoding.Encode(client.Secret)

	if err != nil {
		return err
	}
	_, err = cs.Db.Exec(ctx, `
		INSERT INTO client(id, secret, domain, type)
		VALUES($1, $2, $3, $4)
	`, &client.Id, &secretEnc, &client.Domain, &client.Type)

	return err
}
