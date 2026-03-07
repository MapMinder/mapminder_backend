package dto

import "github.com/MapMinder/mapminder_backend/internal/status"

// GoogleSignIn
type GoogleSignInToken struct {
	IdToken string `json:"id_token" validate:"required"`
}

type GoogleSignInRes struct {
	Status status.Status
	Result string `json:"result"`
}
