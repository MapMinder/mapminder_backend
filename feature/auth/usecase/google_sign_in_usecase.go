package usecase

import (
	"context"
	"os"
	"time"

	"github.com/MapMinder/mapminder_backend/feature/auth/domain"
	"github.com/MapMinder/mapminder_backend/feature/auth/repository"
	usrDomain "github.com/MapMinder/mapminder_backend/feature/user/domain"
	usrRepository "github.com/MapMinder/mapminder_backend/feature/user/repository"
	"github.com/MapMinder/mapminder_backend/internal/logger"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/MapMinder/mapminder_backend/shared/tx"
	uuidgenerator "github.com/MapMinder/mapminder_backend/shared/uuid_manager"
	"github.com/golang-jwt/jwt/v5"
)

type GoogleSignInUsecase interface {
	GoogleSignIn(ctx context.Context, idToken string) (signedJWTToken string, err error)
}

type googleSignInUsecase struct {
	TxManager            tx.Manager
	OauthTokenRepository repository.OauthTokenRepository
	ApiRepository        repository.GoogleSignInApiRepository
	UserRepository       usrRepository.UserRepository
	UUIDGenerator        uuidgenerator.UUIDManager
}

func NewGoogleSignInUsecase(
	oauthTokenRepository repository.OauthTokenRepository,
	apiRepository repository.GoogleSignInApiRepository,
	userRepository usrRepository.UserRepository,
	uuidGenerator uuidgenerator.UUIDManager,
	txManager tx.Manager,
) GoogleSignInUsecase {
	return &googleSignInUsecase{
		TxManager:            txManager,
		OauthTokenRepository: oauthTokenRepository,
		ApiRepository:        apiRepository,
		UserRepository:       userRepository,
		UUIDGenerator:        uuidGenerator,
	}
}

// GoogleSignIn
func (u *googleSignInUsecase) GoogleSignIn(ctx context.Context, idToken string) (signedJWTToken string, err error) {
	logger.Info("google sign in usecase")

	err = u.TxManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		// get Google client id
		clientId := os.Getenv("GOOGLE_CLIENT_ID")
		creds := domain.GoogleSignInCreds{
			ClientId: clientId,
			IdToken:  idToken,
		}

		googleSignInClaims, err := u.ApiRepository.GetGoogleUserInfo(txCtx, creds)
		if err != nil {
			return err
		}

		// check user exists with subject
		oauthToken, err := u.OauthTokenRepository.GetOauthInformation(txCtx, googleSignInClaims.Claims.Subject)
		if err != nil {
			return err
		}

		var user usrDomain.User
		if len(oauthToken.UserId) <= 0 {
			user, err = u.createUserRelatedData(txCtx, googleSignInClaims)
			if err != nil {
				logger.Errorw("Error occurred while creating user uuid: ", err)
				return err
			}
		} else {
			user, err = u.UserRepository.GetUser(oauthToken.UserId)
			if err != nil {
				return err
			}
		}

		signedJWTToken, err = u.createJWTToken(user.UserId)
		if err != nil {
			return err
		}

		return err
	})
	return
}

// createUserRelatedData
func (u *googleSignInUsecase) createUserRelatedData(ctx context.Context, googleSignInClaims domain.GoogleIdToken) (user usrDomain.User, err error) {
	logger.Info("Creating New User")

	newUserUUID, err := u.UUIDGenerator.NewV7()
	if err != nil {
		logger.Errorw("Error occurred while creating user uuid: ", err)
		err = apperror.Internal()
		return
	}

	user = usrDomain.User{
		UserId:   newUserUUID,
		Email:    googleSignInClaims.Claims.Email,
		Username: googleSignInClaims.Claims.Name,
	}

	err = u.UserRepository.Create(ctx, user)
	if err != nil {
		return
	}

	// create record for oauth_token table
	oauthToken := domain.OauthToken{
		UserId:          newUserUUID,
		OauthProvider:   domain.GoogleProvider,
		OauthProviderId: googleSignInClaims.Subject,
	}

	err = u.OauthTokenRepository.CreateOauthToken(ctx, oauthToken)
	if err != nil {
		return
	}

	return
}

// createJWTToken creates jwt token
func (u *googleSignInUsecase) createJWTToken(userId string) (signedJWTToken string, err error) {
	logger.Infof("Create JWTToken for user %s", userId)

	signKey := os.Getenv("JWT_SIGN_KEY")

	// HACK: expiry for jwt token is 30 days
	// TODO: in the future make it smaller and implement refresh token
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "mapminder",
		Subject:   userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// token.SignedString expects byte
	signedJWTToken, err = token.SignedString([]byte(signKey))
	if err != nil {
		logger.Errorw("Failed to sign JWT token: ", err)
	}
	return
}
