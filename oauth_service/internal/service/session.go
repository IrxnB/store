package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type SessionService interface {
	Save(ctx context.Context, CodeId, UserId, ClientId uuid.UUID) error
}

type SessionServiceImpl struct {
	Db *pgxpool.Pool
}

func (ss *SessionServiceImpl) Save(ctx context.Context, CodeId, UserId, ClientId uuid.UUID) error {
	_, err := ss.Db.Exec(ctx, `
		INSERT INTO session(id, client_id, user_id)
		VALUES($1, $2, $3)
	`, CodeId, ClientId, UserId)

	return err
}
