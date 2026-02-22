package repository

import (
	"context"
	"errors"

	"github.com/MapMinder/mapminder_backend/feature/auth/domain"
	tx "github.com/MapMinder/mapminder_backend/internal/infrastructure/database"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"gorm.io/gorm"
)

type OauthTokenRepository interface {
	GetOauthInformation(ctx context.Context, sub string) (oauthToken domain.OauthToken, err error)
	CreateOauthToken(ctx context.Context, oauthToken domain.OauthToken) (err error)
}

type oauthTokenRepository struct {
	db *gorm.DB
}

func NewOauthTokenRepository(db *gorm.DB) OauthTokenRepository {
	return &oauthTokenRepository{
		db: db,
	}
}

// GetOauthInformation()
func (r *oauthTokenRepository) GetOauthInformation(ctx context.Context, sub string) (oauthToken domain.OauthToken, err error) {
	logger.Infof("oauthToken repository")

	if err = r.db.Where("oauth_provider_id = ?", sub).Take(&oauthToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Infof("No records found for given sub")
			err = nil
			return
		} else {
			logger.Errorw("Internal error occurred: ", err)
			err = apperror.Internal()
			return
		}
	}
	return
}

// CreateOauthToken
func (r *oauthTokenRepository) CreateOauthToken(ctx context.Context, oauthToken domain.OauthToken) (err error) {
	db := tx.ExtractTx(ctx, r.db)
	if db == nil {
		db = r.db
	}
	if err = db.Create(&oauthToken).Error; err != nil {
		logger.Errorw("Error creating oauth: ", err)
		err = apperror.Internal()
		return
	}
	return
}
