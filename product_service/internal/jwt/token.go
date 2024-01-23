package jwt

import "github.com/google/uuid"

type Token struct {
	TokenId       uuid.UUID
	UserId        uuid.UUID
	ClientId      uuid.UUID
	Scope         string
	ExpiresAt     int64
	UserRoleNames []string
}
