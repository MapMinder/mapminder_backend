package usecase

type ReminderUsecase interface{}

type reminderUsecase struct{}

func NewReminderUsecase() *reminderUsecase {
	return &reminderUsecase{}
}

// Create Reminder
func (u *reminderUsecase) CreateReminder() (err error) {
	return
}
