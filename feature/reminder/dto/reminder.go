package dto

import (
	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	"github.com/MapMinder/mapminder_backend/internal/status"
)

type Reminder struct {
	Title       string  `json:"title" validate:"required,min=1"`
	Description string  `json:"description" validate:"required,min=1"`
	Latitude    float64 `json:"latitude" validate:"required,latitude"`
	Longitude   float64 `json:"longitude" validate:"required,longitude"`
	// Radius      float64 `json:"radius" validate:"required"` the use's will not be able to set the radius it is a default value (atleast for the mvp)
}

// TODO: the dto should not know about the domain objects adding the
// domain.Reminder here means dto is accessing a business object
// will fix this in separate ticket(or a dedicated ticket for refactor)
type ReminderRes struct {
	Status status.Status
	Result domain.Reminder `json:"result"`
}
