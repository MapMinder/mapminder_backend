package dto

// GoogleSignIn
type GoogleSignInToken struct {
	IdToken string `json:"id_token" validate:"required"`
}
