package middlewares

import (
	"github.com/KhoshMaze/khoshmaze-backend/api/utils"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/permission/model"
	resModel "github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/model"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/restaurant/port"
	"github.com/gofiber/fiber/v2"
)

type ResourceType string

const (
	ResourceRestaurant ResourceType = "restaurant"
	ResourceBranch     ResourceType = "branch"
	ResourceMenu       ResourceType = "menu"
	ResourceFood       ResourceType = "food"
)

func ResourceControl(restaurantSvc port.Service, resourceType ResourceType, paramName string) fiber.Handler {

	return func(c *fiber.Ctx) error {

		userClaims := utils.UserClaimsFromLocals(c)
		if userClaims == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if userClaims.HasRole(model.SuperAdmin + model.Founder) {
			return c.Next()
		}

		resourceID, err := c.ParamsInt(paramName)

		if err != nil {
			return fiber.ErrBadRequest
		}
		ctx := c.UserContext()
		logger := context.GetLogger(ctx)
		context.SetLogger(ctx, logger.With("resource", resourceType, "resource_id", resourceID))
		switch resourceType {
		case ResourceRestaurant:

			restaurant, err := restaurantSvc.GetRestaurantByFilter(ctx, &resModel.RestaurantFilter{
				ID: uint(resourceID),
			})

			if err != nil {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			if restaurant.OwnerID != userClaims.UserID {
				logger.Warn("Forbidden access attempt")
				return c.SendStatus(fiber.StatusForbidden)
			}

		default:

			branch, err := restaurantSvc.GetBranchByFilter(ctx, &resModel.BranchFilter{
				ID: uint(resourceID),
			})

			if err != nil {
				return c.SendStatus(fiber.StatusBadRequest)
			}

			if userClaims.HasRole(model.RestaurantOwner) {
				restaurant, err := restaurantSvc.GetRestaurantByFilter(ctx, &resModel.RestaurantFilter{
					ID: branch.RestaurantID,
				})

				if err != nil {
					logger.Error("Unexpected Error in restaurant service", "error", err.Error())
					return c.SendStatus(fiber.StatusInternalServerError)
				}

				if restaurant.OwnerID != userClaims.UserID {
					logger.Warn("Forbidden access attempt")
					return c.SendStatus(fiber.StatusForbidden)
				}
				return c.Next()
			}

			// TODO: IMPLEMENT EMPLOYEES LOGIC FOR BRANCHES


		}
		return c.Next()
	}

}
