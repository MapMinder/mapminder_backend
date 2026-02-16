package auth

import (
	"github.com/MapMinder/mapminder_backend/feature/auth/handler"
	"github.com/MapMinder/mapminder_backend/feature/auth/repository"
	"github.com/MapMinder/mapminder_backend/feature/auth/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAuthHandler(router *gin.RouterGroup, db *gorm.DB) {
	googleRepository := repository.NewGoogleSignInRepository(db)
	googleApiRepository := repository.NewGoogleSignInApiRepository()
	googleUsecase := usecase.NewHealthUsecase(googleRepository, googleApiRepository)
	googleHandler := handler.NewGoogleSignInhandler(googleUsecase)

	// initialize routes
	googleHandler.RegisterRoutes(router)
}
