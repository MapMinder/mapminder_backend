package reminder

import (
	"github.com/MapMinder/mapminder_backend/feature/reminder/handler"
	"github.com/MapMinder/mapminder_backend/feature/reminder/repository"
	"github.com/MapMinder/mapminder_backend/feature/reminder/usecase"
	infrastructure "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	uuidgenerator "github.com/MapMinder/mapminder_backend/shared/uuid_manager"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewReminderHandler(router *gin.RouterGroup, db *gorm.DB) {
	// repositories
	reminderRepository := repository.NewReminderRepository(db)

	// manager
	// NOTE: i need to know if this is valid this looks like i am initializing this per feature need to cross check
	transactionManager := infrastructure.NewTransactionManager(db)
	uuidManager := uuidgenerator.NewUUIDGenerator()

	// usecase
	reminderUsecase := usecase.NewReminderUsecase(transactionManager, uuidManager, reminderRepository)

	// handler
	reminderHandler := handler.NewReminderHandler(reminderUsecase)

	// initialize routes
	reminderHandler.RegisterRoutes(router)
}
