package oauth

import (
	"slices"

	"github.com/google/uuid"
)

type OauthUser struct {
	Id    uuid.UUID
	Roles []string
}

func (u *OauthUser) HasRole(roleName string) bool {

	return slices.Contains(u.Roles, roleName)
}
