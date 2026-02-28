package handler

import (
	"net/http"

	"github.com/MapMinder/mapminder_backend/feature/reminder/dto"
	"github.com/MapMinder/mapminder_backend/feature/reminder/usecase"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/MapMinder/mapminder_backend/internal/status"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/MapMinder/mapminder_backend/shared/validator"
	"github.com/gin-gonic/gin"
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
	r.POST("/", h.CreateReminder)
}

func (h *ReminderHandler) CreateReminder(r *gin.Context) {
	logger.Info("Create Reminder Handler")
	ctx := r.Request.Context()

	var req dto.Reminder
	if err := r.ShouldBind(&req); err != nil {
		logger.Errorw("Find To Bind Request", err)
		r.Error(apperror.BadRequest())
		return
	}

	if err := validator.ValidateStruct(req); err != nil {
		logger.Errorw("Failed To Bild Request", err)
		r.Error(err)
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
		Result: reminder,
	}

	r.JSON(http.StatusOK, res)
}
