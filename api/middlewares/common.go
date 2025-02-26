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
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		KeyGenerator: func(c *fiber.Ctx) string {
			switch key {
			case "refreshToken":
				key = c.Cookies("refreshToken")
			case "phone":
				key = c.FormValue("phone", c.IP())
			default:
				key = c.IP()

				return key
			}
			return key
		},
		Expiration:         time.Duration(exp) * time.Second,
		Max:                max,
		SkipFailedRequests: true,
	},
	)
}
