package mapper

import (
	"github.com/MapMinder/mapminder_backend/feature/reminder/domain"
	"github.com/MapMinder/mapminder_backend/feature/reminder/dto"
)

// MapReminders maps multiple []domain.Reminder to []dto.ReminderResStruct
func MapReminders(targetReminders []domain.Reminder) []dto.ReminderResStruct {
	reminders := make([]dto.ReminderResStruct, 0, len(targetReminders))

	for _, reminder := range targetReminders {
		reminders = append(reminders, dto.ReminderResStruct{
			ReminderId:  reminder.ReminderId,
			Title:       reminder.Title,
			Description: reminder.Description,
			Latitude:    reminder.Latitude,
			Longitude:   reminder.Longitude,
			Radius:      reminder.Radius,
			Status:      reminder.Status,
		})
	}
	return reminders
}

// MapReminder maps domain.Reminder to dto.ReminderResStruct
func MapReminder(targetReminder domain.Reminder) (reminder dto.ReminderResStruct) {
	return dto.ReminderResStruct{
		ReminderId:  reminder.ReminderId,
		Title:       reminder.Title,
		Description: reminder.Description,
		Latitude:    reminder.Latitude,
		Longitude:   reminder.Longitude,
		Radius:      reminder.Radius,
		Status:      reminder.Status,
	}
}
