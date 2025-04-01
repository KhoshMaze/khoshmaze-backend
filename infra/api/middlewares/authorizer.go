package middlewares

import (
	"github.com/KhoshMaze/khoshmaze-backend/api/utils"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/permission/model"
	"github.com/gofiber/fiber/v2"
)

func Authorizer(role model.UserRoles,permissions ...model.UserPermissions) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userClaims := utils.UserClaimsFromLocals(c)
		if userClaims == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if !userClaims.HasRole(role) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		if !userClaims.HasAllPermissions(permissions...) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		return c.Next()
	}
}
