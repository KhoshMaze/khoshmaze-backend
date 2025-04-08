package middlewares

import (
	cc "context"
	"fmt"
	"time"

	"github.com/KhoshMaze/khoshmaze-backend/api/utils"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"
	ccache "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cache"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/context"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/permission/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ResourceType string

const (
	ResourceRestaurant ResourceType = "restaurant"
	ResourceBranch     ResourceType = "branch"
	ResourceFood       ResourceType = "food"
)

func ResourceControl(db *gorm.DB, cache cache.Provider, resourceType ResourceType, paramName string) fiber.Handler {

	return func(c *fiber.Ctx) error {

		userClaims := utils.UserClaimsFromLocals(c)
		if userClaims == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if userClaims.HasRole(model.SuperAdmin | model.Founder) {
			return c.Next()
		}

		resourceID, err := c.ParamsInt(paramName)

		if err != nil {
			return fiber.ErrBadRequest
		}
		ctx := c.UserContext()
		logger := context.GetLogger(ctx)
		context.SetLogger(ctx, logger.With("resource", resourceType, "resource_id", resourceID))

		// Check if resource exists in cache
		oc := ccache.NewJsonObjectCacher[*resourceHierarchy](cache)
		cached, err := oc.Get(ctx, fmt.Sprintf("resources:%s:%d", resourceType, resourceID))
		if err != nil {
			logger.Error("Unexpected Error in resource cache", "error", err.Error())
		}

		if cached != nil {
			if userClaims.HasRole(model.RestaurantOwner) {
				if cached.OwnerID == userClaims.UserID {
					return c.Next()
				}
			}
			// TODO: Check for employees
		}

		// Checks for owners
		if userClaims.HasRole(model.RestaurantOwner) {
			repo := &resourceRepo{db: db}
			resourceHierarchy, err := repo.getResourceHierarchy(ctx, resourceType, uint(resourceID))

			if err != nil {
				logger.Error("Unexpected Error in food service", "error", err.Error())
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			
			oc.Set(ctx, fmt.Sprintf("resources:%s:%d", resourceType, resourceID), 5*time.Minute, resourceHierarchy)
			if resourceHierarchy.OwnerID != userClaims.UserID {
				logger.Warn("Forbidden access attempt")
				return c.SendStatus(fiber.StatusForbidden)
			}
			return c.Next()
		}

		// TODO: Check for employees
		switch resourceType {

		case ResourceBranch:

		case ResourceFood:
		}
		return c.Next()
	}

}

type resourceHierarchy struct {
	FoodID       uint
	BranchID     uint
	RestaurantID uint
	OwnerID      uint
}

type resourceRepo struct {
	db *gorm.DB
}

func (r *resourceRepo) getResourceHierarchy(ctx cc.Context, resourceType ResourceType, resourceID uint) (*resourceHierarchy, error) {

	var hierarchy resourceHierarchy

	query := r.db.WithContext(ctx)

	switch resourceType {

	case ResourceFood:
		query = query.Raw(`
			SELECT 
				f.id as food_id,
				b.id as branch_id,
				r.id as restaurant_id,
				r.owner_id as owner_id
			FROM foods f 
			INNER JOIN branches b ON f.branch_id = b.id
			INNER JOIN restaurants r ON b.restaurant_id = r.id
			WHERE f.id = ?
		`, resourceID)

	case ResourceBranch:
		query = query.Raw(`
			SELECT 
				NULL as food_id,
				b.id as branch_id,
				r.id as restaurant_id,
				r.owner_id as owner_id
			FROM branches b 
			INNER JOIN restaurants r ON b.restaurant_id = r.id
			WHERE b.id = ?
		`, resourceID)

	case ResourceRestaurant:
		query = query.Raw(`
			SELECT 
				NULL as food_id,
				NULL as branch_id,
				r.id as restaurant_id,
				r.owner_id as owner_id
			FROM restaurants r
			WHERE r.id = ?
		`, resourceID)
	}

	if err := query.Scan(&hierarchy).Error; err != nil {
		return nil, err
	}

	return &hierarchy, nil
}
