package middlewares

import (
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/api/utils"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/jwt"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/logger"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	ErrVerifyNewIP = errors.New("unauthorized ip, please verify it with otp (get otp with /api/v1/send-otp and then request to /api/v1/login)")
)

func AuthMiddleware(secret []byte, checkIp bool, aesSecret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: secret},
		Claims:     &jwt.UserClaims{},
		SuccessHandler: func(ctx *fiber.Ctx) error {
			userClaims := utils.UserClaimsFromLocals(ctx)
			if userClaims == nil {
				return fiber.ErrUnauthorized
			}

			unhashedIP, err := utils.DecryptString(userClaims.IP, []byte(aesSecret))
			if err != nil {
				return fiber.ErrInternalServerError
			}
			
			if checkIp && userClaims.IP != unhashedIP {
				return fiber.NewError(fiber.StatusUnauthorized, ErrVerifyNewIP.Error())
			}
			logger := context.GetLogger(ctx.UserContext())
			context.SetLogger(ctx.UserContext(), logger.With("user_id", userClaims.UserID))
			return ctx.Next()
		},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		},
	})
}

func SetUserContext(c *fiber.Ctx) error {
	c.SetUserContext(context.NewAppContext(c.UserContext(), context.WithLogger(logger.NewLogger().With("ip", c.IP()))))
	return c.Next()
}

func SetTransaction(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tx := db.Begin()

		context.SetDB(c.UserContext(), tx, true)

		err := c.Next()

		if c.Response().StatusCode() >= 300 {
			return context.Rollback(c.UserContext())
		}

		if err := context.CommitOrRollback(c.UserContext(), true); err != nil {
			return err
		}

		return err
	}
}
