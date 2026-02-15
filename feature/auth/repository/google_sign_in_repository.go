package repository

type GoogleSignInRepository interface{}

type googleSignInRepository struct{}

func NewGoogleSignInRepository() GoogleSignInRepository {
	return &googleSignInRepository{}
}
