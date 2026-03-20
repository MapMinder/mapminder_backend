package mapper

import (
	"github.com/MapMinder/mapminder_backend/feature/user/domain"
	"github.com/MapMinder/mapminder_backend/feature/user/dto"
)

func MapUser(targetUser domain.User) dto.UserResStruct {
	return dto.UserResStruct{
		UserId:   targetUser.UserId,
		Username: targetUser.Username,
		Email:    targetUser.UserId,
	}
}
