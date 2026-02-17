package usecase

import (
	"context"
	"os"

	"github.com/MapMinder/mapminder_backend/feature/auth/domain"
	"github.com/MapMinder/mapminder_backend/feature/auth/repository"
	"github.com/MapMinder/mapminder_backend/internal/logger"
)

type GoogleSignInUsecase interface {
	GoogleSignIn(ctx context.Context, idToken string) (googleSignInClaims domain.GoogleIdToken, err error)
}

type googleSignInUsecase struct {
	GoogleSignInRepository repository.GoogleSignInRepository
	ApiRepository          repository.GoogleSignInApiRepository
}

func NewHealthUsecase(googleSignInRepository repository.GoogleSignInRepository, apiRepository repository.GoogleSignInApiRepository) GoogleSignInUsecase {
	return &googleSignInUsecase{
		GoogleSignInRepository: googleSignInRepository,
		ApiRepository:          apiRepository,
	}
}

func (u *googleSignInUsecase) GoogleSignIn(ctx context.Context, idToken string) (googleSignInClaims domain.GoogleIdToken, err error) {
	logger.Info("google ")

	// get Google client id
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	creds := domain.GoogleSignInCreds{
		ClientId: clientId,
		IdToken:  idToken,
	}

	googleSignInClaims, err = u.ApiRepository.GetGoogleUserInfo(ctx, creds)
	if err != nil {
		return
	}

	return
}
