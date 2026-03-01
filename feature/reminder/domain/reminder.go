package domain

import "time"

type (
	ReminderStatus string
	Radius         float64
)

const (
	CreatedStatus = ReminderStatus("created")
	DefaultRadius = 200.00
)

type Reminder struct {
	ReminderId      string     `json:"reminder_id" gorm:"reminder_id;primaryKey"`
	UserId          string     `json:"user_id" gorm:"user_id"`
	Title           string     `json:"title" gorm:"title"`
	Description     string     `json:"description" gorm:"description"`
	Latitude        float64    `json:"latitude" gorm:"latitude"`
	Longitude       float64    `json:"longitude" gorm:"longitude"`
	Radius          float64    `json:"radius" gorm:"radius"`
	Status          string     `json:"status" gorm:"status"`
	LastTriggeredAt *time.Time `json:"last_triggered_at" gorm:"last_triggered_at"`
	CompletedAt     *time.Time `json:"completed_at" gorm:"completed_at"`
}

func (Reminder) TableName() string {
	return "reminder"
}
