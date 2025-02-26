package http

import (
	"fmt"

	"github.com/KhoshMaze/khoshmaze-backend/api/http/handlers"
	"github.com/KhoshMaze/khoshmaze-backend/api/middlewares"
	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/app"
	json "github.com/goccy/go-json"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {

	router := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

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
	// return router.ListenTLS(fmt.Sprintf(":%d", cfg.Port), cfg.SSLCertPath, cfg.SSLKeyPath)
	return router.Listen(fmt.Sprintf(":%d", cfg.Port))

}

func registerAuthAPI(appContainer app.App, cfg config.ServerConfig, router fiber.Router) {
	userSvcGetter := handlers.UserServiceGetter(appContainer, cfg)
	router.Post("/register", middlewares.SetTransaction(appContainer.DB()), handlers.SignUp(userSvcGetter))
	router.Post("/refresh", middlewares.RateLimiter("refreshToken", int((cfg.AuthExpMinute-1)*60), 1), handlers.RefreshToken(userSvcGetter))
	router.Post("/send-otp", middlewares.RateLimiter("", 140, 1), handlers.SendOTP(userSvcGetter))
	router.Post("/login", handlers.Login(userSvcGetter))

	router.Use(middlewares.AuthMiddleware([]byte(cfg.AuthSecret), []byte(cfg.RefreshSecret)))
	router.Post("/logout", handlers.Logout(userSvcGetter))
	router.Post("/qrcode", handlers.GenerateQrCode())
	router.Get("/test", handlers.Test())
}
