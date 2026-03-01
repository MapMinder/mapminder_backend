package usecase

import (
	"context"

	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	"github.com/MapMinder/mapminder_backend/feature/reminder/dto"
	"github.com/MapMinder/mapminder_backend/feature/reminder/repository"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/MapMinder/mapminder_backend/shared/middleware"
	timeProvider "github.com/MapMinder/mapminder_backend/shared/time"
	"github.com/MapMinder/mapminder_backend/shared/tx"
	uuidgenerator "github.com/MapMinder/mapminder_backend/shared/uuid_manager"
)

type ReminderUsecase interface {
	CreateReminder(ctx context.Context, params dto.Reminder) (reminder domain.Reminder, err error)
}

type reminderUsecase struct {
	TxManager     tx.Manager
	UUIDgenerator uuidgenerator.UUIDManager
	TimeProvider  timeProvider.RealTimeProvider

	ReminderRepository repository.ReminderRepository
}

func NewReminderUsecase(txManager tx.Manager, uuidGenerator uuidgenerator.UUIDManager, reminderRepository repository.ReminderRepository, timeProvider timeProvider.RealTimeProvider) *reminderUsecase {
	return &reminderUsecase{
		TxManager:          txManager,
		UUIDgenerator:      uuidGenerator,
		TimeProvider:       timeProvider,
		ReminderRepository: reminderRepository,
	}
}

// Create Reminder
func (u *reminderUsecase) CreateReminder(ctx context.Context, params dto.Reminder) (reminder domain.Reminder, err error) {
	logger.Info("Create Reminder Usecase")

	err = u.TxManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// get user data from the auth
		userId := middleware.UserIDFromContext(txCtx)

		reminderId, err := u.UUIDgenerator.NewV7()
		if err != nil {
			logger.Errorw("Error Creating reminder_id with error: ", err, "userId", userId)
			err = apperror.Internal()
			return err
		}

		reminder = domain.Reminder{
			ReminderId:  reminderId,
			UserId:      userId,
			Title:       params.Title,
			Description: params.Description,
			Latitude:    params.Latitude,
			Longitude:   params.Longitude,
			Radius:      domain.DefaultRadius,
			Status:      string(domain.CreatedStatus),
		}

		err = u.ReminderRepository.Create(txCtx, reminder)
		if err != nil {
			err = apperror.Internal()
			return err
		}
		return nil
	})
	if err != nil {
		return
	}

	return
}
