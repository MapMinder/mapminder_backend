package middleware

import (
	"errors"

	"github.com/MapMinder/mapminder_backend/internal/status"
	apperror "github.com/MapMinder/mapminder_backend/shared/appError"
	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.Error
		if errors.As(err, &appErr) {
			c.JSON(appErr.Status.Code, gin.H{
				"error": appErr.Status.Message,
			})
			return
		}

		c.JSON(status.InternalError.Code, gin.H{
			"error": status.InternalError.Message,
		})
	}
}
