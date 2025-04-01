package jwt

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/permission/model"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt5.RegisteredClaims
	UserID      uint
	Permissions uint64
	Roles       uint64
	Phone       string
	IP          string
}

func (c *UserClaims) ConvertToAuthority() model.Authority {
	return model.Authority(c.Permissions)
}

func (c *UserClaims) HasRole(role model.UserRoles) bool {
	return model.Authority(c.Roles).HasSpecific(role)
}

func (c *UserClaims) HasAllPermissions(permissions ...model.Authority) bool {
	return c.ConvertToAuthority().HasAll(permissions...)
}
