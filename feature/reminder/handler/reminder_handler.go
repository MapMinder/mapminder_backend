package handler

import "github.com/gin-gonic/gin"

type ReminderHandler struct{}

func NewReminderHandler() *ReminderHandler {
	return &ReminderHandler{}
}

func (h *ReminderHandler) CreateReminder() (err error) {
	return
}

func (h *ReminderHandler) RegisterRoutes(r *gin.RouterGroup) {}
