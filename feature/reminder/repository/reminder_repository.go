package repository

import (
	"context"
	"errors"

	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	tx "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"gorm.io/gorm"
)

type ReminderRepository interface {
	Create(ctx context.Context, reminder domain.Reminder) (err error)
	GetReminder(ctx context.Context, reminderId string) (reiminder domain.Reminder, err error)
	GetReminders(ctx context.Context, userId string, status string) (reminders []domain.Reminder, err error)
	DeleteReminder(ctx context.Context, reminderId string) (err error)
	UpdateReminder(ctx context.Context, reminder domain.Reminder) (err error)
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

// GetReminder
func (r *reminderRepository) GetReminder(ctx context.Context, reminderId string) (reminder domain.Reminder, err error) {
	logger.Infof("reminder repository: GetReminder")

	db := tx.ExtractTx(ctx, r.db)
	if db == nil {
		db = r.db
	}

	if err = db.Where("reminder_id = ?", reminderId).Take(&reminder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Infof("No records found for given reminderId")
			err = apperror.NotFound()
			return
		} else {
			logger.Errorw("Internal error occurred: ", err)
			err = apperror.Internal()
			return
		}
	}
	return
}

func (r *reminderRepository) GetReminders(ctx context.Context, userId string, status string) (reminders []domain.Reminder, err error) {
	logger.Infof("reminder repository : GetReminders")

	db := tx.ExtractTx(ctx, r.db)
	if db == nil {
		db = r.db
	}

	query := db.Where("user_id = ?", userId)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err = query.Find(&reminders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Infof("No records found for given reminderId")
			err = apperror.NotFound()
			return
		} else {
			logger.Errorw("Internal error occurred: ", err)
			err = apperror.Internal()
			return
		}
	}

	return
}

func (r *reminderRepository) DeleteReminder(ctx context.Context, reminderId string) (err error) {
	logger.Infof("reminder repository: DeleteReminder")

	db := tx.ExtractTx(ctx, r.db)
	if db == nil {
		db = r.db
	}

	res := db.WithContext(ctx).Where("reminder_id = ?", reminderId).Delete(&domain.Reminder{})
	if res.Error != nil {
		logger.Errorw("Internal error occurred: ", res.Error)
		err = apperror.Internal()
		return
	}

	if res.RowsAffected == 0 {
		logger.Infow("No records found", "remidner id: ", reminderId)
		err = apperror.NotFound()
		return
	}
	return
}

func (r *reminderRepository) UpdateReminder(ctx context.Context, reminder domain.Reminder) (err error) {
	logger.Infof("reminder repository: UpdateReminder")

	logger.Infof("reminder repository: UpdateReminder: %+v", reminder) // ← add this
	db := tx.ExtractTx(ctx, r.db)
	if db == nil {
		db = r.db
	}

	if err = db.Model(domain.Reminder{}).Where("reminder_id = ?", reminder.ReminderId).Updates(&reminder).Error; err != nil {
		logger.Errorw("Internal error occurred: ", err)
		err = apperror.Internal()
		return
	}
	return
}
