package cron

import (
	"log"
	"time"

	"github.com/KhoshMaze/khoshmaze-backend/internal/adapters/storage/types"
	timeutils "github.com/KhoshMaze/khoshmaze-backend/internal/adapters/time"
	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

func SetTokenDeleterJob(db *gorm.DB, interval int) {
	job := gocron.NewScheduler(timeutils.TehranLoc)
	job.Every(interval).Hours().Do(func() {
		log.Println("started cron job")
		err := db.Where("expires_at < ?", time.Now()).Delete(&types.TokenBlacklist{}).Error
		if err != nil {
			log.Fatal("TOKEN DELETER CRON JOB FAILED", err)
		}
	})
	job.StartAsync()
}
