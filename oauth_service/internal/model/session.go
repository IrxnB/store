package model

import "github.com/google/uuid"

type Session struct {
	Id       uuid.UUID
	UserId   uuid.UUID
	ClientId uuid.UUID
	IsUsed   bool
}
