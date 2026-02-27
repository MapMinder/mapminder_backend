package domain

import "time"

type ReminderStatus string

const CreatedStatus = ReminderStatus("created")

type Reminder struct {
	ReminderId      string    `gorm:"reminder_id;primaryKey"`
	UserId          string    `gorm:"user_id"`
	Title           string    `gorm:"title"`
	Description     string    `gorm:"description"`
	Latitude        float32   `gorm:"latitude"`
	Longitude       float32   `gorm:"longitude"`
	Radius          float32   `gorm:"radius"`
	Status          string    `gorm:"status"`
	LastTriggeredAt time.Time `gorm:"last_triggered_at"`
	CompletedAt     time.Time `gorm:"completed_at"`
}

func (Reminder) TableName() string {
	return "reminder"
}
