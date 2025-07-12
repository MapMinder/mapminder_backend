package health

import (
	"github.com/MapMinder/mapminder_backend/feature/health/handler"
	"github.com/MapMinder/mapminder_backend/feature/health/repository"
	"github.com/MapMinder/mapminder_backend/feature/health/usecase"

	"github.com/gin-gonic/gin"
)

func NewHealthHandler(r *gin.Engine) {
	healthRepo := repository.NewHealthRepository()
	healthUsecase := usecase.NewHealthUsecase(healthRepo)
	healthHandler := handler.NewHealthHandler(healthUsecase)

	healthHandler.RegisterRoutes(r)
}
