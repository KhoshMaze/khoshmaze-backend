package utils

import (
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/jwt"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/permission/model"
	"github.com/gofiber/fiber/v2"
	jwt5 "github.com/golang-jwt/jwt/v5"
)

func UserClaimsFromLocals(ctx *fiber.Ctx) *jwt.UserClaims {
	if u := ctx.Locals("user"); u != nil {
		userClaims, ok := u.(*jwt5.Token).Claims.(*jwt.UserClaims)
		if ok {
			return userClaims
		}
	}
	return nil
}

func UserClaimsFromCookies(token string, secret []byte) (*jwt.UserClaims, error) {
	if token != "" {
		userclaims, err := jwt.ParseToken(token, secret)
		if err != nil {
			return nil, err
		}
		return userclaims, nil
	}

	return nil, errors.New("empty refresh cookie")
}

func UserRoleToString(role model.UserRoles) string {
	switch {
	case role.HasSpecific(model.Founder):
		return "founder"
	case role.HasSpecific(model.SuperAdmin):
		return "super_admin"
	case role.HasSpecific(model.Accountant):
		return "accountant"
	case role.HasSpecific(model.Support):
		return "support"
	case role.HasSpecific(model.RestaurantOwner):
		return "restaurant_owner"
	case role.HasSpecific(model.BranchManager):
		return "branch_manager"
	case role.HasSpecific(model.Cashier):
		return "cashier"
	case role.HasSpecific(model.Waiter):
		return "waiter"
	case role.HasSpecific(model.BranchAccountant):
		return "branch_accountant"

	default:
		return "unknown"
	}
}
