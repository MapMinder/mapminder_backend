package user

import (
	"github.com/MapMinder/mapminder_backend/feature/user/handler"
	"github.com/MapMinder/mapminder_backend/feature/user/repository"
	"github.com/MapMinder/mapminder_backend/feature/user/usecase"
	infrastructure "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewUserHandler(router *gin.RouterGroup, db *gorm.DB) {
	// repositories
	userRepository := repository.NewUserRepository(db)

	// manager
	transactionManager := infrastructure.NewTransactionManager(db)

	// usecase
	userUsecase := usecase.NewUserUsecase(transactionManager, userRepository)

	// handler
	userHandler := handler.NewUserHandler(userUsecase)

	userHandler.RegisterRoutes(router)
}
