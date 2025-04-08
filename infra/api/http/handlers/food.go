package handlers

import (
	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/gofiber/fiber/v2"
)

func GetFoods(svcGetter ServiceGetter[*service.MenuService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())

		branchID, err := c.ParamsInt("branchID")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		page := c.QueryInt("page")
		pageSize := c.QueryInt("size")

		foods, err := svc.GetFoods(c.UserContext(), uint(branchID), page, pageSize)
		return c.Status(fiber.StatusOK).JSON(foods)
	}
}

func GetFood(svcGetter ServiceGetter[*service.MenuService]) fiber.Handler {

	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		foodId, err := c.ParamsInt("foodID")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		food, err := svc.GetFood(c.UserContext(), uint(foodId))

		if err != nil {
			// TODO: Add custom errors for handling errors
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.Status(fiber.StatusOK).JSON(food)
	}
}

func CreateFood(svcGetter ServiceGetter[*service.MenuService]) fiber.Handler {

	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var food pb.CreateFoodRequest

		if err := c.BodyParser(&food); err != nil {
			return fiber.ErrBadRequest
		}

		id, err := svc.AddFood(c.UserContext(), &food)

		if err != nil {
			return fiber.ErrInternalServerError
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "food created",
			"id":      id,
		})
	}
}

func UpdateFood(svcGetter ServiceGetter[*service.MenuService]) fiber.Handler {

	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())

		var food pb.Food

		if err := c.BodyParser(&food); err != nil {
			return fiber.ErrBadRequest
		}

		err := svc.UpdateFood(c.UserContext(), &food)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
func DeleteFood(svcGetter ServiceGetter[*service.MenuService]) fiber.Handler {

	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())

		id, err := c.ParamsInt("foodID")

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		err = svc.DeleteFood(c.UserContext(), uint(id))

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
