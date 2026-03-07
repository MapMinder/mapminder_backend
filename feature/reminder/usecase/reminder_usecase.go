package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	"github.com/MapMinder/mapminder_backend/feature/reminder/dto"
	"github.com/MapMinder/mapminder_backend/feature/reminder/mapper"
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
	GetReminder(ctx context.Context, reminderId string) (reminder domain.Reminder, err error)
	GetReminders(ctx context.Context, status string) (reminders []domain.Reminder, err error)
	DeleteReminder(ctx context.Context, reminderId string) (err error)
	UpdateReminder(ctx context.Context, reminderId string, reminder dto.UpdateReminder) (updatedReminder domain.Reminder, err error)
	UpdateLastTriggeredAt(ctx context.Context, reminderId string) (shouldNotify bool, err error)
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
			Status:      string(domain.ActiveStatus),
		}

		err = u.ReminderRepository.Create(txCtx, reminder)
		if err != nil {
			err = apperror.Internal()
			return err
		}
		return nil
	})

	return
}

func (u *reminderUsecase) GetReminder(ctx context.Context, reminderId string) (reminder domain.Reminder, err error) {
	logger.Infof("reminder usecase: GetReminder")

	reminder, err = u.ReminderRepository.GetReminder(ctx, reminderId)
	if err != nil {
		return
	}
	return
}

func (u *reminderUsecase) GetReminders(ctx context.Context, status string) (reminders []domain.Reminder, err error) {
	logger.Infof("reminder usecase: GetReminders")
	userId := middleware.UserIDFromContext(ctx)

	if status != "" && !domain.IsValidStatus(status) {
		logger.Errorw("Invalid status provided: ", apperror.BadRequest(), "status", status)
		err = apperror.BadRequest()
		return
	}

	reminders, err = u.ReminderRepository.GetReminders(ctx, userId, status)
	if errors.Is(err, apperror.NotFound()) {
		// if records are not found it could mean that the user has never created a reminder
		// this is not an error so we return no errors
		err = nil
		return
	}
	if err != nil {
		return
	}
	return
}

func (u *reminderUsecase) DeleteReminder(ctx context.Context, reminderId string) (err error) {
	logger.Infof("reminder usecase: DeleteReminder")

	err = u.TxManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		err = u.ReminderRepository.DeleteReminder(ctx, reminderId)
		return nil
	})
	return
}

func (u reminderUsecase) UpdateReminder(ctx context.Context, reminderId string, reminder dto.UpdateReminder) (updatedReminder domain.Reminder, err error) {
	logger.Infof("reminder usecase: UpdateReminder")

	userId := middleware.UserIDFromContext(ctx)

	err = u.TxManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		oldReminder, err := u.ReminderRepository.GetReminder(txCtx, reminderId)
		if err != nil {
			// ReminderRepoisitory.GetReminder returns not found error and logs so no need to explictly handle here
			return err
		}

		if oldReminder.UserId != userId {
			logger.Errorw("Invalid access reminder does not belong to user: ", apperror.Unauthorized(), "reminder_id: ", reminderId, "user_id: ", userId)
			err = apperror.Unauthorized()
			return err
		}

		err = u.ReminderRepository.UpdateReminder(txCtx, mapper.MapReminderFromDTOForUpdate(reminderId, reminder))
		if err != nil {
			return err
		}

		updatedReminder, err = u.ReminderRepository.GetReminder(txCtx, reminderId)
		if err != nil {
			return err
		}

		return nil
	})
	return
}

func (u reminderUsecase) UpdateLastTriggeredAt(ctx context.Context, reminderId string) (shouldNotify bool, err error) {
	logger.Infof("reminder usecase: UpdateReminder")

	userId := middleware.UserIDFromContext(ctx)

	err = u.TxManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		reminder, err := u.ReminderRepository.GetReminder(txCtx, reminderId)
		if err != nil {
			return err
		}

		if reminder.UserId != userId {
			logger.Errorw("Invalid access reminder does not belong to user: ", apperror.Unauthorized(), "reminder_id: ", reminderId, "user_id: ", userId)
			err = apperror.Unauthorized()
			return err
		}

		now := u.TimeProvider.Now()
		shouldNotify = u.shouldNotify(reminder, now)
		if !shouldNotify {
			logger.Infof("Cooldown in process or reminder's status is paused no notification will be sent")
			return nil
		}

		err = u.ReminderRepository.UpdateLastTriggeredAt(txCtx, reminderId, now)
		if err != nil {
			return err
		}

		return nil
	})
	return
}

func (u reminderUsecase) shouldNotify(reminder domain.Reminder, now time.Time) (shouldNotify bool) {
	if reminder.Status != string(domain.ActiveStatus) {
		shouldNotify = false
		return
	}

	if reminder.LastTriggeredAt == nil {
		shouldNotify = true
		return
	}

	tenMinsBeforeNow := now.Add(-10 * time.Minute)
	shouldNotify = reminder.LastTriggeredAt.Before(tenMinsBeforeNow)

	return
}
