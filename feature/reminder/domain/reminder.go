package domain

import "time"

type (
	ReminderStatus string
	Radius         float64
)

const (
	CreatedStatus = ReminderStatus("created")
	DefaultRadius = 200
)

type Reminder struct {
	ReminderId      string    `gorm:"reminder_id;primaryKey"`
	UserId          string    `gorm:"user_id"`
	Title           string    `gorm:"title"`
	Description     string    `gorm:"description"`
	Latitude        float64   `gorm:"latitude"`
	Longitude       float64   `gorm:"longitude"`
	Radius          float64   `gorm:"radius"`
	Status          string    `gorm:"status"`
	LastTriggeredAt time.Time `gorm:"last_triggered_at"`
	CompletedAt     time.Time `gorm:"completed_at"`
}

func (Reminder) TableName() string {
	return "reminder"
}
