package service

import (
	"context"
	"oauth2_provider/internal/model"

	"github.com/jackc/pgx/v4/pgxpool"
)

type RoleService interface {
	Save(ctx context.Context, role model.Role) error
}

type RoleServiceImpl struct {
	Db *pgxpool.Pool
}

func (rs *RoleServiceImpl) Save(ctx context.Context, role model.Role) (err error) {
	_, err = rs.Db.Exec(ctx, `
		INSERT INTO role(id, name)
		VALUES ($1, $2)
	`, &role.Id, &role.Name)

	return err
}
