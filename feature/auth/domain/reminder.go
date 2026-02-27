package domain

import (
	"time"
)

// Reminder status constants
const (
	ReminderStatusActive    = "active"
	ReminderStatusCompleted = "completed"
)

type Reminder struct {
	ReminderId     string    `gorm:"reminder_id;primaryKey"`
	UserId         string    `gorm:"column:user_id"`
	Title          string    `gorm:"column:title"`
	Description    string    `gorm:"column:description"`
	Latitude       float64   `gorm:"column:latitude"`
	Longitude      float64   `gorm:"column:longitude"`
	Radius         float64   `gorm:"column:radius"`
	Status         string    `gorm:"column:status"`
	LastTriggeredAt *time.Time `gorm:"column:last_triggered_at"`
	CompletedAt    *time.Time `gorm:"column:completed_at"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Reminder) TableName() string {
	return "reminder"
}