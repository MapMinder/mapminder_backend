package domain

type User struct {
	UserId   string `gorm:"user_id"`
	Username string `gorm:"username"`
	Email    string `gorm:"email"`
}

func (User) TableName() string {
	return "user"
}
