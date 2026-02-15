package auth

import (
	"github.com/MapMinder/mapminder_backend/feature/auth/handler"
	"github.com/MapMinder/mapminder_backend/feature/auth/repository"
	"github.com/MapMinder/mapminder_backend/feature/auth/usecase"
	"github.com/gin-gonic/gin"
)

func NewAuthHandler(r *gin.Engine, router *gin.RouterGroup) {
	googleRepository := repository.NewGoogleSignInRepository()
	googleUsecase := usecase.NewHealthUsecase(googleRepository)
	googleHandler := handler.NewGoogleSignInhandler(googleUsecase)

	// initialize routes
	googleHandler.RegisterRoutes(router)
}
