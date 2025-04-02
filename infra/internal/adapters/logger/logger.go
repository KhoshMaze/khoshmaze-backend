package logger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/google/uuid"
	slogbetterstack "github.com/samber/slog-betterstack"
)

type LoggerConfig struct {
	Token    string `json:"token"`
	Endpoint string `json:"endpoint"`
	Level    int    `json:"level"`
	DevMode  bool   `json:"dev_mode"`
}

var LoggerConfiguration *LoggerConfig

func init() {
	data, err := os.ReadFile("./logger-config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &LoggerConfiguration)
	if err != nil {
		panic(err)
	}
}

func NewLogger() *CustomLogger {
	if LoggerConfiguration.DevMode {
		return &CustomLogger{
			Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil).WithAttrs([]slog.Attr{
				slog.String("trace_id", uuid.NewString()),
			})),
			Level: LoggerConfiguration.Level,
		}
	}
	return &CustomLogger{
		Logger: slog.New(
			slogbetterstack.Option{
				Token:    LoggerConfiguration.Token,
				Endpoint: LoggerConfiguration.Endpoint,
			}.NewBetterstackHandler().WithAttrs([]slog.Attr{
				slog.String("trace_id", uuid.NewString()),
			})),
		Level: LoggerConfiguration.Level,
	}
}

type CustomLogger struct {
	*slog.Logger
	Level int
}

func (c *CustomLogger) With(args ...interface{}) *CustomLogger {
	return &CustomLogger{
		Logger: c.Logger.With(args...),
		Level:  c.Level,
	}
}

func (c *CustomLogger) Debug(msg string, args ...interface{}) {
	if slog.Level(c.Level) <= slog.LevelDebug {
		c.Logger.Debug(msg, args...)
	}
}

func (c *CustomLogger) Info(msg string, args ...interface{}) {
	if slog.Level(c.Level) <= slog.LevelInfo {
		c.Logger.Info(msg, args...)
	}
}

func (c *CustomLogger) Warn(msg string, args ...interface{}) {
	if slog.Level(c.Level) <= slog.LevelWarn {
		c.Logger.Warn(msg, args...)
	}
}

func (c *CustomLogger) Error(msg string, args ...interface{}) {
	if slog.Level(c.Level) <= slog.LevelError {
		c.Logger.Error(msg, args...)
	}
}

func (c *CustomLogger) Printf(msg string, args ...interface{}) {
	c.Warn("DATABASE GORM LOG", "msg", fmt.Sprintf("%v", args[1]), "elapsed_time", fmt.Sprintf("%.3fms", args[2]), "rows", fmt.Sprintf("%v", args[3]), "query", fmt.Sprintf("%v", args[4:]...))
}
