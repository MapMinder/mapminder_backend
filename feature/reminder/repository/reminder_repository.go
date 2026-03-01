package repository

import (
	"context"

	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	tx "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"gorm.io/gorm"
)

type ReminderRepository interface {
	Create(ctx context.Context, reminder domain.Reminder) (err error)
}

type reminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{
		db: db,
	}
}

// Create
func (r *reminderRepository) Create(ctx context.Context, reminder domain.Reminder) (err error) {
	db := tx.ExtractTx(ctx, r.db)
	if db == nil {
		db = r.db
	}

	if err = db.Create(&reminder).Error; err != nil {
		logger.Errorw("Error creating reminder: ", err)
		err = apperror.Internal()
		return
	}
	return
}
