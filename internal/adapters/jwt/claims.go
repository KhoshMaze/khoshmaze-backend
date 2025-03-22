package jwt

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/permission/model"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt5.RegisteredClaims
	UserID      uint
	Permissions uint64
	Phone       string
	IP          string
}

func (c *UserClaims) ConvertToPermissionDomain() model.UserPermissions {
	return model.UserPermissions(c.Permissions)
}

// func (c *UserClaims) HasPermission(permission model.Permission) bool {
// 	return c.ConvertToPermissionDomain().HasPermission(permission)
// }

func (c *UserClaims) HasAllPermissions(permissions ...model.Permission) bool {
	return c.ConvertToPermissionDomain().HasAllPermissions(permissions...)
}
