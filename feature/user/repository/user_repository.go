package repository

import (
	"github.com/MapMinder/mapminder_backend/feature/user/domain"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user domain.User, tx *gorm.DB) (err error)
	GetUser(userId string) (user domain.User, err error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user domain.User, tx *gorm.DB) (err error) {
	if tx != nil {
		r.db = tx
	}

	if err = r.db.Create(&user).Error; err != nil {
		logger.Errorw("Error creating user: ", err)
		err = apperror.Internal()
		return
	}

	return
}

func (r *userRepository) GetUser(userId string) (user domain.User, err error) {
	if err = r.db.Take(&user).Where("user_id = ?", userId).Error; err != nil {
		logger.Errorw("Error finding user: ", err)
		return
	}
	return
}
