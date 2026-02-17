package repository

import "gorm.io/gorm"

type GoogleSignInRepository interface{}

type googleSignInRepository struct {
	db *gorm.DB
}

func NewGoogleSignInRepository(db *gorm.DB) GoogleSignInRepository {
	return &googleSignInRepository{
		db: db,
	}
}
