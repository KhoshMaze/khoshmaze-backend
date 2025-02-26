package handlers

import (
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
	"github.com/gofiber/fiber/v2"
)

func SignUp(svcGetter ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req pb.UserSignUpRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		resp, err := svc.SignUp(c.UserContext(), &req, c.IP())
		if err != nil {
			if errors.Is(err, service.ErrUserAlreadyExists) || errors.Is(err, service.ErrWrongOTP) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		c.Cookie(&fiber.Cookie{
			Name:     "refreshToken",
			Value:    resp.RefreshToken,
			HTTPOnly: true,
			Secure:   true,
			SameSite: "strict",
			MaxAge:   int(resp.RefreshMaxAge),
		})

		logger := context.GetLogger(c.UserContext())
		logger.Info("new user created")
		resp.RefreshToken = ""
		resp.RefreshMaxAge = 0
		return c.Status(fiber.StatusCreated).JSON(resp)

	}
}

func SendOTP(svcGetter ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req pb.OtpRequest

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}
		if err := svc.SendOTP(c.UserContext(), &req); err != nil {
			return err
		}

		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"message": "sent otp",
				"status":  "ok",
			},
		)
	}
}

func Logout(svcGetter ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())

		err := svc.Logout(c.UserContext(), c.Cookies("refreshToken"))
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		logger := context.GetLogger(c.UserContext())
		logger.Info("user logged out")

		return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "ok",
		})
	}
}

func Login(svcGetter ServiceGetter[*service.UserService]) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		svc := svcGetter(ctx.UserContext())
		var req pb.UserLoginRequest

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		resp, err := svc.Login(ctx.UserContext(), &req, ctx.IP())
		if err != nil {
			return fiber.ErrInternalServerError
		}

		logger := context.GetLogger(ctx.UserContext())
		logger.Info("user logged in")

		ctx.Cookie(&fiber.Cookie{
			Name:     "refreshToken",
			Value:    resp.RefreshToken,
			HTTPOnly: true,
			Secure:   true,
			SameSite: "strict",
			MaxAge:   int(resp.RefreshMaxAge),
		})
		resp.RefreshToken = ""
		resp.RefreshMaxAge = 0
		return ctx.Status(fiber.StatusAccepted).JSON(resp)
	}
}
func RefreshToken(svcGetter ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		svc := svcGetter(c.UserContext())

		resp, err := svc.RefreshToken(c.UserContext(), c.Cookies("refreshToken"), c.IP())

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if resp.GetRefreshToken() != "" {

			c.Cookie(&fiber.Cookie{
				Name:     "refreshToken",
				Value:    resp.RefreshToken,
				HTTPOnly: true,
				Secure:   true,
				SameSite: "strict",
				MaxAge:   int(resp.RefreshMaxAge),
			})
		}
		resp.RefreshToken = ""
		resp.RefreshMaxAge = 0
		return c.Status(fiber.StatusOK).JSON(resp)

	}
}

func Test() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger := context.GetLogger(ctx.UserContext())
		logger.Info("test")
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"test": "ok",
		})
	}
}
