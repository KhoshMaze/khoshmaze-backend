package handlers

import (
	"errors"

	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
	"github.com/gofiber/fiber/v2"
)

// SignUp handles user registration requests.
// @Summary Register a new user
// @Description Register a new user with phone and OTP
// @Tags users
// @Accept json
// @Produce json
// @Param request body pb.UserSignUpRequest true "User sign-up request"
// @Success 201 {object} pb.UserTokenResponse
// @Failure 400
// @Failure 500
// @Router /register [post]
func SignUp(svcGetter ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req pb.UserSignUpRequest
		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		resp, err := svc.SignUp(c.UserContext(), &req, c.IP())
		if err != nil {
			if errors.Is(err, service.ErrUserAlreadyExists) {
				return fiber.NewError(fiber.StatusConflict, err.Error())
			}
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
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

// SendOTP handles requests to send a one-time password (OTP) to the user.
// @Summary Send OTP to user
// @Description Send a one-time password to the user's phone
// @Tags users
// @Accept json
// @Produce json
// @Param request body pb.OtpRequest true "OTP request"
// @Success 200 {object} map[string]string
// @Failure 400
// @Router /send-otp [post]
func SendOTP(svcGetter ServiceGetter[*service.UserService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req pb.OtpRequest

		if err := c.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		category, err := svc.SendOTP(c.UserContext(), &req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(
			fiber.Map{
				"message":  "sent otp",
				"status":   "ok",
				"category": category,
			},
		)
	}
}

// Logout handles user logout requests.
// @Summary Logout user
// @Description Invalidate the user's refresh token and log them out
// @Tags users
// @Produce json
// @Success 202 {object} map[string]string
// @Failure 400
// @Router /logout [post]
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

// Login handles user login requests.
// @Summary Login user
// @Description Authenticate the user and issue a refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param request body pb.UserLoginRequest true "User login request"
// @Success 202 {object} pb.UserTokenResponse
// @Failure 400
// @Failure 500
// @Router /login [post]
func Login(svcGetter ServiceGetter[*service.UserService]) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		svc := svcGetter(ctx.UserContext())
		var req pb.UserLoginRequest

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		resp, err := svc.Login(ctx.UserContext(), &req, ctx.IP())
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
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

// RefreshToken handles requests to refresh the user's authentication token.
// @Summary Refresh authentication token
// @Description Validate the current refresh token and issue a new one
// @Tags users
// @Produce json
// @Success 200 {object} pb.UserTokenResponse
// @Failure 400
// @Router /refresh [post]
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

// Test is a simple handler for testing purposes.
// @Summary Test endpoint
// @Description A simple test endpoint
// @Tags test
// @Produce json
// @Success 200 {object} map[string]string
// @Router /test [get]
func Test() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger := context.GetLogger(ctx.UserContext())
		logger.Info("test")
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"test": "ok",
		})
	}
}
