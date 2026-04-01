package utils

import (
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
)

// ErrCode is a custom type for error codes used in the application.
type ErrCode string

// Define constants for different error codes that can be used throughout the application.
const (
	ErrCodeBadRequest      ErrCode = "BAD_REQUEST"
	ErrCodeNotFound        ErrCode = "NOT_FOUND"
	ErrCodeInternal        ErrCode = "INTERNAL_SERVER_ERROR"
	ErrCodeConflict        ErrCode = "CONFLICT"
	ErrCodeUnauthorized    ErrCode = "UNAUTHORIZED"
	ErrCodeTooManyRequests ErrCode = "TOO_MANY_REQUESTS"
)

// AppError is a custom error type that includes a message, an error code, and an optional underlying error.
type AppError struct {
	Message string
	Code    ErrCode
	Err     error
}

func (ae *AppError) Error() string {
	if ae.Err != nil {
		return ae.Message + ": " + ae.Err.Error()
	}
	return ae.Message
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

func ResponseErrorAbort(ctx *gin.Context, err error) {
	ResponseError(ctx, err)
	ctx.Abort()
}

func WriteGRPCErrorToGin(c *gin.Context, err error) {
	st, ok := status.FromError(err)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch st.Code() {

	case codes.InvalidArgument:

		errors := map[string]string{}

		for _, detail := range st.Details() {

			switch d := detail.(type) {

			case *errdetails.BadRequest:

				for _, v := range d.FieldViolations {
					errors[v.Field] = v.Description
				}
			}
		}

		if len(errors) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "VALIDATION_ERROR",
				"message": st.Message(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "VALIDATION_ERROR",
			"detail": errors,
		})
		return

	case codes.Unavailable:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":  "Service is currently unavailable. Please try again later.",
			"detail": st.Message(),
			"retry":  true,
		})
		return

	case codes.Internal:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "An internal server error occurred. Please try again later.",
			"detail": st.Message(),
		})
		return

	case codes.FailedPrecondition:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Failed precondition. Please check your request and try again.",
			"detail": st.Message(),
		})
		return

	default:
		c.JSON(httpStatusFromGrpcCode(st.Code()), gin.H{
			"error":  st.Code().String(),
			"detail": st.Message(),
		})
		return
	}
}

func ResponseGRPCErrorAbort(ctx *gin.Context, err error) {
	WriteGRPCErrorToGin(ctx, err)
	ctx.Abort()
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
	case ErrCodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// httpStatusFromGrpcCode maps gRPC codes to HTTP status codes.
func httpStatusFromGrpcCode(code codes.Code) int {
	switch code {
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
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
