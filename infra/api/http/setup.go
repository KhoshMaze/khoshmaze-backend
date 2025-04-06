package http

import (
	"fmt"

	"github.com/KhoshMaze/khoshmaze-backend/api/http/handlers"
	"github.com/KhoshMaze/khoshmaze-backend/api/middlewares"
	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/app"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/permission/model"
	json "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

/*

	Define permissions here to avoid redundancy
	R = Read
	W = Write
	D = Delete
	U = Update

*/

var (
	userSvcGetter       handlers.ServiceGetter[*service.UserService]
	restaurantSvcGetter handlers.ServiceGetter[*service.RestaurantService]
	menuSvcGetter       handlers.ServiceGetter[*service.MenuService]
	authSecret          []byte
	refreshSecret       []byte
	aesSecret           string

	authenticator   = middlewares.AuthMiddleware
	authorizer      = middlewares.Authorizer
	setTransaction  = middlewares.SetTransaction
	setContext      = middlewares.SetUserContext
	rateLimiter     = middlewares.RateLimiter
	resourceControl = middlewares.ResourceControl
)

func Run(appContainer app.App, cfg config.ServerConfig) error {

	// services init
	userSvcGetter = handlers.UserServiceGetter(appContainer, cfg)
	restaurantSvcGetter = handlers.RestaurantServiceGetter(appContainer)
	menuSvcGetter = handlers.MenuServiceGetter(appContainer)

	authSecret = []byte(cfg.AuthSecret)
	refreshSecret = []byte(cfg.RefreshSecret)
	aesSecret = cfg.AESSecret

	router := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	router.Use(healthcheck.New())
	router.Get("/docs/*", fiberSwagger.WrapHandler)

	api := router.Group("/api/v1", setContext)

	registerGlobalRoutes(appContainer, api)
	registerProtectedRoutes(appContainer, api)
	// return router.ListenTLS(fmt.Sprintf(":%d", cfg.Port), cfg.SSLCertPath, cfg.SSLKeyPath)
	return router.Listen(fmt.Sprintf(":%d", cfg.Port))

}

func registerGlobalRoutes(appContainer app.App, router fiber.Router) {

	router.Post("/register", setTransaction(appContainer.DB()), handlers.SignUp(userSvcGetter))
	router.Post("/login", handlers.Login(userSvcGetter))
	router.Post("/send-otp", setTransaction(appContainer.DB()), rateLimiter("", 140, 1), handlers.SendOTP(userSvcGetter))
	router.Get("/:name/:id<int>", handlers.GetBranch(restaurantSvcGetter))
	router.Get("/menus/:menuID<int>/foods", handlers.GetFoods(menuSvcGetter))
}

func registerProtectedRoutes(appContainer app.App, router fiber.Router) {
	cfg := appContainer.Config().Server
	anomalyDetection := appContainer.AnomalyDetectionService()
	router.Use(anomalyDetection.DetectAnomalyMiddleware(refreshSecret))

	router.Post("/refresh", rateLimiter("refreshToken", int((cfg.AuthExpMinute-1)*60), 1),
		handlers.RefreshToken(userSvcGetter))

	router.Use(authenticator(authSecret, false, aesSecret))
	router.Get("/metrics", authorizer(model.Founder), monitor.New())

	registerUserEndpoints(appContainer, router)
	registerRestaurantEndpoints(appContainer, router)
	registerFoodEndpoints(appContainer, router)
}

func registerUserEndpoints(appContainer app.App, router fiber.Router) {
	router.Post("/logout", setTransaction(appContainer.DB()), handlers.Logout(userSvcGetter))
	router.Post("/qrcode", handlers.GenerateQrCode())
}

func registerRestaurantEndpoints(appContainer app.App, router fiber.Router) {

	router = router.Group("/restaurants")
	router.Post("/", setTransaction(appContainer.DB()), handlers.CreateRestaurant(restaurantSvcGetter))
	router.Get("/", authorizer(model.Founder), handlers.GetRestaurants(restaurantSvcGetter))
}

func registerFoodEndpoints(appContainer app.App, router fiber.Router) {
	router = router.Group("/branch/:branchID<int>")
	// add resource controller middleware
	// router.Use()

	router.Get("/food/:foodID<int>", handlers.GetFood(menuSvcGetter))
	router.Post("/food", setTransaction(appContainer.DB()), handlers.CreateFood(menuSvcGetter))
	router.Put("/food", setTransaction(appContainer.DB()), handlers.UpdateFood(menuSvcGetter))
	router.Delete("/food/:foodID<int>", handlers.DeleteFood(menuSvcGetter))
}
