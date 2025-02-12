package http

import (
	"fmt"

	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/app"
	"github.com/gofiber/fiber/v2"
)

func Run(appContainer app.App, cfg config.ServerConfig) error {

	router := fiber.New()

	api := router.Group("/api/")
	_ = api
	return router.Listen(fmt.Sprintf("Server is running on port: %d", cfg.Port))

}
