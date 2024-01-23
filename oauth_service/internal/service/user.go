package service

import (
	"context"
	"fmt"
	"oauth2_provider/internal/encoding"
	"oauth2_provider/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserService interface {
	GrantRole(ctx context.Context, user_id, role_id uuid.UUID) error
	GetWithValidation(ctx context.Context, credentials model.Credentials) (*model.User, error)
	Save(ctx context.Context, user *model.User) error
}

type UserServiceImpl struct {
	Db *pgxpool.Pool
}

func (us *UserServiceImpl) GetWithValidation(ctx context.Context, credentials model.Credentials) (*model.User, error) {

	rows, err := us.Db.Query(ctx, `
		SELECT u.id as user_id, u.username, u.password, r.id as rolde_id, r.name 
		FROM usr u
		LEFT JOIN role_user ru
		ON ru.user_id = u.id
		JOIN role r
		ON r.id = ru.role_id
		WHERE username = $1
	`, credentials.Username)

	fromDb := model.User{Id: uuid.Nil, Roles: make([]model.Role, 0, 1)}

	for rows.Next() {
		var user model.UserWithRole
		rows.Scan(&user.User.Id, &user.User.Username, &user.User.PasswordHash, &user.Role.Id, &user.Role.Name)

		if fromDb.Id == uuid.Nil {
			fromDb.Id = user.User.Id
			fromDb.Username = user.User.Username
			fromDb.PasswordHash = user.User.PasswordHash
		}
		if user.Role.Id != uuid.Nil {
			fromDb.Roles = append(fromDb.Roles, user.Role)
		}
	}

	if err != nil {
		return nil, err
	}

	if !encoding.Check(credentials.Password, fromDb.PasswordHash) {
		return nil, fmt.Errorf("wrong password")
	}

	return &fromDb, nil
}

func (us *UserServiceImpl) GrantRole(ctx context.Context, user_id, role_id uuid.UUID) (err error) {
	_, err = us.Db.Exec(ctx, `
		INSERT INTO role_user(user_id, role_id)
		VALUES ($1, $2)
	`, user_id, role_id)
	return
}

func (us *UserServiceImpl) Save(ctx context.Context, user *model.User) (err error) {

	passwordEnc, err := encoding.Encode(user.PasswordHash)

	if err != nil {
		return
	}

	trx, err := us.Db.Begin(ctx)
	if err != nil {
		return
	}
	defer trx.Commit(ctx)

	_, err = trx.Exec(ctx, `
		INSERT INTO usr(id, username, password)
		VALUES($1, $2, $3)
	`, &user.Id, &user.Username, passwordEnc)

	if err != nil {
		trx.Rollback(ctx)
		return
	}

	if len(user.Roles) > 0 {
		roleNames := make([]string, len(user.Roles))
		for _, r := range user.Roles {
			roleNames = append(roleNames, r.Name)
		}
		statement := `INSERT INTO role_user (user_id, role_id)
			SELECT $1, id FROM role WHERE name = ANY($2)`
		_, err = trx.Exec(ctx, statement, user.Id, roleNames)

		if err != nil {
			trx.Rollback(ctx)
			return
		}
	}
	return
}
