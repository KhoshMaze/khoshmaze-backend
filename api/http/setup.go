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

	router.Get("/docs/*", fiberSwagger.WrapHandler)
	api := router.Group("/api/v1", middlewares.SetUserContext)
	api.Get("/metrics", monitor.New())

	userSvcGetter := handlers.UserServiceGetter(appContainer, cfg)
	registerGlobalRoutes(appContainer, cfg, api, userSvcGetter)

	registerUserEndpoints(appContainer, cfg, api, userSvcGetter)
	registerRestaurantEndpoints(appContainer, cfg, api)
	// return router.ListenTLS(fmt.Sprintf(":%d", cfg.Port), cfg.SSLCertPath, cfg.SSLKeyPath)
	return router.Listen(fmt.Sprintf(":%d", cfg.Port))

}

func registerUserEndpoints(appContainer app.App, cfg config.ServerConfig, router fiber.Router, userSvcGetter handlers.ServiceGetter[*service.UserService]) {
	secret := []byte(cfg.AuthSecret)

	router.Use(middlewares.AuthMiddleware(secret, userR))
	router.Post("/logout", middlewares.SetTransaction(appContainer.DB()), handlers.Logout(userSvcGetter))
	router.Post("/qrcode", handlers.GenerateQrCode())
	router.Get("/test", handlers.Test())
}

func registerRestaurantEndpoints(appContainer app.App, cfg config.ServerConfig, router fiber.Router) {
	router = router.Group("/restaurant") 
	// router.Get("/:name/:id<int>", nil) 
}

func registerGlobalRoutes(appContainer app.App, cfg config.ServerConfig, router fiber.Router, userSvcGetter handlers.ServiceGetter[*service.UserService]) {
	router.Post("/register", middlewares.SetTransaction(appContainer.DB()), handlers.SignUp(userSvcGetter))
	router.Post("/login", handlers.Login(userSvcGetter))
	router.Post("/send-otp", middlewares.RateLimiter("", 140, 1), handlers.SendOTP(userSvcGetter))
	router.Post("/refresh", middlewares.RateLimiter("refreshToken", int((cfg.AuthExpMinute-1)*60), 1), handlers.RefreshToken(userSvcGetter))
}
