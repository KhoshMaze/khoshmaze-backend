package notification

import (
	"encoding/json"

	"github.com/KhoshMaze/khoshmaze-backend/config"
	"github.com/KhoshMaze/khoshmaze-backend/internal/domain/notification/model"
	"github.com/gofiber/fiber/v2"
)

func SendSMS(data *model.OutboxData, cfg config.SMSConfig) {
	agent := fiber.Post(cfg.ApiBaseURL + cfg.Verification.URL)
	// agent.Debug()
	agent.Request().Header.Add("x-api-key", cfg.ApiKey)
	agent.ContentType("application/json")
	body := fiber.Map{
		"mobile":     data.Dest,
		"templateId": cfg.Verification.TemplateID,
		"parameters": []fiber.Map{
			{"name": "Code",
				"value": data.Content,
			},
		},
	}
	d, _ := json.Marshal(body)
	agent.Body(d)
	go agent.Bytes()
}
