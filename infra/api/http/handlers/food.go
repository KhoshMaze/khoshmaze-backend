package handlers

import (
	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/gofiber/fiber/v2"
)

func GetFoods(svcGetter ServiceGetter[*service.MenuService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())

		menuID, err := c.ParamsInt("menuID")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		page := c.QueryInt("page")
		pageSize := c.QueryInt("pagesize")

		foods, err := svc.GetFoods(c.Context(), uint(menuID), page, pageSize)
		return c.Status(fiber.StatusOK).JSON(foods)
	}
}
