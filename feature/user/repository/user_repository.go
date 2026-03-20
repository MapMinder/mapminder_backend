package repository

import (
	"context"

	"github.com/MapMinder/mapminder_backend/feature/user/domain"
	tx "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (err error)
	GetUser(ctx context.Context, userId string) (user domain.User, err error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) (err error) {
	db := tx.ExtractTx(ctx, r.db)
	if db == nil {
		db = r.db
	}

	if err = db.Create(&user).Error; err != nil {
		logger.Errorw("Error creating user: ", err)
		err = apperror.Internal()
		return
	}

	return
}

func (r *userRepository) GetUser(ctx context.Context, userId string) (user domain.User, err error) {
	db := tx.ExtractTx(ctx, r.db)
	if db == nil {
		db = r.db
	}

	if err = db.Where("user_id = ?", userId).Take(&user).Error; err != nil {
		logger.Errorw("Error finding user: ", err)
		err = apperror.Internal()
		return
	}
	return
}
