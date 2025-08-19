package infrastructure

import (
	"fmt"
	"log"

	"github.com/MapMinder/mapminder_backend/internal/config"
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to initilize DB: ", err)
	}

	return db
}
