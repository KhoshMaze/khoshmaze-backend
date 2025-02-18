package http

import (
	"fmt"

	"github.com/KhoshMaze/khoshmaze-backend/api/http/handlers"
	"github.com/KhoshMaze/khoshmaze-backend/api/middlewares"
	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/app"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {

	router := fiber.New()

	router.Use(swagger.New(
		swagger.Config{
			// BasePath: "/docs",
			FilePath: "./docs/swagger.yaml",
			// Path: "v1",
			CacheAge: 0,
		},
	))
	api := router.Group("/api/v1", middlewares.SetUserContext)
	api.Get("/metrics", monitor.New())
	// api.Use(limiter.New())

	registerAuthAPI(appContainer, cfg, api)

	return router.Listen(fmt.Sprintf(":%d", cfg.Port))

}

func registerAuthAPI(appContainer app.App, cfg config.ServerConfig, router fiber.Router) {
	userSvcGetter := handlers.UserServiceGetter(appContainer, cfg)
	router.Post("/register", middlewares.SetTransaction(appContainer.DB()), handlers.SignUp(userSvcGetter))
	router.Post("/refresh", handlers.RefreshToken(userSvcGetter))
	router.Post("/logout", middlewares.AuthMiddleware([]byte(cfg.AuthSecret)), handlers.Logout(userSvcGetter))
	router.Get("/test", middlewares.AuthMiddleware([]byte(cfg.AuthSecret)), handlers.Test())
	// router.Post("/login")
}
