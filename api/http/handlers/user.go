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

		resp, err := svc.SignUp(c.UserContext(), &req)
		if err != nil {
			if errors.Is(err, errors.New("")) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		logger := context.GetLogger(c.UserContext())
		logger.Info("new user created")
		return c.Status(fiber.StatusCreated).JSON(resp)

	}
}

func RefreshToken(svcGetter ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		svc := svcGetter(c.UserContext())

		resp, err := svc.RefreshToken(c.UserContext(), c.Cookies("refreshToken"))

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

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
