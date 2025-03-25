package http

import (
	"fmt"

	"github.com/KhoshMaze/khoshmaze-backend/api/http/handlers"
	"github.com/KhoshMaze/khoshmaze-backend/api/middlewares"
	"github.com/KhoshMaze/khoshmaze-backend/api/service"
	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/app"
	perm "github.com/KhoshMaze/khoshmaze-backend/internal/domain/permission/model"
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

const (
	userR    = perm.ReadUser
	userRW   = perm.ReadUser + perm.WriteUser
	adminRWD = perm.ReadAdmin + perm.WriteAdmin + perm.DeleteAdmin
	adminRW  = perm.ReadAdmin + perm.WriteAdmin
)

func Run(appContainer app.App, cfg config.ServerConfig) error {

	router := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	router.Use(healthcheck.New())
	router.Get("/docs/*", fiberSwagger.WrapHandler)
	api := router.Group("/api/v1", middlewares.SetUserContext)
	api.Get("/metrics", monitor.New())

	userSvcGetter := handlers.UserServiceGetter(appContainer, cfg)
	restaurantSvcGetter := handlers.RestaurantServiceGetter(appContainer)
	registerGlobalRoutes(appContainer, cfg, api, userSvcGetter)

	// Authentication Required Endpoints
	secret := []byte(cfg.AuthSecret)
	router.Use(middlewares.AuthMiddleware(secret))

	registerUserEndpoints(appContainer, api, userSvcGetter)
	registerRestaurantEndpoints(appContainer, api, restaurantSvcGetter)
	// return router.ListenTLS(fmt.Sprintf(":%d", cfg.Port), cfg.SSLCertPath, cfg.SSLKeyPath)
	return router.Listen(fmt.Sprintf(":%d", cfg.Port))

}

func registerUserEndpoints(appContainer app.App, router fiber.Router, userSvcGetter handlers.ServiceGetter[*service.UserService]) {
	router.Post("/logout", middlewares.SetTransaction(appContainer.DB()), handlers.Logout(userSvcGetter))
	router.Post("/qrcode", handlers.GenerateQrCode())
	router.Get("/test", handlers.Test())
}

func registerRestaurantEndpoints(appContainer app.App, router fiber.Router, restaurantSvcGetter handlers.ServiceGetter[*service.RestaurantService]) {
	
	router = router.Group("/restaurants")
	router.Get("/:name/:id<int>", handlers.GetBranch(restaurantSvcGetter))
	router.Post("/", middlewares.SetTransaction(appContainer.DB()), handlers.CreateRestaurant(restaurantSvcGetter))
	router.Get("/", middlewares.Authorizer(adminRW), handlers.GetRestaurants(restaurantSvcGetter))
}

func registerGlobalRoutes(appContainer app.App, cfg config.ServerConfig, router fiber.Router, userSvcGetter handlers.ServiceGetter[*service.UserService]) {
	router.Post("/register", middlewares.SetTransaction(appContainer.DB()), handlers.SignUp(userSvcGetter))
	router.Post("/login", handlers.Login(userSvcGetter))
	router.Post("/send-otp", middlewares.RateLimiter("", 140, 1), handlers.SendOTP(userSvcGetter))
	router.Post("/refresh", middlewares.RateLimiter("refreshToken", int((cfg.AuthExpMinute-1)*60), 1), handlers.RefreshToken(userSvcGetter))
}
