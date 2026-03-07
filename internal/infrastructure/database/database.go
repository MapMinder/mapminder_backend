package infrastructure

import (
	"fmt"

	"github.com/MapMinder/mapminder_backend/internal/config"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// NewDBHanler
func NewDBHandler(cfg *config.DBConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	logger.Info("connecting to the database...")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("Failed to initilize DB: %v", err)
	}

	logger.Info("connected to the database")

	return db
}
