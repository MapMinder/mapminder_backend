package usecase

import (
	"context"

	"github.com/MapMinder/mapminder_backend/feature/user/domain"
	"github.com/MapMinder/mapminder_backend/feature/user/repository"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	"github.com/MapMinder/mapminder_backend/shared/middleware"
	"github.com/MapMinder/mapminder_backend/shared/tx"
)

type UserUsecase interface {
	GetUser(ctx context.Context) (user domain.User, err error)
}

type userUsecase struct {
	TxManager      tx.Manager
	UserRepository repository.UserRepository
}

func NewUserUsecase(txManager tx.Manager, userRepository repository.UserRepository) *userUsecase {
	return &userUsecase{
		TxManager:      txManager,
		UserRepository: userRepository,
	}
}

func (u *userUsecase) GetUser(ctx context.Context) (user domain.User, err error) {
	logger.Info("Create User Usecase")

	err = u.TxManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		userId := middleware.UserIDFromContext(txCtx)

		user, err = u.UserRepository.GetUser(txCtx, userId)
		if err != nil {
			return err
		}
		return nil
	})
	return
}
