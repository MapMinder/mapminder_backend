package reminder

import (
	"github.com/MapMinder/mapminder_backend/feature/reminder/handler"
	"github.com/MapMinder/mapminder_backend/feature/reminder/repository"
	"github.com/MapMinder/mapminder_backend/feature/reminder/usecase"
	infrastructure "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	timeProvider "github.com/MapMinder/mapminder_backend/shared/time"
	uuidgenerator "github.com/MapMinder/mapminder_backend/shared/uuid_manager"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewReminderHandler(router *gin.RouterGroup, db *gorm.DB) {
	// repositories
	reminderRepository := repository.NewReminderRepository(db)

	// manager
	transactionManager := infrastructure.NewTransactionManager(db)
	uuidManager := uuidgenerator.NewUUIDGenerator()
	timeProvider := timeProvider.NewTimeProvider()

	// usecase
	reminderUsecase := usecase.NewReminderUsecase(transactionManager, uuidManager, reminderRepository, timeProvider)

	// handler
	reminderHandler := handler.NewReminderHandler(reminderUsecase)

	// initialize routes
	reminderHandler.RegisterRoutes(router)
}
