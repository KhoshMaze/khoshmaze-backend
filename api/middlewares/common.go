package middlewares

import (
	"time"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func RequestsLogger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		logger := context.GetLogger(ctx.UserContext())
		context.SetLogger(ctx.UserContext(), logger.With("ip", ctx.IP()))
		return ctx.Next()

	}
}

func RateLimiter(key string, exp, max int) fiber.Handler {
	return limiter.New(limiter.Config{
		KeyGenerator: func(c *fiber.Ctx) string {
			if key == "" {
				key = c.IP()
			}
			return key
		},
		Expiration: time.Duration(exp) * time.Second,
		Max:        max,
	},
	)
}
