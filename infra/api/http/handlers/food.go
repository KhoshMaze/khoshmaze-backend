package handlers

import (
	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/menu/model"
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

func CreateFoodImage(svcGetter ServiceGetter[*service.MenuService]) fiber.Handler {

	return func(c *fiber.Ctx) error {
		foodID, err := c.ParamsInt("foodID")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "couldn't parse foodID")
		}

		img := &model.FoodImage{
			FoodID: uint(foodID),
		}

		// TODO: CHECK MIMTYPES AND FILE EXTs later
		data := c.Body()
		// if err != nil {
		// 	return fiber.NewError(fiber.StatusBadRequest, err.Error())
		// }

		// f, err := file.Open()
		// if err != nil {
		// 	return fiber.NewError(fiber.StatusBadRequest, err.Error())
		// }

		// defer f.Close()

		// data, err := io.ReadAll(f)
		// if err != nil {
		// 	return fiber.NewError(fiber.StatusBadRequest, err.Error())
		// }
		img.MIMEType = model.MIMEType(c.Get("Content-Type"))
		img.Image = data

		svc := svcGetter(c.UserContext())

		if err := svc.AddFoodImageToFood(c.UserContext(), img); err != nil {
			// TODO: fix error handling
			return fiber.ErrUnsupportedMediaType
		}

		return c.Status(fiber.StatusCreated).JSON("image has been added")
	}

}

func GetFoodImages(svcGetter ServiceGetter[*service.MenuService]) fiber.Handler {

	return func(c *fiber.Ctx) error {

		foodID, err := c.ParamsInt("foodID")

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		svc := svcGetter(c.UserContext())
		
		page := c.QueryInt("page")
		size := c.QueryInt("size")
		
		// TODO: FIX GET IMAGES (should I return binary or actual image??)
		result, err := svc.GetImagesByFoodID(c.UserContext(), uint(foodID), page, size)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		c.Set("Cache-Control", "public, max-age=1209600")
		return c.Status(fiber.StatusOK).JSON(result)
	}
}

func DeleteFoodImageFromFood(svcGetter ServiceGetter[*service.MenuService]) fiber.Handler {

	return func(c *fiber.Ctx) error {

		id, err := c.ParamsInt("id")

		if err != nil {
			return fiber.ErrBadRequest
		}

		svc := svcGetter(c.UserContext())

		if err := svc.DeleteFoodImageFromFood(c.UserContext(), uint(id)); err != nil {
			return fiber.NewError(fiber.StatusNotFound, "image not found")
		}

		return c.SendStatus(fiber.StatusOK)

	}
}
