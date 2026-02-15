package usecase

import "github.com/MapMinder/mapminder_backend/feature/auth/repository"

type GoogleSignInUsecase interface{}

type googleSignInUsecase struct {
	GoogleSignInRepository repository.GoogleSignInRepository
}

func NewHealthUsecase(googleSignInRepository repository.GoogleSignInRepository) GoogleSignInUsecase {
	return &googleSignInUsecase{
		GoogleSignInRepository: googleSignInRepository,
	}
}
