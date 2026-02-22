package repository

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/MapMinder/mapminder_backend/feature/auth/domain"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
)

type GoogleSignInApiRepository interface {
	GetGoogleUserInfo(ctx context.Context, googleSignInCreds domain.GoogleSignInCreds) (googleIdToken domain.GoogleIdToken, err error)
}

type googleSignInApiRepository struct{}

func NewGoogleSignInApiRepository() GoogleSignInApiRepository {
	return &googleSignInApiRepository{}
}

// GetGoogleUserInfo get user's info form google
func (r *googleSignInApiRepository) GetGoogleUserInfo(ctx context.Context, googleSignInCreds domain.GoogleSignInCreds) (googleIdToken domain.GoogleIdToken, err error) {
	payload, err := idtoken.Validate(ctx, googleSignInCreds.IdToken, googleSignInCreds.ClientId)
	if err != nil {
		logger.Error(err)
		err = apperror.GoogleError()
		return
	}

	// idtoken validate returns the claims as map which is not clean marshal and unmarshal to make claims cleaner
	marshalPayload, err := json.Marshal(payload.Claims)
	if err != nil {
		logger.Errorw("Marshal error occurred", apperror.Internal())
		return
	}

	var claims domain.GoogleIdTokenClaims
	if err = json.Unmarshal(marshalPayload, &claims); err != nil {
		logger.Errorw("Unmarshal error occurred", apperror.Internal())
		return
	}

	googleIdToken = domain.GoogleIdToken{
		Issuer:    payload.Issuer,
		Audience:  payload.Audience,
		Subject:   payload.Subject,
		IssuedAt:  payload.IssuedAt,
		ExpiresAt: payload.Expires,
		Claims:    claims,
	}

	return
}
