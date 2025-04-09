package config

import (
	"github.com/sandroJayas/user-service/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn := AppConfig.DatabaseURL
	var db *gorm.DB
	var err error

	maxRetries := 10
	retryDelay := time.Second * 2

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			utils.Logger.Sugar().Infof("🚀 Connected to database")
			return db
		}

		utils.Logger.Sugar().Warnf("❌ Failed to connect to DB (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	utils.Logger.Sugar().Fatalf("❌ Could not connect to database after %d attempts: %v", maxRetries, err)
	return nil
}
