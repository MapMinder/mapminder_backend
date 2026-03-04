package handler

import (
	"net/http"

	"github.com/MapMinder/mapminder_backend/feature/reminder/dto"
	"github.com/MapMinder/mapminder_backend/feature/reminder/mapper"
	"github.com/MapMinder/mapminder_backend/feature/reminder/usecase"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/MapMinder/mapminder_backend/internal/status"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/MapMinder/mapminder_backend/shared/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReminderHandler struct {
	ReminderUsecase usecase.ReminderUsecase
}

func NewReminderHandler(reminderUsecase usecase.ReminderUsecase) *ReminderHandler {
	return &ReminderHandler{
		ReminderUsecase: reminderUsecase,
	}
}

func (h *ReminderHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("", h.CreateReminder)
	r.GET("/:reminder_id", h.GetReminder)
	r.GET("", h.GetReminders)
	r.DELETE("/:reminder_id", h.DeleteReminder)
}

func (h *ReminderHandler) CreateReminder(r *gin.Context) {
	logger.Info("Create Reminder Handler")
	ctx := r.Request.Context()

	var req dto.Reminder
	if err := r.ShouldBind(&req); err != nil {
		logger.Errorw("Failed To Bind Request", err)
		r.Error(apperror.BadRequest())
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		logger.Errorw("Failed To Validate Request", err)
		r.Error(apperror.BadRequest())
		return
	}

	reminder, err := h.ReminderUsecase.CreateReminder(ctx, req)
	if err != nil {
		logger.Error(err)
		r.Error(err)
		return
	}

	res := dto.ReminderRes{
		Status: status.Created,
		Result: dto.ReminderResStruct{
			ReminderId:  reminder.ReminderId,
			Title:       reminder.Title,
			Description: reminder.Description,
			Latitude:    reminder.Latitude,
			Longitude:   reminder.Longitude,
			Radius:      reminder.Radius,
			Status:      reminder.Status,
		},
	}

	r.JSON(http.StatusCreated, res)
}

func (h *ReminderHandler) GetReminder(r *gin.Context) {
	logger.Infof("reminde handler: GetReminder")
	ctx := r.Request.Context()

	reminderId := r.Param("reminder_id")
	_, err := uuid.Parse(reminderId)
	if err != nil {
		logger.Errorw("Invalid ReminderId: ", err)
		err = apperror.BadRequest()
		r.Error(err)
		return
	}

	reminder, err := h.ReminderUsecase.GetReminder(ctx, reminderId)
	if err != nil {
		r.Error(err)
		return
	}

	res := dto.ReminderRes{
		Status: status.Success,
		Result: mapper.MapReminder(reminder),
	}

	r.JSON(http.StatusOK, res)
}

func (h *ReminderHandler) GetReminders(r *gin.Context) {
	logger.Infof("reminde handler: GetReminders")
	ctx := r.Request.Context()

	requestParam := r.Query("status")

	reminders, err := h.ReminderUsecase.GetReminders(ctx, requestParam)
	if err != nil {
		r.Error(err)
		return
	}

	res := dto.RemindersRes{
		Status: status.Success,
		Total:  len(reminders),
		Result: mapper.MapReminders(reminders),
	}

	r.JSON(http.StatusOK, res)
}

func (h *ReminderHandler) DeleteReminder(r *gin.Context) {
	logger.Infof("reminde handler: DeleteReminder")
	ctx := r.Request.Context()

	reminderId := r.Param("reminder_id")
	_, err := uuid.Parse(reminderId)
	if err != nil {
		logger.Errorw("Invalid ReminderId: ", err)
		err = apperror.BadRequest()
		r.Error(err)
		return
	}

	err = h.ReminderUsecase.DeleteReminder(ctx, reminderId)
	if err != nil {
		r.Error(err)
		return
	}

	r.JSON(http.StatusOK, status.Success)
}
