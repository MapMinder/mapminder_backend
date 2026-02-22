package uuidgenerator

import (
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/google/uuid"
)

type UUIDManager interface {
	NewV7() (string, error)
}

type uuidManager struct{}

func NewUUIDGenerator() UUIDManager {
	return &uuidManager{}
}

func (g *uuidManager) NewV7() (string, error) {
	u, err := uuid.NewV7()
	if err != nil {
		logger.Errorw("Error occurred while creating user uuid: ", err)
		err = apperror.Internal()
		return u.String(), err
	}
	return u.String(), err
}
