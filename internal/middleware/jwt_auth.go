package middleware

import (
	"context"
	"os"
	"strings"

	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// TODO: move this to the shared
const userIDKey contextKey = "user_id"

// JWTAuthHandler
func JWTAuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(apperror.JWTError())
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(apperror.JWTError())
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(parts[1], &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			signKey := os.Getenv("JWT_SIGN_KEY")
			if signKey == "" {
				c.Error(apperror.JWTError())
				c.Abort()
				return nil, apperror.JWTError()
			}
			return []byte(signKey), nil
		})
		if err != nil {
			c.Error(apperror.JWTError())
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			c.Error(apperror.JWTError())
			c.Abort()
			return
		}
		userId := claims.Subject
		c.Set(string(userIDKey), userId)
		ctx := context.WithValue(c.Request.Context(), userIDKey, userId)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// TODO: move this to the shared with userIdKey
// UserIDFromContext
func UserIDFromContext(ctx context.Context) string {
	userId, _ := ctx.Value(userIDKey).(string)
	return userId
}
