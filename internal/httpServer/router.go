package router

import (
	"github.com/MapMinder/mapminder_backend/feature/health"

	"github.com/gin-gonic/gin"
)

// 各機能のハンドラーの呼び出し
func NewRoutes(r *gin.Engine) {
	health.NewHealthHandler(r)
}
