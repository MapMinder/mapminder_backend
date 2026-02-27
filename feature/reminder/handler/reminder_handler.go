package handler

import (
	"github.com/MapMinder/mapminder_backend/feature/reminder/usecase"
	"github.com/MapMinder/mapminder_backend/internal/logger"
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
}
