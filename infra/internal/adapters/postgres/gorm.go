package postgres

import (
	"fmt"
	"time"

	LG "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConnOptions struct {
	User   string
	Pass   string
	Host   string
	Port   uint
	DBName string
}

func (o *DBConnOptions) PostgresDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		o.Host, o.Port, o.User, o.Pass, o.DBName)
}

func NewPsqlGormConnection(opt DBConnOptions) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(opt.PostgresDSN()), &gorm.Config{
		Logger: logger.New(LG.NewLogger(), logger.Config{
			SlowThreshold: 800 * time.Millisecond,
			// IgnoreRecordNotFoundError: true,
			ParameterizedQueries: true,
			LogLevel:             logger.LogLevel(LG.LoggerConfiguration.Level),
			Colorful:             false,
		}),
	})
}
