package main

import (
	"flag"
	"log"
	"os"

	"github.com/KhoshMaze/khoshmaze-backend/api/http"
	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/cron"
	"github.com/KhoshMaze/khoshmaze-backend/internal/app"
)

var cfg = flag.String("config", "config.json", "server configuration file")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*cfg = v
	}

	c := config.MustReadConfig(*cfg)

	appContainer := app.MustNewApp(c)

	cron.SetTokenDeleterJob(appContainer.DB(), c.Jobs.TokenCheckerInterval)

	log.Fatal(http.Run(appContainer, c.Server))
}
