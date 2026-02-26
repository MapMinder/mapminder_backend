package repository

import (
	"gorm.io/gorm"
)

type ReminderRepository interface{}

type reminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{
		db: db,
	}
}

// Create
func (r *reminderRepository) Create() (err error) {
	return
}
