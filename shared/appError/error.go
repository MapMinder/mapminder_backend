package apperror

import (
	"github.com/MapMinder/mapminder_backend/internal/status"
)

type Error struct {
	Status status.Status
	Err    error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Status.Message
}

func newError(s status.Status, args ...any) *Error {
	var err error

	return &Error{
		Status: s,
		Err:    err,
	}
}

func BadRequest(args ...any) *Error {
	return newError(status.BadRequest, args...)
}

func Unauthorized(args ...any) *Error {
	return newError(status.Unauthorized, args...)
}

func NotFound(args ...any) *Error {
	return newError(status.NotFound, args...)
}

func Internal(args ...any) *Error {
	return newError(status.InternalError, args...)
}

func GoogleError(args ...any) *Error {
	return newError(status.GoogleError, args...)
}

func JWTError(args ...any) *Error {
	return newError(status.JWTError, args...)
}
