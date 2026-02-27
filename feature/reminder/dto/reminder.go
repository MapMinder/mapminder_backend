package dto

import (
	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	"github.com/MapMinder/mapminder_backend/internal/status"
)

type Reminder struct {
	Title       string  `json:"title" validate:"required;min=1"`
	Description string  `json:"description" validate:"required;min=1"`
	Latitude    float32 `json:"latitude" validate:"required"`
	Longitude   float32 `json:"longitude" validate:"required"`
	Radius      float32 `json:"radius" validate:"required"`
}

type ReminderRes struct {
	Status status.Status
	Result domain.Reminder `json:"result"`
}
