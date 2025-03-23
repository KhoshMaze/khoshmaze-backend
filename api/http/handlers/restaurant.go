package handlers

import (
	"github.com/KhoshMaze/khoshmaze-backend/api/pb"
	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/KhoshMaze/khoshmaze-backend/api/utils"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrRestaurantAlreadyExists = fiber.NewError(fiber.StatusBadRequest, "restaurant already exists")
	ErrRestaurantOnCreate      = fiber.NewError(fiber.StatusInternalServerError, "failed to create restaurant")
	ErrBranchOnCreate          = fiber.NewError(fiber.StatusInternalServerError, "failed to create branch")
	ErrBranchNotFound          = fiber.NewError(fiber.StatusNotFound, "branch not found")
)

// CreateRestaurant handles restaurant creation requests.
// @Summary Create a new restaurant
// @Description Create a new restaurant with name and URL
// @Tags restaurants
// @Accept json
// @Produce json
// @Param request body pb.CreateRestaurantRequest true "Restaurant creation request"
// @Success 201 {object} map[string]uint
// @Failure 400
// @Failure 500
// @Router /restaurants [post]
func CreateRestaurant(svcGetter ServiceGetter[*service.RestaurantService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req pb.CreateRestaurantRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		owner := utils.UserClaimsFromLocals(c)

		if owner == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		id, err := svc.CreateRestaurant(c.Context(), owner.UserID, &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrRestaurantOnCreate,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": id,
		})
	}
}

// GetRestaurants handles requests to get a paginated list of restaurants.
// @Summary Get all restaurants
// @Description Retrieve a paginated list of all restaurants
// @Tags restaurants
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Number of items per page"
// @Success 200 {object} pb.GetAllRestaurantsResponse
// @Failure 500
// @Router /restaurants [get]
func GetRestaurants(svcGetter ServiceGetter[*service.RestaurantService]) fiber.Handler {
	return func(c *fiber.Ctx) error {

		svc := svcGetter(c.UserContext())

		page := c.QueryInt("page")

		pageSize := c.QueryInt("pageSize")
		restaurants, err := svc.GetAllRestaurants(c.Context(), page, pageSize)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(restaurants)
	}
}

// GetBranch handles requests to get a specific branch.
// @Summary Get branch details
// @Description Retrieve details of a specific branch by ID
// @Tags restaurants
// @Produce json
// @Param name path string true "Restaurant name"
// @Param id path int true "Branch ID"
// @Success 200 {object} pb.GetBranchResponse
// @Failure 404
// @Failure 500
// @Router /restaurants/{name}/branches/{id} [get]
func GetBranch(svcGetter ServiceGetter[*service.RestaurantService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		branchID, err := c.ParamsInt("id")

		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		name := c.Params("name")
		branch, err := svc.GetBranch(c.Context(), name, uint(branchID))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrBranchNotFound,
			})
		}
		return c.Status(fiber.StatusOK).JSON(branch)
	}
}

// CreateBranch handles branch creation requests.
// @Summary Create a new branch
// @Description Create a new branch for a restaurant
// @Tags restaurants
// @Accept json
// @Produce json
// @Param request body pb.CreateBranchRequest true "Branch creation request"
// @Success 201 {object} map[string]uint
// @Failure 400
// @Failure 500
// @Router /restaurants/branches [post]
func CreateBranch(svcGetter ServiceGetter[*service.RestaurantService]) fiber.Handler {
	return func(c *fiber.Ctx) error {
		svc := svcGetter(c.UserContext())
		var req pb.CreateBranchRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": ErrBranchOnCreate,
			})
		}

		owner := utils.UserClaimsFromLocals(c)
		if owner == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		id, err := svc.CreateBranch(c.Context(), &req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": ErrBranchOnCreate,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id": id,
		})
	}
}
