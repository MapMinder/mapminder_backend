package dto

import "github.com/MapMinder/mapminder_backend/internal/status"

type UserRes struct {
	Status status.Status
	Result UserResStruct `json:"result"`
}

type UserResStruct struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
