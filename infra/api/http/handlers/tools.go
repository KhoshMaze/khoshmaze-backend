package handlers

import (
	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/gofiber/fiber/v2"
)

func GenerateQrCode() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		ctx.Set("Content-Type", "image/png")
		var req pb.QrCodeRequest

		if err := ctx.BodyParser(&req); err != nil {
			return fiber.ErrBadRequest
		}

		data, err := service.CreateQR(&req)

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return ctx.Send(data)
	}
}
