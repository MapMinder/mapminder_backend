package router

import (
	"github.com/MapMinder/mapminder_backend/feature/health"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// 各機能のハンドラーの呼び出し
func NewRoutes(r *gin.Engine, dbHandler *gorm.DB) {
	health.NewHealthHandler(r)
}
