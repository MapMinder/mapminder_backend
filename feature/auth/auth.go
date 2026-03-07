package auth

import (
	"github.com/MapMinder/mapminder_backend/feature/auth/handler"
	"github.com/MapMinder/mapminder_backend/feature/auth/repository"
	"github.com/MapMinder/mapminder_backend/feature/auth/usecase"
	usrRepository "github.com/MapMinder/mapminder_backend/feature/user/repository"
	infrastructure "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	uuidgenerator "github.com/MapMinder/mapminder_backend/shared/uuid_manager"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewAuthHandler(router *gin.RouterGroup, db *gorm.DB) {
	// repositories
	oauthRepository := repository.NewOauthTokenRepository(db)
	googleApiRepository := repository.NewGoogleSignInApiRepository()
	userRespository := usrRepository.NewUserRepository(db)

	// manager
	transactionManager := infrastructure.NewTransactionManager(db)
	uuidManager := uuidgenerator.NewUUIDGenerator()

	// usecse
	googleUsecase := usecase.NewGoogleSignInUsecase(oauthRepository, googleApiRepository, userRespository, uuidManager, transactionManager)

	// handler
	googleHandler := handler.NewGoogleSignInhandler(googleUsecase)

	// initialize routes
	googleHandler.RegisterRoutes(router)
}
