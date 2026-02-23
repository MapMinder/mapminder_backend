package status

import (
	"net/http"
)

type Status struct {
	Code    int
	Message string
}

// NOTE: every time error's are added they should reflect on the appError package
var (
	Success       = Status{Code: http.StatusOK, Message: "Success"}
	Created       = Status{Code: http.StatusCreated, Message: "Created"}
	BadRequest    = Status{Code: http.StatusBadRequest, Message: "Invalid request"}
	Unauthorized  = Status{Code: http.StatusUnauthorized, Message: "Unauthorized"}
	NotFound      = Status{Code: http.StatusNotFound, Message: "Resource not found"}
	InternalError = Status{Code: http.StatusInternalServerError, Message: "Internal server error"}
	GoogleError   = Status{Code: http.StatusBadRequest, Message: "Falied to get required information"}
	JWTError      = Status{Code: http.StatusBadRequest, Message: "JWT token is malformed or invalid"}
)
