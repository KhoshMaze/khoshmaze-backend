package notification

import (
	"encoding/json"
	"strconv"

	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/model"
	"github.com/gofiber/fiber/v2"
)

func SendSMS(data *model.OutboxData, cfg *config.SMSConfig) {
	agent := fiber.Post(cfg.ApiBaseURL + cfg.Verification.URL)
	// agent.Debug()
	agent.Request().Header.Add("apikey", cfg.ApiKey)
	agent.ContentType("application/json")
	code, _ := strconv.Atoi(data.Content)
	body := fiber.Map{
		"code":      cfg.Verification.TemplateID,
		"recipient": data.Dest,
		"sender":    cfg.Sender,
		"variable": fiber.Map{
			"verification-code": code,
		},
	}
	d, _ := json.Marshal(body)
	agent.Body(d)
	go agent.Bytes()
}
