package dto

import (
	"github.com/MapMinder/mapminder_backend/internal/status"
)

type Reminder struct {
	Title       string  `json:"title" validate:"required,min=1"`
	Description string  `json:"description" validate:"required,min=1"`
	Latitude    float64 `json:"latitude" validate:"required,latitude"`
	Longitude   float64 `json:"longitude" validate:"required,longitude"`
	// Radius      float64 `json:"radius" validate:"required"` the use's will not be able to set the radius it is a default value (atleast for the mvp)
}

type UpdateReminder struct {
	Title       string  `json:"title" validate:"omitempty"`
	Description string  `json:"description" validate:"omitempty"`
	Latitude    float64 `json:"latitude" validate:"latitude,omitempty"`
	Longitude   float64 `json:"longitude" validate:"longitude,omitempty"`
}

type ReminderRes struct {
	Status status.Status
	Result ReminderResStruct `json:"result"`
}
type RemindersRes struct {
	Status status.Status
	Total  int
	Result []ReminderResStruct `json:"result"`
}

type ReminderResStruct struct {
	ReminderId  string  `json:"reminder_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Radius      float64 `json:"radius"`
	Status      string  `json:"status"`
}
