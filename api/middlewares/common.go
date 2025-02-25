package middlewares

import (
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
	"github.com/gofiber/fiber/v2"
)

func RequestsLogger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		logger := context.GetLogger(ctx.UserContext())
		context.SetLogger(ctx.UserContext(), logger.With("ip", ctx.IP()))
		return ctx.Next()

	}
}
