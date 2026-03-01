package router

import (
	"github.com/MapMinder/mapminder_backend/feature/auth"
	"github.com/MapMinder/mapminder_backend/feature/health"
	"github.com/MapMinder/mapminder_backend/feature/reminder"
	"github.com/MapMinder/mapminder_backend/shared/middleware"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// 各機能のハンドラーの呼び出し
func NewRoutes(r *gin.Engine, dbHandler *gorm.DB) {
	// error handler middleware
	r.Use(middleware.ErrorHandler())

	// routes
	// health
	health.NewHealthHandler(r)

	// auth
	authGroup := r.Group("/auth")
	auth.NewAuthHandler(authGroup, dbHandler)

	// reminder
	reminderGroup := r.Group("/reminder")
	reminderGroup.Use(middleware.JWTAuthHandler())
	reminder.NewReminderHandler(reminderGroup, dbHandler)
}
