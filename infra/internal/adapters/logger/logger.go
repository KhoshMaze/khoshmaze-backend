package logger

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/google/uuid"
	slogbetterstack "github.com/samber/slog-betterstack"
)

type LoggerConfig struct {
	Token    string `json:"token"`
	Endpoint string `json:"endpoint"`
}

var betterstackConfig *LoggerConfig

func init() {
	data, err := os.ReadFile("./logger-config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &betterstackConfig)
	if err != nil {
		panic(err)
	}
}

func NewLogger() *slog.Logger {
	return slog.New(
		slogbetterstack.Option{
			Token:    betterstackConfig.Token,
			Endpoint: betterstackConfig.Endpoint,
		}.NewBetterstackHandler().WithAttrs([]slog.Attr{
			slog.String("trace_id", uuid.NewString()),
		}))
}
