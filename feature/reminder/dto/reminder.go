package dto

import (
	"time"

	"github.com/MapMinder/mapminder_backend/internal/status"
)

type Reminder struct {
	Title        string  `json:"title" validate:"required,min=1,max=65"`
	Description  string  `json:"description" validate:"required,min=1"`
	Latitude     float64 `json:"latitude" validate:"required,latitude"`
	Longitude    float64 `json:"longitude" validate:"required,longitude"`
	LocationName string  `json:"location_name" validate:"required,min=1"`
	// Radius      float64 `json:"radius" validate:"required"` the user's will not be able to set the radius it is a default value (atleast for the mvp)
}

type UpdateReminder struct {
	Title       string  `json:"title" validate:"max=65,omitempty"`
	Description string  `json:"description" validate:"omitempty"`
	Latitude    float64 `json:"latitude" validate:"latitude"`
	Longitude   float64 `json:"longitude" validate:"longitude"`
	Status      string  `json:"status" validate:"status"`
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

type UpdateLastReminderRes struct {
	Status status.Status
	Result Notify `json:"result"`
}

type Notify struct {
	Notify bool `json:"notify"`
}

type ReminderResStruct struct {
	ReminderId   string     `json:"reminder_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Latitude     float64    `json:"latitude"`
	Longitude    float64    `json:"longitude"`
	LocationName string     `json:"location_name"`
	Radius       float64    `json:"radius"`
	Status       string     `json:"status"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
