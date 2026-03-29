package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrCode is a custom type for error codes used in the application.
type ErrCode string

// Define constants for different error codes that can be used throughout the application.
const (
	ErrCodeBadRequest   ErrCode = "BAD_REQUEST"
	ErrCodeNotFound     ErrCode = "NOT_FOUND"
	ErrCodeInternal     ErrCode = "INTERNAL_SERVER_ERROR"
	ErrCodeConflict     ErrCode = "CONFLICT"
	ErrCodeUnauthorized ErrCode = "UNAUTHORIZED"
)

// AppError is a custom error type that includes a message, an error code, and an optional underlying error.
type AppError struct {
	Message string
	Code    ErrCode
	Err     error
}

func (ae *AppError) Error() string {
	return ""
}

// NewError creates a new AppError with the given message and error code.
func NewError(message string, code ErrCode) error {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

// WrapError creates a new AppError that wraps an existing error with a new message and error code.
func WrapError(err error, message string, code ErrCode) error {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// ResponseError is a helper function that takes a Gin context and an error, checks if the error is an AppError,
// and responds with the appropriate HTTP status code and JSON response based on the error code and message.
// If the error is not an AppError, it responds with a generic internal server error message.
func ResponseError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appErr.Code)
		response := gin.H{
			"error": appErr.Message,
			"code":  appErr.Code,
		}

		if appErr.Err != nil {
			response["details"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrCodeInternal,
	})
}

// httpStatusFromCode is a helper function that maps custom error codes to corresponding HTTP status codes.
func httpStatusFromCode(code ErrCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// ResponseSuccess is a helper function that takes a Gin context, an HTTP status code, and any data,
// and responds with a JSON object containing a "status" field set to "success" and a "data" field containing the provided data.
func ResponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, gin.H{
		"status": "success",
		"data":   data,
	})
}

func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.JSON(status, gin.H{})
}

func ResponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
}
